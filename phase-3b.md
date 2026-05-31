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

(empty)
