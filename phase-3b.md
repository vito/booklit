# Phase 3b: mdx-as-template

Replace Go `html/template` as the user-facing template engine with mdx
(Markdown + JSX) templates. Authors write one syntax — the same JSX +
Dang they already use for content — and templates become discoverable
by dropping a `.md` file in `html/` (or `templates/`). The tier-3 Dang
dispatch added in Phase 3a is the closure-capturing mechanism templates
build on.

## Why

Phase 3a left two dialects in play for user-facing extension:

1. A Dang function (`pub Foo(...)`) dispatched via tier-3.
2. A Go `html/template` file (`Foo.tmpl`) wrapping a `booklit.Styled`
   emitted by tier-4 fallback.

Pick one. The Go-template option pulls authors into a separate mini-
language with its own scoping, escaping rules, and helper function
registry — exactly the kind of accidental complexity the original
pivot set out to remove. mdx templates reuse the JSX+Dang machinery
that already exists. Templates become "a JSX/MDX file that gets
evaluated with the component's props in scope and `children` as a
block."

Concrete payoff for the dogfooded docs: **`<Define>` in
`docs/booklitdoc/` can move to `Define.md`**, removing 30+ lines of
content-construction Go. The remaining helpers in `cmd/booklit-docs`
either follow the same pattern or fall into the Dagger / "punted"
buckets enumerated in `jsx-dang.md`'s Phase 3a entry.

## Target shape

A template file under `html/Card.md`:

```jsx
<div class="card">
  <h3>{title}</h3>
  <div class="body">{children}</div>
</div>
```

A doc author's invocation:

```jsx
<Card title="Hello">
  Welcome to the page.
</Card>
```

Renders as:

```html
<div class="card">
  <h3>Hello</h3>
  <div class="body"><p>Welcome to the page.</p></div>
</div>
```

`title` is a Dang value (the prop, bridged from JSX). `children` is a
block arg — a Dang callable holding the original JSX children
("Welcome to the page."). Dang auto-calls 0-arity callables when
interpolated, so `{children}` evaluates the children inline at that
position. Lowercase `<div>` / `<h3>` pass through as raw HTML. The
prose inside the body still goes through Markdown.

A more substantial template — what `Define.md` should look like:

```jsx
<div class="definition">
  <Target tag={tag}/>
  <pre class="definition-thumb">
    <Syntax language="html">{sig}</Syntax>
  </pre>
  <div class="definition-content">
    {children}
  </div>
</div>
```

Used as `<Define tag="card" sig="<Card title=…>x</Card>">A card.</Define>`,
this gives the same output as today's Go helper, with no Go.

## What stays vs. what changes vs. what goes away

| Stays | Changes | Goes away |
|---|---|---|
| JSX syntax in content | Template *engine* | Go `html/template` for user-facing |
| Tier-1 (built-in) and tier-3 (Dang) dispatch | Tier-4 (template) — now mdx, not `.tmpl` | `booklit.Styled` as the auto-wrap fallback |
| Renderer-internal templates (`render/html/*.tmpl`) | Template file extension: `.tmpl` → `.md` | `{{.Partial "x"}}` / `{{.Content | render}}` syntax for user files |
| `<Children/>` (new built-in) and `{children}` interpolation are both valid ways to emit body | Dispatch precedence (built-in → Dang → mdx template) | Implicit "render anything to Styled and let the renderer hunt for a template" fallback |

The renderer-internal templates (`styled.tmpl`, `section.tmpl`,
`page.tmpl`, etc., bundled in `render/html/`) stay as Go templates.
They're framework infrastructure, not user extension. Converting them
is a separate question and may never be worth it.

## Architecture sketch

Pipeline change (only tier-4 is new):

```
parse (goldmark + JSX inline parser)
  → AST (Section tree with JSXElement invocations)
  → evaluate
      tier-1: built-in dispatch (Go function in builtins/)
      tier-2: deferred (currently Dang in dispatch order, see Q1)
      tier-3: Dang function dispatch (pub Foo with optional &body)
      tier-4: mdx template — look up html/Foo.md, parse, evaluate
              with props bound + `children` as a block
      tier-5: error (no implicit Styled fallback)
  → render (existing HTML engine; renderer-internal templates only)
```

The mdx template engine is small: it composes the existing JSX parser
(marklit), the existing JSX dispatcher (stages.Evaluate), and the
existing dangeval scope-pushing machinery. The template's evaluation
is *just another JSX evaluation*, with props pre-bound in Dang scope
and `children` as a callable block. No new template language, no new
syntax.

## File-level impact

### New code

- `templates/` (or similar): a small package that wraps "load + cache
  + evaluate" of `.md` template files. Reuses `marklit.Parse` and
  `stages.Evaluate`.

- `builtins/children.go`: a `<Children/>` built-in. Looks up the
  current scope's `children` block (the same one tier-3 dispatch
  sets up) and invokes it. Useful for cases where the value-style
  `{children}` doesn't fit (e.g. as a JSX child position).

### Modified

- `stages/evaluate.go::VisitJSXElement`: add a fourth tier after Dang.
  If a `.md` template exists for the component name, dispatch to it
  by constructing an evaluator scope with props bound as Dang vars
  and `children` as a componentBlock-style closure over the original
  JSX children. The renderer-fallback Styled path goes away.

- `dangeval/component.go`: factor the props-bind + body-block setup
  out of tier-3 dispatch so the template path can reuse it. The
  componentBlock pattern is the same; only the "what we call" differs
  (a Dang function vs. an mdx template's evaluator).

- `docs/booklitdoc/booklitdoc.go`: shrink further as `Define.md`,
  `Godoc.md`, and any other portable helpers move to templates. Final
  state: only helpers that are genuinely blocked on Dagger (LitSyntax)
  or AST introspection (Columns) and the chroma palette init().

### Deletions

- `docs/html/Card.tmpl`, `docs/html/columns.tmpl`,
  `docs/html/definition.tmpl`, `docs/html/output-frame.tmpl` (already
  renamed to `OutputFrame.tmpl` in Phase 3a), etc. — converted to
  `.md` siblings. Each is a small mechanical migration.

- `tests/fixtures/*.tmpl` — the test fixtures used by template-
  fallback tests. Some can become `.md` templates; others might stay
  (the test suite is its own concern).

## Phases

### Phase 3b-1: template engine MVP

End-to-end with one fixture. Don't migrate any real templates yet.

- Load `html/Hello.md` (or whatever).
- Dispatch `<Hello name="world"/>` to it.
- Verify `{name}` interpolates and `{children}` emits the JSX
  children (test both empty and non-empty children).
- Tests in `tests/template_test.go` covering: prop interpolation,
  `{children}` with and without content, `<Children/>` built-in
  equivalence, nested JSX inside the template, multi-line templates.

Outcome: an `.md` template can serve as a JSX component with the
same expressive range as a Go template + Styled today.

### Phase 3b-2: port `<Define>`

The forcing function. Define exercises everything: prop interpolation
(`{tag}`, `{sig}`), `<Target>` invocation from within a template,
`<Syntax>` invocation with prop-as-content, `{children}` for the
description, lowercase HTML pass-through, and multi-line content.

- Write `docs/html/Define.md` matching the target-shape sketch.
- Remove `defineFunc` from `docs/booklitdoc/booklitdoc.go`.
- Rebuild docs site; diff rendered output against pre-port to confirm
  parity.

Outcome: `<Define>` runs with zero Go.

### Phase 3b-3: port remaining easy helpers

- `<Godoc ref="pkg.Symbol"/>` → `Godoc.md`. Needs a string-split
  helper somewhere — either a tiny Dang helper in scope
  (`pub splitOnce(s: String!, sep: String!): [String!]! { … }`) or a
  template-side helper if we add one.
- `<Columns>` — stays a built-in (AST introspection, see jsx-dang.md
  Phase 3a entry).
- `<LitSyntax>` — stays a built-in until Dagger. (Wait, can we
  punt this differently? The chroma+regex pipeline could be its own
  small Dang function library… but that's a separate scope.)

Outcome: `cmd/booklit-docs` shrinks to the genuinely-Go-bound
helpers (`<Columns>`, `<LitSyntax>`, chroma palette init).

### Phase 3b-4: decide the renderer-internal templates

`render/html/*.tmpl` are bundled into the binary and handle base
primitives: `styled.tmpl`, `section.tmpl`, `page.tmpl`, `link.tmpl`,
etc. Two options:

- **Leave as Go templates.** They're framework infrastructure;
  authors don't touch them. The dialect distinction is invisible to
  users.
- **Convert to mdx.** Removes one engine entirely; codebase is
  consistent. But there's a chicken-and-egg moment: mdx templates
  rely on the JSX evaluator, which rely on these primitives. May
  need a small bootstrap.

Recommendation: defer. Leave as Go templates until there's a concrete
reason to convert.

## Open questions

### Q1. Dispatch order — does mdx template go before or after Dang?

Phase 3a order is: built-in → Dang → template-fallback. For Phase
3b, two reasonable orderings:

- (a) built-in → Dang → mdx template. Dang wins if both exist.
  Authors who explicitly write `pub Foo(...)` get it. Templates are
  the "default presentation" tier.
- (b) built-in → mdx template → Dang. Templates win if both exist.
  Matches the original plan ("template first, Dang second").

Recommendation: (a). A Dang function is more specific than a
template (the user wrote actual code); templates are the "drop a
file" path. Matches Phase 3a's existing ordering.

### Q2. Where do templates live?

Today: `html/*.tmpl` (configured via `--html-templates`). Keep that
path? Rename to `templates/`? Allow both Go and mdx during a grace
period?

Recommendation: keep `html/` for the directory; `.md` extension for
mdx templates. Hard cutover — no Go templates accepted from the user
side. (Same call we made in Phase 3a Q5.)

### Q3. How does `children` interact with markdown re-parsing?

JSX children are AST nodes (already parsed). When `children()` is
invoked from inside a template, the closure re-enters
`stages.Evaluate` on those nodes. That evaluation happens in the
template's evaluator context, with whatever scope is active at the
invocation site.

Question: does the re-evaluation preserve the original parsing
context, or does it re-process the children's text? Per Phase 3a
the children are `[]ast.Node` (already parsed by marklit). So
re-evaluation visits the same nodes. No re-parsing. Should be safe
but worth a fixture.

### Q4. Template caching and invalidation

Per Phase 1: marklit parsing is cached by modification time. Same
strategy for templates: parse `html/Foo.md` once, cache the AST, re-
parse on file change (matters for serve mode).

Question: is parsing actually a hot path? Even uncached, a small
`.md` template is microseconds. Defer caching until benchmarks
justify it.

### Q5. Side-effect semantics of `{children}`

Phase 3a's componentBlock returns `NullValue` and appends content
via a captured accumulator. For tier-3, the accumulator IS the
element's content.

For Phase 3b templates, `{children}` interpolation evaluates as a
Dang expression. Dang's auto-call invokes the 0-arity block. The
block re-enters the JSX evaluator with the original children. The
content needs to land at the `{children}` position in the template's
output.

Mechanism: the componentBlock's closure captures the template's
*current* evaluator, and the block's invocation appends content to
that evaluator's Result. The Dang expression returns `NullValue`;
the bridge gives empty content; the side-effect content was already
appended at the right position. Should work because Dang evaluates
the expression eagerly before VisitJSXExpression returns.

Edge case: `{1 + children}` or `{someTransform(children)}` — if
children's value is needed downstream rather than just rendered in
place, the side-effect model breaks. For Phase 3b, error on this or
just produce wrong output; v1 doesn't need to handle it.

### Q6. Naming: `children` everywhere?

The user-facing convention is `children`. Implementation-side we
already have `componentBlock` (the Go type). The block parameter
*name* in tier-3 Dang functions is whatever the author writes
(`&body(item: ...)`, `&block(x: ...)`, etc.). For mdx templates, the
block is always called `children` — Booklit picks the name, not the
template author. So template authors can always write `{children}`
or `<Children/>` and expect it to work.

### Q7. `<Children/>` built-in: how does it find the block?

The block lives in the current Dang eval scope (bound by name
`children`). `<Children/>` is a built-in: it can call into the
current evaluator's Dang env to look up `children`, invoke it, and
emit the resulting content. Confirm during implementation that
this lookup path is clean — might want a dedicated helper on
`dangeval.Evaluator`.

### Q8. Error reporting

Templates are real files with line numbers. When `Define.md` has a
syntax error, the error should point at `Define.md:line:col`. When
a Dang expression inside the template fails, same. This dovetails
with the source-mapped-errors follow-up from Phase 3.

Probably defer the polish to the same future phase; for v1, plain
error wrapping is OK.

### Q9. Coexistence with the current template fallback

The current tier-4 wraps in `booklit.Styled{Style: Name}`. Three
options:

- (a) Hard cutover: tier-4 is mdx template lookup; no Styled fallback.
  If no template, error.
- (b) Try mdx first, fall back to Styled if no `.md` exists.
- (c) Reverse: try Styled first (renderer looks up `.tmpl`), then
  mdx.

(a) matches the Phase 0 / 3a posture (hard cutover). (b)/(c) preserve
backwards compat at the cost of complexity. Recommendation: (a).
Tests will catch any regression and we own all the templates.

### Q10. What about `<Children name="foo"/>` (partials)?

Go templates use `{{.Partial "thumb"}}` for named partials —
sub-sections of the children. JSX doesn't have an analogous concept
natively; the author would have to pass partials as props
(`<Card title={…} footer={…}>body</Card>`).

For Phase 3b, drop partials as a concept. Anything you'd use a
partial for becomes a prop. Simpler model, fewer features, but
covers everything we've seen.

## Relevant existing code

Before starting, read:

- `stages/evaluate.go::VisitJSXElement` — the dispatch site. Tier-4
  goes here.
- `stages/evaluate.go::dispatchDang` — the closure-capture pattern
  for the body block. Template dispatch is structurally similar.
- `dangeval/component.go` — componentBlock, WithBindings,
  CallComponent. The block plumbing is reusable.
- `marklit/marklit.go` — how documents get parsed; the same path
  parses templates.
- `render/html.go` — existing template engine (for the renderer-
  internal templates that stay). Worth understanding so we don't
  accidentally break it.
- `builtins/styled.go` — pattern for adding `<Children/>`.
- `docs/booklitdoc/booklitdoc.go` — the `<Define>` helper that's
  the forcing function for 3b-2.
- `docs/html/Card.tmpl`, `docs/html/columns.tmpl`,
  `docs/html/definition.tmpl` — what existing templates look like,
  for migration reference.

For Dang's auto-call semantics, look at `pkg/dang/eval.go` near
`IsAutoCallable` (mentioned in the Phase 3a survey). A 0-arity
callable interpolated as a value should self-invoke.

## Definition of done

- `Define` (in `docs/booklitdoc/`) is replaced by `Define.md` and the
  docs site builds with identical visible output.
- Any other booklitdoc helper that doesn't need chroma/regex/AST-
  introspection is replaced by an `.md` template.
- A new component is addable by dropping `html/Foo.md` (no Go, no
  rebuild needed in serve mode).
- `go test ./...` green; existing JSX tests still pass; new template
  engine tests cover the matrix of prop interpolation, body
  emission via `{children}` and `<Children/>`, and nested JSX.
- `cmd/booklit-docs` shrinks to just `<Columns>` + `<LitSyntax>` +
  chroma palette init (all blocked on Dagger or AST introspection).

## Out of scope

- Renderer-internal template conversion (deferred, see Phase 3b-4).
- Source-mapped errors for Dang/template failures (Phase 3 follow-up).
- Per-section scope for `{expr}` (still deferred from Phase 3).
- Dagger story for the gnarly helpers (separate phase entirely).
- mdx templates outside the user-extension boundary (e.g. allowing
  authors to override `section.tmpl`).
- A separate template registry or namespace; templates are just files
  in the templates directory.

## What to do first

1. Read `dangeval/component.go` and
   `stages/evaluate.go::dispatchDang` end to end. The block-arg
   mechanism is the load-bearing piece; the template engine is a
   thin wrapper.

2. Build a minimal `templates/` package: `Load(dir, name) (ast.Node, error)`
   parses `<dir>/<name>.md` via marklit and returns the AST.

3. Wire a tier-4 path in `VisitJSXElement`: if a template exists,
   construct a dispatch shape that mirrors `dispatchDang` but uses
   the template's AST as the "function body" and binds `children`
   as the block. Props become Dang scope bindings (use
   `dangeval.WithBindings`).

4. Write `tests/template_test.go` with the matrix from Phase 3b-1.

5. Once green, port `<Define>` (Phase 3b-2). If it works, this whole
   plan is validated.

6. Iterate on remaining helpers.

## Progress log

This section is append-only. Each session that makes meaningful
progress should add a new dated `###` entry at the bottom — do not
edit earlier entries or earlier sections of this plan, even when
later decisions contradict them.

### 2026-05-31 — Phase 3b-1 through 3b-3 landed

Tier-4 mdx-template dispatch is live, `<Define>` and `<Godoc>` are
ported, and the docs site builds with two fewer Go helpers.

**Done.**

- New `templates/` package: `Registry` loads `<dir>/<Name>.md` lazily,
  caches the parsed AST by mtime, and trims a single trailing newline
  from the file so multiple invocations don't pile up blank lines.
  `Parse` is a hand-written tokenizer (not goldmark) that recognizes
  three things and nothing else: `<Pascal ...>` JSX elements, `{expr}`
  expressions, and raw HTML text. No Markdown processing — templates
  are HTML scaffolding around prop holes.

- Raw-HTML text chunks emit `ast.JSXElement{Name: "RawHTML"}` whose
  body is the literal text. A new `builtins/raw_html.go` registers
  `<RawHTML>` which wraps in `Styled{Style: "raw-html"}`. Templates
  emit this via a private `templates.rawHTML` ast.Node whose `Visit`
  forwards to the JSXElement dispatch. The public `<RawHTML>` builtin
  doubles as a content-author escape hatch.

- `<Children/>` built-in (`builtins/children.go`) looks up `children`
  in the active Dang scope and emits its content. Equivalent to
  `{children}` as an expression but works in JSX child position where
  an expression sometimes reads awkwardly. Plumbed via a new `Dang`
  field on `builtins.Context`.

- `dangeval.ContentValue` is a Dang `Value` that carries booklit.Content
  verbatim. The bridge unwraps it without going through string-form
  flattening. `{children}` resolves to a `ContentValue` so nested
  styling like `<Italic>{children}</Italic>` survives.

- `stages.Evaluate.dispatchTemplate` is the new tier-4: render children
  eagerly into a `ContentValue`, bridge props by name, push everything
  via `WithBindings`, and visit the template's AST with a sub-evaluator.
  Dispatch order is built-in → Dang → mdx template → legacy Styled
  fallback. Templates win over the Styled fallback (Q9 = (a)-ish:
  Styled stays during this phase as a safety net, but templates take
  precedence). Dang functions still beat templates (Q1 = (a)).

- `tests/template_test.go` (9 cases): prop interpolation, `{children}`,
  `<Children/>`, empty children, nested JSX, multi-line + raw HTML,
  template-beats-fallback, Dang-beats-template, expression props.

- `docs/html/Define.md` (5 lines of mdx) replaces the 40-line `defineFunc`.
  Uses `<Target tag={tag}/><Syntax language="html">{sig}</Syntax>` inside
  the existing `definition` styled wrapper. The rendered definition
  blocks are byte-identical to the Go version.

- `docs/html/Godoc.md` (one line) replaces `godocFunc`. String-split
  logic moved into `docs/lit/helpers.dang` as `godocPkg`,
  `godocPkgDisplay`, `godocSymbol`, `godocURL` Dang functions. They
  compose `String.split` / `List.takeFirst` / `List.dropFirst` /
  `List.join` / `String.trimLeft` to do what `strings.SplitN` +
  `strings.TrimLeft` did. The rendered Godoc links are byte-identical.

`cmd/booklit-docs` now registers only `<Columns>` and `<LitSyntax>`
plus the chroma palette init() — matching the "definition of done"
inventory in this plan.

**Decisions worth knowing.**

- *Custom template parser, not goldmark.* Goldmark's HTML block parser
  swallows `<div class="card">...</div>` as opaque raw HTML — `{title}`
  and `{children}` inside would never be processed. Writing a small
  byte-tokenizer specific to templates was cleaner than extending
  goldmark to recognize `{expr}` inside HTML blocks. The trade-off is
  no Markdown features in templates (no `# heading`, no `- list`, no
  fenced code blocks); templates are HTML scaffolding, not prose.

- *Eager children render, not side-effect closure.* The plan's Q5
  proposes a side-effect closure that captures the template's evaluator
  and appends content at the `{children}` site. The simpler approach —
  evaluate children once into a `ContentValue` and bind it — fits the
  current Dang shape better: a custom auto-callable Callable would need
  a real `FunctionType` in the type env, and the side-effect model
  raises subtle questions about what `{children + "x"}` should mean.
  `ContentValue` declares its type as `String!` so the inferrer accepts
  it; the bridge takes the `Content` field directly. Pure-emission
  cases (`{children}`, `<Children/>`) work; expressions that try to
  treat children as a real string will fail at runtime, which is fine
  for v1.

- *RawHTML as a public builtin.* The template parser needs to emit raw
  HTML chunks. The cleanest representation is a JSX element with a
  registered built-in. Making `<RawHTML>` public means content authors
  can use it as a Markdown escape hatch — that's a feature, not a
  leak. Templates emit it via a private `templates.rawHTML` node whose
  Visit forwards to the same JSX dispatch, so the synthesis is opaque
  to anyone reading a template.

- *Trailing newline trim in the loader.* Templates conventionally end
  with `\n`. Without trimming, every `<Define>` invocation produced a
  trailing blank line, drifting the rendered HTML from the original
  Go version. The loader strips one trailing `\n` (interior whitespace
  is preserved verbatim because templates are HTML-significant).

- *Define dropped the "highlighted title for references" feature.* The
  old `defineFunc` set `Target.Title` to a syntax-highlighted
  `<TagName>` so `<Reference tag="title"/>` displayed as
  `<a href="..."><code>&lt;Title&gt;</code></a>`. The mdx
  `<Target tag={tag}/>` falls back to `Title = String(tag)`, so
  references now show plain text (`<a href="...">title</a>`). Visually
  the definitions themselves are byte-identical; only references
  to-definitions look different. The fix would be a small docs-side
  built-in that constructs a Target with a chroma-highlighted title,
  but that contradicts the "Define runs with zero Go" outcome.
  Trade-off accepted; the docs still work, links still link.

- *Helpers.dang lives at `docs/lit/helpers.dang`.* That's where
  `dangeval.New` scans for *.dang files (non-recursive, alphabetical;
  matches what Phase 3a established). No new convention introduced.

- *Dang comments use `#`, not `//`.* Discovered on first run of
  `helpers.dang`. Not Dang-specific knowledge to keep, just a note for
  the next person.

**Known limitations / follow-ups.**

- *Card.tmpl and TemplateLink.tmpl in `tests/fixtures/` and
  `docs/html/` still use the legacy Styled fallback.* Templates win
  over the fallback when both exist, but the fallback is still in place
  as a safety net. Phase 3b's "Goes away" column for `booklit.Styled`
  as auto-wrap is therefore aspirational — flipping the hard cutover
  (`Q9 = (a)`) means porting every `.tmpl` to `.md` first or accepting
  test breakage. Defer until there's a reason to push.

- *Plan Q1 (built-in → Dang → mdx) is implemented; the legacy Styled
  fallback sits as tier-5.* The plan envisioned tier-5 = error. Same
  reasoning as the bullet above — flip when the fallback is unused.

- *Reference-display regression for Define.* See decisions above.

- *No source-mapped errors yet.* If `Define.md:3:5` has a Dang typo,
  the error wraps as `evaluating template <Define>: evaluating {bad}:
  ...`. Useful enough; mapping to template-file line/col is a future
  polish.

- *Renderer-internal templates (`render/html/*.tmpl`, plus
  `docs/html/page.tmpl` etc.) remain Go templates.* Phase 3b-4
  recommended deferring; no reason found to convert.

**Reasonable next steps.** (a) Decide whether to push for the hard
cutover on the Styled fallback — would need a sweep of remaining
`.tmpl` files and a deletion of the tier-5 path in `stages/evaluate.go`.
(b) Source-mapped errors for template/Dang failures (deferred from
Phase 3). (c) Restore the highlighted-title-for-references feature if
it turns out to matter (small docs-side built-in, ~20 lines).

### 2026-05-31 — late session: cutover + control flow + Target rework

Three commits landed picking up where Phase 3b-3 left off:

- `0ce7b03` (`refactor(builtins): make Target children the rich title`):
  `<Target>` now uses its children as the rich Title content (matching
  baselit's legacy `\target{tag}{title}{content}` shape). The `title`
  prop becomes a string-shorthand that loses to children when both are
  given. Define.md uses
  `<Target tag={tag}><Syntax language="html"><{componentName(tag: tag)}>
  </Syntax></Target>` to restore the syntax-highlighted `<TagName>`
  title that the original Go `defineFunc` set, eliminating the
  reference-display regression flagged in the previous entry. The
  kebab→PascalCase conversion lives in a new `componentName` Dang
  helper (composed from `String.split` / `List.map` / a small
  `capitalize` that leans on `split(separator: "", limit: 2)` since
  Dang has no first-character method).

- `97de1f4` (`refactor: hard-cutover the JSX dispatch — no more Styled
  fallback`): implements plan Q9 (a). The dispatch order is now
  built-in → Dang → mdx template → error. Unknown JSX components are
  a Load-time error, not a Render-time "missing .tmpl" panic. To
  unblock, the three remaining JSX-fallback templates ported to
  mdx: `tests/fixtures/Card.tmpl` → `Card.md`,
  `docs/html/OutputFrame.tmpl` → `OutputFrame.md`,
  `docs/html/TemplateLink.tmpl` → `TemplateLink.md`.
  `templates.Registry` now takes a list of search directories so the
  test harness can layer the per-test tempdir over `tests/fixtures/`.
  `docs/lit/plugins.md` was rewritten to teach the `.md` form;
  `tests/dang_dispatch_test.go` cases that asserted on the render-time
  fallback failure now assert on the Load-time error. Renderer-internal
  templates (page, section, sidebar, splash, big-code, code-*,
  lit-*, columns) stay as Go templates — they're framework
  infrastructure rendered by built-ins that explicitly emit
  `Styled{Style: …}`. `2fb1896` follows up by removing `definition.tmpl`
  since `Define.md` now emits the HTML scaffolding directly.

- `0ab313e` (`feat(builtins): add <For>, <If>, <Unless>`): packages the
  tier-3 body-block + `WithBindings` pattern as built-ins so authors
  don't need a per-project `pub Each(items, &body)` helper. Five test
  cases cover string iteration, record iteration with a custom `as`
  binding name, conditional truthiness for `<If>` / `<Unless>`, and
  empty-list behavior. `eb5aaaf` was the proof-point test that this
  pattern handles record iteration end-to-end via tier-3 before the
  built-in landed.

**Decisions worth knowing.**

- *Target children are Title, not Content.* The prior shape (children
  → Content) was a JSX-migration accident; baselit's legacy variadic
  put title first. Existing tests only assert on Target's rendered
  anchor (`<a id="middle"></a>`), so the migration was safe — but it
  does mean Target no longer carries a Content field for search-index
  previews. If/when that matters, add a `content` prop.

- *Hard cutover was uneventful.* Two test cases needed updating from
  RenderErr to LoadErr; nothing else broke. Confirms the Styled
  fallback was load-bearing for fewer cases than the safety-net
  rhetoric suggested. The Q9 = (a) choice is right.

- *Built-in control-flow over per-project Dang helpers.* `<For>` and
  `<If>` could be left to author-provided Dang helpers (the test
  fixtures showed `pub Each(items: [a!]!, &body(item: a!): Boolean!)`
  works). Making them built-in avoids the body-return-type quirk
  ("null is not allowed for Boolean!") that the explicit-return Dang
  idiom papers over, gives a stable name across projects, and removes
  one piece of mandatory per-project boilerplate. Trade-off: two more
  reserved JSX names. Acceptable.

- *`<For>` binding name is `item` by default, overridable with
  `as="…"`.* Mirrors JSX/MDX conventions (`for…in`, `useState`); the
  default keeps simple cases terse and the override avoids shadowing
  when nested.

**Known limitations / follow-ups.**

- *Phase 4 (Dagger dispatch) intentionally deferred from this
  autonomous session.* See `decisions.md` 2026-05-31 entry: the
  `<Foo from="github.com/.../mod"/>` syntax requires runtime
  construction of `dang.ImportConfig` with a Dagger graphql client +
  schema. Existing dangeval `dagger.json` auto-import already covers
  the in-tree case; the missing piece is JSX-tag syntax for
  out-of-tree modules. Deferred to a co-design pass with the user.

- *JSX-in-Dang (the `{items.map(t => <Foo>{t.x}</Foo>)}` headline
  from the very top of `jsx-dang.md`) is no longer load-bearing.* The
  iteration use case is now covered by `<For each={items} as="t">…
  </For>` without Dang grammar changes. Conditionals by `<If
  cond={…}>`. Ad-hoc JSX construction in arbitrary Dang expressions
  remains unsupported; left for if-and-when it matters.

- *PascalCase Partial keys remain in `docs/html/columns.tmpl`*
  (`.Partial "Columns"`). That template is still alive because
  `<Columns>` is a Go built-in that emits `Styled{Style: "columns",
  Partials: {"Columns": ...}}`. Harmless; cleanup paired with any
  future columns refactor.

- *Reference-display regression flagged in the previous entry is
  resolved.* Rendered docs are byte-identical to the pre-port
  baseline.
