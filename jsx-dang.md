# Booklit pivot: JSX syntax + Dang embedded language + Dagger plugins

A hard cutover from Go-native Booklit to a Dang-native, Dagger-native
documentation system. No deprecation path; the existing `\foo{}` syntax and
Go plugin loader go away.

## Why

Three pressures lined up:

1. **The plugin recompilation dance is awkward.** Today `--plugin` causes
   Booklit to write a `main.go` that imports the plugin packages, runs
   `go install`, and re-execs itself with `BOOKLIT_REEXEC=1`. See
   `booklitcmd/command.go` lines 63–74 and 215–265. It works, but it
   makes plugins a heavyweight commitment and locks the ecosystem to
   Go.

2. **Most "plugins" exist purely to bridge AST invocations to HTML
   templates.** The pattern is everywhere:

   ```go
   func (p Plugin) Foo(content booklit.Content) booklit.Content {
     return booklit.Styled{Style: "foo", Content: content, Partials: ...}
   }
   ```

   …paired with `html/foo.tmpl`. The Go code adds nothing the framework
   can't infer. If invocations defaulted to looking up a template by
   name, ~90% of plugin code would evaporate.

3. **JSX gives `\foo{a}{b}{c}` a native named-prop syntax.** Booklit
   currently fakes named args via `Partials`. JSX makes them
   first-class:

   ```jsx
   <Card title="Hello" icon="star">body</Card>
   ```

   …directly mirrors `Styled{Style: "Card", Content: body, Partials:
   {Title: ..., Icon: ...}}`.

The pivot also dogfoods Dang: by embedding Dang as the expression
language inside `{}` interpolations, every Booklit-built site becomes
a live consumer of the Dang interpreter. That's a feedback loop on
Dang's ergonomics that no test suite can replicate.

## Target shape

A source file under `lit/foo.md`:

```mdx
# Types {#types}

The standard library exposes a few primitive types:

<DefinitionList>
  {primitiveTypes.map(t => <Definition term={t.name}>{t.docs}</Definition>)}
</DefinitionList>

<Aside>
  See <Reference tag="nullability"/> for how `!` works.
</Aside>
```

A presentation file under `html/Card.html` (or `html/DefinitionList.html`):

```html
<div class="card">
  <h3>{{.title}}</h3>
  <div class="body">{{.children | render}}</div>
</div>
```

A Dang module under `dang/` (or referenced via a Dagger module URL) that
exports `primitiveTypes`:

```dang
pub primitiveTypes = [
  {{ name: "Int!",    docs: "64-bit signed integer" }}
  {{ name: "String!", docs: "UTF-8 text" }}
  ...
]
```

Build: `booklit -i lit/index.md -o dist/`. No Go compile step, no
`--plugin` flag. Plugins resolve in order: (a) template in `html/`,
(b) Dang function in scope, (c) Dagger function via `<Foo
from="github.com/.../module">`.

## What changes vs. what stays

| Stays | Changes | Goes away |
|---|---|---|
| Section tree model | Invocation syntax (JSX) | `\foo{}` syntax |
| `\split-sections` / `\single-page` semantics | Expression language (Dang) | Go plugin registry |
| `Styled` + `Partials` content model | Plugin resolution (template-first) | `--plugin` flag + reexec |
| HTML template renderer | Embedded language eval (Dang) | `RegisterPlugin` API |
| Markdown for prose (goldmark) | Invocation parsing path | `baselit` as a Go plugin |
| `\include-section` (re-spelled `<IncludeSection>`) | Per-file plugin scoping | The `BOOKLIT_REEXEC` dance |

The baselit functions (`Title`, `Section`, `IncludeSection`,
`SplitSections`, `UsePlugin`, `Aux`, `Styled`, `SetPartial`,
`TableOfContents`, etc.) become **built-in JSX components** rendered by
the evaluator directly, not user-replaceable plugins. They're the
language, not the standard library.

## Architecture sketch

```
lit/*.md                Markdown + JSX content
html/*.html             Per-component templates (Go html/template)
dang/*.dang             Dang modules (optional; for content data + helpers)
booklit.toml            Project config: Dagger modules, paths, output settings
dist/*.html             Build output
```

Pipeline:

```
parse (goldmark + JSX inline parser)
  → AST (Section tree with JSXNode invocations)
  → resolve (built-in components, then template lookup, then Dang, then Dagger)
  → evaluate (run Dang for {expr}; invoke Dagger for <X from="..."/>)
  → render (existing HTML engine; Styled+Partials still maps to .tmpl)
```

## File-level impact

### Deletions

- `booklitcmd/command.go` lines 63–74, 215–265 — reexec apparatus
- `plugin.go` — `RegisterPlugin` / `LookupPlugin` / `PluginFactory`
- `baselit/` as a separate package — folded into a `builtins/` directory
  of JSX components
- `chroma/` plugin — re-implement as a built-in `<Code>` component;
  Chroma is just a Go dep, not a plugin
- `docs/hello/`, `tests/fixtures/erroring-plugin/`, etc. — example
  Go plugins, replaced by Dagger module examples

### Major rewrites

- `marklit/invoke_parser.go` (143 lines) — replace with JSX inline
  parser. Trigger on `<` followed by uppercase letter (component) or
  `</` (close tag); fall through to goldmark otherwise.
- `marklit/invoke_node.go`, `convert.go` — adjust to emit a new
  `ast.JSXElement` node instead of `ast.Invoke`.
- `ast/node.go` — replace `Invoke` with `JSXElement{Name, Props,
  Children}`. `Method()` and the dash-to-CamelCase logic go away (JSX
  names are already PascalCase).
- `stages/evaluate.go` — replace plugin-method reflection with the
  4-tier dispatch (built-in → template → Dang → Dagger).
- `booklitcmd/command.go` — strip plugin flags, add Dagger client
  bootstrap.

### New code

- `jsxparse/` — JSX inline parser (PEG via pigeon, matching dang's tooling
  choice). Handles attribute values that are strings, `{dangExpr}`,
  boolean shorthand, spread, self-closing, children.
- `dangeval/` — wrapper around the Dang interpreter. Loads `dang/`
  modules, exposes scopes per-section so `{expr}` can reference
  in-scope names. Bridges Dang values ↔ `booklit.Content`.
- `daggerdispatch/` — Dagger client wrapper. Resolves `<X
  from="github.com/.../mod">` to a module function call. Caches
  responses across a single build session.
- `builtins/` — `<Title>`, `<Section>`, `<IncludeSection>`,
  `<SplitSections>`, `<SinglePage>`, `<TableOfContents>`,
  `<Reference>`, `<Aux>`, `<Code>`, `<Inset>`, `<Aside>`, etc. Each is
  a Go function that returns `booklit.Content`. Built-in lookup is the
  first dispatch tier.
- `booklit.toml` schema — `[modules]` mapping local names to Dagger
  module refs; `[paths]` for `lit/`, `html/`, `dang/`, output dir;
  `[render]` for engine-level settings.

## Phases

Each phase ends with something runnable. Don't merge a phase that
breaks the build.

### Phase 0: clear the decks

- Move `baselit/` functionality inline (we'll re-express it as
  built-ins in Phase 2; for now keep it as a package, just stop
  referring to it as a "plugin").
- Delete the Go plugin loader, `--plugin` flag, and reexec code. Sites
  that referenced custom Go plugins will stop building. We're a hard
  cutover.
- Delete example Go plugins (`docs/hello/`, the test-fixture plugins).
- All existing tests that depend on `--plugin` get deleted or rewritten
  in later phases.

Outcome: Booklit builds and renders, but only with the built-in
functions. No way for users to extend it yet.

### Phase 1: JSX parser

- Write `jsxparse/` as a pigeon PEG grammar. Reuse patterns from
  Dang's grammar tooling (`pkg/dang/dang.peg` → `dang.peg.go`).
- Register a goldmark inline parser that triggers on `<` (uppercase
  follow) and `</`. Block-level JSX (`<Foo>...</Foo>` as its own
  paragraph) needs a block parser too.
- Emit `ast.JSXElement{Name, Props, Children}`. Props are
  `map[string]ast.Node` (string literal, expression `{...}`, or
  boolean shorthand `disabled`).
- Convert `ast.Invoke` references throughout the codebase to
  `ast.JSXElement`. The existing `\function-name` →
  `FunctionName` method-name logic moves into JSX (where the name is
  already `FunctionName`, so it's a no-op).
- Update `marklit/convert.go`'s thematic-break and heading conversions
  to emit JSX (e.g., `<ThematicBreak/>` instead of
  `Invoke{Function: "thematic-break"}`).
- Update all existing `lit/` content (Booklit's own docs) to use JSX.
  This is mostly mechanical: `\foo{x}{y}` → `<Foo prop="x">y</Foo>`
  where the prop name conventions need to be figured out (see open
  questions).

Outcome: Booklit's own docs build using JSX syntax instead of `\foo{}`.
No new capabilities yet.

### Phase 2: template-default dispatch

- Implement the resolution order:
  1. **Built-in**: look up `Name` in `builtins/` registry. Built-ins
     are Go functions with typed signatures, matching today's plugin
     methods. They handle the section-tree fundamentals
     (`<Section>`, `<IncludeSection>`, etc.).
  2. **Template**: look for `html/Name.html`. If present, render with
     props as `.propName` and children as `.children`. This is the
     "unknown component" default — equivalent to `Styled{Style: "Name",
     Content: children, Partials: props}`.
  3. (defer Dang and Dagger to phase 3 and 4)
- Re-implement chroma as a `<Code>` built-in. Chroma stays a Go
  dep (we're not pretending to evict Go), it's just not a "plugin"
  anymore.
- Delete `Styled` if no one is using it directly anymore — or keep it
  as the internal data type that templates render.

Outcome: writing a new "plugin" means dropping an `.html` template
into `html/`. No Go code, no recompile.

### Phase 3: embed Dang as the expression language

This is the dependency-flip moment. Booklit currently has no
relationship to Dang. After this phase, Booklit imports Dang.

- Add `github.com/vito/dang` as a Go module dependency.
- Build `dangeval/` to wrap the Dang interpreter. Each section gets a
  Dang environment that:
  - has all in-scope props/locals as Dang values
  - exposes Booklit's section tree for queries (e.g., `sections()` to
    enumerate)
  - can evaluate a Dang expression in that environment, returning a
    Dang value
- Bridge Dang values ↔ `booklit.Content`:
  - `String` → `String`
  - `Int!`, `Float!`, `Boolean!` → stringified `String`
  - `[T]!` → `Sequence`
  - records (`{{}}`) → keyed `Partials` if used as a JSX child? Or
    error? (open question)
  - JSX element values? Dang doesn't have those — Dang functions can't
    construct `<Foo/>` directly. (TBD: should they be able to, via
    something like a `BooklitContent` opaque type?)
- Wire JSX attribute `attr={expr}` and child `{expr}` to dangeval.
- Allow `dang/` directory at the project root: top-level Dang module
  whose `pub` bindings are available in every section.
- Provide `<UseDang src="./helpers.dang"/>` for per-file imports if
  needed.

Outcome: docs can reference dynamic data:

```mdx
<DefinitionList>
  {primitiveTypes.map(t => <Definition term={t.name}>{t.docs}</Definition>)}
</DefinitionList>
```

…with `primitiveTypes` defined in a `dang/` module.

### Phase 4: Dagger dispatch

- Add `github.com/dagger/dagger/sdk/go` (or the right import path).
- `booklit.toml` `[modules]` section maps local names to Dagger module
  refs:
  ```toml
  [modules]
  diagrams = "github.com/vito/booklit-diagrams"
  ```
- `<Diagrams.Mermaid src="..."/>` resolves to calling the `mermaid`
  function in the `diagrams` module via Dagger.
- Anonymous form: `<Foo from="github.com/.../mod"/>` for one-off use.
- The whole build runs in a single Dagger session (boot once, reuse
  for every call).
- Per-call Dagger overhead is small (function dispatch within an
  already-warm session); the cost we're avoiding is the per-`booklit`-
  invocation Go compile + reexec.
- Caching: Dagger memoizes function calls automatically — no plugin
  cache needed.

Outcome: ecosystem opens up. Plugins in Python, Rust, JS, anything
with a Dagger SDK. The Booklit core stays small.

### Phase 5: dogfood

- Rewrite Booklit's own `docs/` site in the new system.
- Move dang docs to use this (this is the application that started the
  whole conversation).
- Write a `booklit-init` command (or document the layout) so new
  projects can be scaffolded.

## Open questions to resolve before / during Phase 1

These need answers before the JSX parser commits to a shape.

1. **Attribute name convention.** Today's `Styled.Partials` uses
   PascalCase keys (`"Title"`, `"Src"`). React/JSX convention is
   camelCase (`title`, `src`). Pick one. Recommendation: camelCase, since
   it matches what people expect from `<Foo title="x">`. Templates would
   reference `{{.title}}` instead of `{{.Title}}`. This means tuist-
   style templates (`{{.Partial "Src"}}`) need a different access pattern,
   probably `{{.props.src}}` or just `{{.src}}` if we flatten.

2. **How are children passed to templates?** Options:
   - `{{.children | render}}` — single content
   - `{{range .children}}{{. | render}}{{end}}` — list
   - Both, depending on whether the JSX had multiple children?
   Recommendation: always a list. Single-child templates just iterate.

3. **`{dangExpr}` inside attributes vs. as a child.** JSX allows both:
   `<Foo bar={x}>{y}</Foo>`. In the attribute case, the value is the
   Dang expression's result. In the child case, the result is content
   that interleaves with siblings. Both need handling.

4. **Conditional and loop rendering.** React idiom is `{cond && <X/>}`
   and `{arr.map(x => <X/>)}`. Dang supports `if (cond) { ... }` and
   `.map { x => ... }`. Both work as expressions, but Dang's syntax is
   different from JS's. Pick: do we limit `{...}` to value expressions
   that return content/strings, or full Dang blocks? Recommendation:
   full Dang. The expression is parsed by the Dang parser; it returns a
   value; we coerce to content.

5. **Tag-name conflicts with HTML.** `<table>`, `<div>`, etc. are real
   HTML. In React, lowercase = HTML, uppercase = component. We can adopt
   that: lowercase JSX names pass through to literal HTML. This means
   markdown `<div>` blocks still work.

6. **Section creation.** Today, headings (`# Foo`) auto-create sections.
   Do we keep that, or require explicit `<Section title="Foo">...</Section>`?
   Recommendation: keep headings. They feed `<Section>` invocations under
   the hood. Authors who want non-heading sections can use `<Section>`
   explicitly.

7. **What's a "page boundary"?** Today, `split-sections` on a parent
   means each child renders as a separate file. In JSX terms, this
   means `<Section title="..."><Section title="..."/></Section>` with
   `<SplitSections/>` on the outer one. Verify the data model still
   carries this through cleanly.

8. **Backwards compat for built-in `\foo{}` content.** None. The pivot
   is hard. Booklit's own `docs/lit/*.md` files get rewritten as part
   of Phase 1.

## Open questions for Phase 3 (Dang embedding)

1. **Where does Dang live?** Booklit imports `github.com/vito/dang`
   directly as a Go package. No new embedding API needed — the public
   surface that Dang's own CLI/REPL composes against is enough.
   Reference uses: `dang.RunFile`, `dang.RunDir`, `dang.ServiceRegistry`,
   `dang.ContextWithServices`, `dang.FindDaggerModule`,
   `dang.ClearSchemaCache`, `dang.FormatFile` (see
   `cmd/dang/main.go`). Booklit composes the same primitives to
   evaluate per-section expressions; the REPL is the working analogue
   to "evaluate a snippet in an environment that holds state across
   calls."

2. **What does Dang see as the "schema"?** Dang's current model assumes
   a GraphQL schema gives it root types and functions. When embedded in
   Booklit, what's the schema? Options:
   - No schema: Dang runs in "schema-less" mode where it only knows
     primitives, lists, records, and user-defined types.
   - A Booklit-provided schema: types like `Section`, `Content`,
     `Reference` are available as built-in Dang types.
   Recommendation: schema-less for v1, add a Booklit schema in a later
   phase if it proves useful.

3. **Threading.** Dang's evaluator is presumably single-threaded per
   environment. Booklit currently evaluates sections somewhat in
   parallel (verify). We may need per-section Dang environments.

## Open questions for Phase 4 (Dagger)

1. **Session lifecycle.** Boot the Dagger session at `booklit` startup,
   tear down at exit. The `-s` server mode needs a long-lived session.
   Verify Dagger's Go SDK supports this cleanly.

2. **Argument marshalling.** JSX attributes are strings or Dang values.
   Dagger function calls have typed arguments. We need a mapping. For
   string args this is trivial; for complex args (records, lists),
   maybe limit to JSON-encodable in v1.

3. **Error reporting.** Dagger function errors need to surface in
   Booklit's existing error-location framework (file path, line/column).
   The Dagger SDK should give us at least the function-name context;
   we attach the source location of the `<Foo from="..."/>` invocation.

4. **Offline / no-Dagger mode.** Should Booklit still work without
   Dagger if a project doesn't use `from=` plugins? Probably yes —
   only boot the session lazily on first Dagger dispatch.

## Risks

- **JSX-in-markdown parsing edge cases.** Goldmark's parser is
  flexible but JSX has corners (raw HTML coexistence, whitespace
  significance, fragments, multi-line attributes). Budget for a long
  tail of test fixtures.
- **Dang's evaluator wasn't designed for embedding from another
  long-running program.** The REPL is the closest existing consumer
  and may not exercise every seam Booklit needs (e.g., evaluating an
  expression in a frozen environment, or feeding Dang values back as
  arguments to subsequent calls). Expect to discover ergonomic gaps
  in Dang's public API during Phase 3 and fix them as found.
- **Loss of Go-plugin authors.** Anyone who wrote a Go plugin against
  the old API needs to port to either a template, a Dang module, or a
  Dagger module. We control all of Booklit's own plugins; external
  plugins are rare enough to be a small concern.
- **Performance regression.** Today's Go-plugin reflection dispatch is
  microseconds. Dang interpretation + Dagger calls are milliseconds.
  For docs with thousands of invocations, build times could grow. Need
  benchmarks at the end of Phase 3 and Phase 4.

## Relevant existing code

Key files to read before starting Phase 1:

- `ast/node.go` — current `Invoke` shape; this is what becomes
  `JSXElement`.
- `marklit/invoke_parser.go` (143 lines) — current `\foo{}` parser;
  the JSX parser replaces it.
- `marklit/marklit.go` (170 lines) — how the goldmark extension is
  wired; the JSX parser plugs in here.
- `marklit/convert.go` (772 lines) — converts goldmark AST to Booklit
  AST. Heading→Section logic and thematic-break→Invoke conversion
  live here.
- `stages/evaluate.go` — `VisitInvoke` dispatches to plugin methods
  via reflection. This becomes the 4-tier resolver.
- `booklitcmd/command.go` — CLI entry point; the reexec dance and
  `--plugin` flag get torn out.
- `baselit/plugin.go` — every function here becomes a built-in
  component in `builtins/`. This is the inventory of language
  fundamentals.
- `render/html.go`, `render/writer.go` — HTML rendering; mostly
  unchanged.

Dang code to read for Phase 3 (lives at `~/src/dang`, module
`github.com/vito/dang`):

- `cmd/dang/main.go` — the script-mode entry points (`RunFile`,
  `RunDir`) and how services/context get wired in. Booklit composes
  the same primitives.
- `cmd/dang/repl_tuist.go` — REPL eval loop; the closest existing
  analogue to "evaluate a snippet in a held-open environment."
- `pkg/dang/dang.peg` — grammar; reference for pigeon usage.
- `pkg/dang/eval.go` — evaluator entry points.
- `pkg/dang/env.go` — module/environment model; this is what we
  populate per-section.
- `pkg/dang/types.go` — value types we bridge to `booklit.Content`.

## Definition of done

- Booklit's own docs site (`booklit/docs/lit/*.md` → `dist/*.html`)
  builds with zero Go plugins, zero `--plugin` flags, zero reexec.
- All current behavior preserved: section tree, split-sections,
  table-of-contents, references, includes, styled content, partials.
- New affordances unlocked: template-only "plugins," Dang for data
  and control flow, Dagger for cross-language extension.
- A new `<Foo>` component is added by either (a) dropping
  `html/Foo.html`, (b) defining `pub Foo(...)` in `dang/`, or
  (c) adding a Dagger module to `booklit.toml`. None require touching
  Booklit's source.
- The Dang docs site is rebuilt on the new Booklit. This is the
  acceptance test from the application side.

## Out of scope

- Search (Pagefind-style). Plugable later; not a v1 concern.
- Live reload / `-s` server mode upgrades. Keep what works.
- Image optimization, asset pipelines. Not Booklit's job.
- Multi-language sites (i18n). Not yet.
- Versioned docs. Not yet.

## What to do first (concrete starting point)

1. Read `marklit/invoke_parser.go`, `marklit/marklit.go`, and
   `ast/node.go` end to end. Form a clear mental model of the existing
   parse pipeline.
2. Stub out `ast.JSXElement` alongside `ast.Invoke`. Don't delete
   `Invoke` yet — let them coexist while the JSX parser stabilizes.
3. Write a minimal JSX inline parser that handles `<Foo/>`,
   `<Foo>body</Foo>`, and `<Foo bar="x"/>`. Register it as a goldmark
   inline parser.
4. Write a handful of `.md` test fixtures and golden outputs.
5. Once the parser is solid, swap `\foo{}` to JSX in
   `booklit/docs/lit/index.md` as a forcing function.
6. Then start ripping out the plugin system.

## Progress log

This section is append-only. Each session that makes meaningful progress
should add a new dated `###` entry at the bottom — do not edit earlier
entries or earlier sections of this plan, even when later decisions
contradict them. A reader scanning top-to-bottom learns the original
intent first, then sees how the work actually unfolded.

### 2026-05-30 — Phase 0 + start of Phase 1

Landed in commit `b748c94` (`feat(marklit): add JSX parsing alongside
\invoke`).

**Done.** `ast.JSXElement` and `ast.JSXExpression` now live alongside
`ast.Invoke` as parallel node kinds. The goldmark layer has a
hand-written JSX inline parser (priority 98) and block parser
(priority 100, ahead of goldmark's HTML block parser at 900, which
would otherwise claim `<UpperCase` lines as CommonMark type-7 HTML and
block inline JSX parsing). The parser handles `<Foo/>`, `<Foo>body</Foo>`,
`<Foo bar="x"/>`, `<Foo bar={expr}/>`, `<Foo>{expr}</Foo>`, nested
elements, multi-line attrs and bodies, escapes, and string-aware brace
balancing inside expressions. camelCase prop names are preserved
verbatim; lowercase tags fall through to goldmark's raw-HTML support.
Markdown inside JSX bodies still works (`<Foo>*emph*</Foo>`) — text
chunks are re-parsed via `ParseInlineArg`. 17 positive table-driven
cases plus 3 negative tests in `marklit/jsx_test.go`; full
`go test ./...` green.

**Decisions made.** Walked through the eight open Phase 1 questions
explicitly: camelCase attrs, lowercase = HTML pass-through, full Dang
in `{...}` (parsed opaquely for now — captured as raw bytes with
string-aware brace balancing, real parsing comes in Phase 3), headings
continue to auto-create `<Section>` invocations, and no backcompat for
`\foo{}`. Q2 (children-to-template format), Q4 (conditional/loop
semantics), and Q7 (split-sections data model) deferred until their
phase actually hits.

**Plugin system status.** Still in place. The `\foo{}` syntax, the Go
plugin loader, `--plugin`, `BOOKLIT_REEXEC`, and the example plugins
under `docs/hello/` and `tests/fixtures/*-plugin/` all still build and
work. Both syntaxes coexist deliberately so file-by-file migration is
possible. Deletion was deferred until "JSX parses real fixtures," which
it now does — so the next session can remove the plugin apparatus in
a single commit if appropriate.

**JSX evaluation is not wired up.** `stages.Evaluate.VisitJSXElement`
and `VisitJSXExpression` return a "not yet implemented" error.
Dispatching a JSX element to a template, built-in, Dang function, or
Dagger module is Phase 2/3/4 territory. A real `.md` file using JSX
will parse but not render until those land.

**Reasonable next steps.** (a) Delete the plugin loader / `--plugin` /
`BOOKLIT_REEXEC` / example Go plugins now that JSX parses real
fixtures — pure cleanup commit, no behavior change for the JSX path.
(b) Phase 2 (template-default dispatch): make a 4-tier resolver in
`stages/evaluate.go` so `<Foo>` first checks built-ins, then a
`html/Foo.html` template, before falling through to errors. (c) Phase 1's
docs migration: rewrite `booklit/docs/lit/*.md` from `\foo{}` to JSX as
a forcing function — this probably waits for Phase 2 since otherwise
the migrated docs would parse but not render.
