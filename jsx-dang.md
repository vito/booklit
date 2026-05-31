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

### 2026-05-30 — Phase 0 cleanup complete

Landed in commit `cbb4d68` (`refactor: remove the Go plugin loader and
\use-plugin`).

**Done.** Removed the `--plugin` CLI flag, the `BOOKLIT_REEXEC` /
`reexec()` apparatus that wrote a generated `main.go` and re-execed
under `go install`, the `booklit.RegisterPlugin` / `LookupPlugin`
registry, the `\use-plugin` function on baselit, and every example Go
plugin: `chroma/`, `docs/hello/`, `docs/go/` (booklitdoc), and the test
fixture plugins under `tests/fixtures/*-plugin/`. Integration test
cases that depended on those fixtures went with them.

**What stays.** `booklit.Plugin` and `booklit.PluginFactory` types
remain because `Section.PluginFactories` is still the evaluator's
dispatch mechanism for `\foo{}` until the 4-tier resolver lands.
baselit stays as a hardcoded base in `booklitcmd.basePluginFactories`,
and its chroma-driven `\code-block` / `\syntax` functions are
untouched, so basic syntax highlighting still works for the existing
`\foo{}` path.

**Known breakage.** Booklit's own `docs/lit/*.md` files reference
`\use-plugin{booklitdoc}` and `\use-plugin{chroma}` at the top of every
file. Those will fail with `undefined function \use-plugin` until the
docs are rewritten to JSX in the Phase 1 docs migration. The test
suite does not exercise those docs and stays green.

**Next.** Phase 2 (template-default dispatch) is the obvious next
substantive step — it's what makes JSX render. Once `<Foo>` can look
up `html/Foo.html` and fall back to built-ins, we have a usable system
end-to-end and can attempt the docs migration.

### 2026-05-30 — Phase 2: dispatch + initial built-ins

Landed in commit `373a6b9` (`feat: dispatch JSX via builtins registry
and template fallback`).

**Done.** JSX elements render end-to-end. The evaluator
(`stages/evaluate.go::VisitJSXElement`) consults a new `builtins/`
package first, then falls back to wrapping the element in a
`booklit.Styled` so dropping `html/<Name>.tmpl` is enough to introduce
a new component. Built-ins receive raw `ast.Node` props and children
plus a `Context.Evaluate` helper, so they can choose whether to
evaluate eagerly or pass raw AST through (used by `<Section>`, which
recurses via `Processor.EvaluateNode`).

The initial built-in set: `<Title>`, `<Section>`, `<Reference>`,
`<Target>`, and the styled-content family (`<Italic>`, `<Bold>`,
`<Larger>`, `<Smaller>`, `<Strike>`, `<Superscript>`, `<Subscript>`,
`<Inset>`, `<Aside>`). Five integration tests in `tests/jsx_test.go`
exercise the dispatcher across builtins, recursion, refs/targets, and
template fallback. Full `go test ./...` green.

**Important refinement.** The parser now distinguishes single-line
from multi-line JSX. `<Title>x</Title>` on a single line parses its
children as inline (matching `\title{x}`); a JSX element whose tags
straddle multiple lines parses children as block content (matching
`\section{...}` with multi-line braces). Multi-line state is tracked
per element via a `crossedLine` flag in the scanner that's saved and
restored around each `parseJSXElement` call, so an inline `<Title>`
nested inside a multi-line `<Section>` keeps inline semantics. This
mirrors the existing single-line/multi-line invoke distinction baked
into the marklit preprocessor.

**Template fallback details.** The fallback wraps an unknown component
in `booklit.Styled{Style: Name, Content: <evaluated children>,
Partials: <evaluated props>}`. Prop keys are kept as authored
(camelCase), so a template uses `{{.Partial "title"}}`. The renderer
looks up `<Name>.tmpl` exactly as written — case-sensitive, no
slug-style conversion. This is a pragmatic reuse of the existing
template machinery; if camelCase-pure templates ever become awkward
the boundary can be revisited without changing user-visible JSX.

**Notable trade-off.** Single-line block JSX with text chunks split
across child elements (`<Foo>Use <Bar/>.</Foo>` if multi-line) wraps
each text chunk in its own paragraph, because chunks are
independently `ParseArg`'d. Acceptable for MVP; the cleaner fix is to
re-serialize children with placeholders for nested elements, then
block-parse the whole body. Deferred.

**Long tail still missing.** Many `\foo{}` invocations don't have JSX
built-in equivalents yet: `<Code>` / `<CodeBlock>` (chroma syntax
highlighting), `<List>` / `<OrderedList>` / `<Definition>` /
`<Definitions>`, `<Image>`, `<Link>`, `<Table>` / `<TableRow>`,
`<TableOfContents>`, `<IncludeSection>`, `<SinglePage>`,
`<SplitSections>`, `<SetPartial>`, `<Aux>`, `<ThematicBreak>`. Each is
a small forwarding wrapper around existing baselit code; adding them
is a mechanical follow-up.

**Next.** Two reasonable directions. (a) Fill in the long tail of
built-ins so JSX matches `\foo{}` coverage; this unblocks the docs
migration. (b) Phase 3 — embed Dang as the expression language, so
`{primitiveTypes.map(...)}` becomes evaluable. (a) is mostly
mechanical; (b) is the bigger unknown and likely surfaces gaps in
Dang's public API.

### 2026-05-30 — Phase 1 docs migration complete

Landed in commits `77be16e`, `b1e7bf9`, `3a575eb`, `5c15f53`.

**Done.** Booklit's own docs site builds end-to-end on the new JSX
pipeline. Concretely:

- `feat(builtins): add long-tail baselit built-ins` (`77be16e`):
  ~30 new components (Aux, Code, Link, Image, ThematicBreak, Styled,
  SetPartial, IncludeSection, SinglePage, SplitSections,
  OmitChildrenFromTableOfContents, TableOfContents, CodeBlock, Syntax,
  List, OrderedList, Item, Definition, Definitions, Table, Row,
  TableRow). Containers use child-element-based shape (`<List>` takes
  `<Item>`, `<Table>` takes `<Row>` of `<Item>`s, etc.).

- `fix(marklit): support single-quoted attrs and backtick code spans`
  (`b1e7bf9`): two parser corner-cases discovered during migration —
  single-quoted attribute values (so a `sig` prop can carry a JSX
  signature with `"` inside), and a backtick-code-span flag in
  `parseChildren` so literal `` `attr={expr}` `` in docs doesn't get
  interpreted as a Dang expression.

- `feat(docs): restore booklitdoc helpers as docs-specific built-ins`
  (`3a575eb`): the booklitdoc helpers (<Define>, <Columns>,
  <OutputFrame>, <Godoc>, <LitSyntax>, <TemplateLink>, <SyntaxHl>,
  <ColumnHeader>, <Column>) live in `docs/booklitdoc/`, registered via
  `builtins.Register`. A new `cmd/booklit-docs` binary imports them on
  top of the standard `booklitcmd`, keeping docs-specific styling out
  of the main `booklit` binary. The chroma `styles.Fallback` override
  for the booklitdoc palette also moved here.

- `docs: migrate docs/lit to JSX` (`5c15f53`): all eight doc files
  rewritten — `\use-plugin{...}` removed, `\foo{...}` rewritten as
  `<Foo>...</Foo>`, prose defaults to Markdown where it can (headings,
  lists, fenced code blocks, `[#tag]` references). `<Define>` API
  changed from "consume the raw `\foo{...}` AST" to "take a string
  `sig` prop"; the documented signatures are now JSX. Makefile target
  switched to `go run ./cmd/booklit-docs`.

The full docs site (`docs/lit/index.md` → `docs/outputs/`) now builds
with zero errors, zero Go plugins, zero `--plugin` flags. Full
`go test ./...` still green.

**Decisions worth knowing.**

- *Container shape*. `<List>`/`<Table>`/`<Definitions>` use explicit
  sub-components (`<Item>`, `<Row>`, `<Definition>`) rather than
  positional children. Whitespace text between entries is ignored.
  Markdown native lists and tables still work and remain the
  recommended form.

- *Prop key case*. Camel-case end to end. Templates that need a prop
  reach for it by name: `{{.Partial "title"}}`.

- *Single-line vs multi-line children*. Carried over from Phase 2:
  `<Title>x</Title>` on one line parses children as inline (matching
  `\title{x}`); a multi-line `<Section>` parses children as block.
  Nested elements use their own line-span, so an inline `<Title>`
  inside a multi-line `<Section>` keeps inline semantics.

- *Backtick code spans*. Inside backticks, `<` and `{` are literal.
  Multi-character fenced spans (`` ``...`` ``) aren't tracked; only
  single-backtick spans. Sufficient for the docs.

**Known follow-ups.**

- A multi-paragraph block JSX whose children mix nested elements with
  text wraps each text chunk in its own paragraph. The cleaner fix is
  to re-serialize children with placeholders for nested elements and
  block-parse the whole body. Deferred — the docs don't trip on this.

- Templates still use PascalCase partial keys in a few places
  (`{{.Partial "Title"}}`, `{{.Partial "URL"}}`) because the docs
  templates predate the camelCase decision. Those are content for
  cleanup but harmless: they're string lookups, not interface-level
  bindings.

**Next.** Phase 3 — embed Dang as the expression language. Lots of the
deferred-questions list (full `{expr}` evaluation, conditional/loop
rendering, dynamic data via `{primitiveTypes.map(...)}`) opens up once
Dang is wired in. Expect to discover ergonomic gaps in Dang's public
embedding API; the REPL in `cmd/dang/repl_tuist.go` is the closest
existing analogue.

### 2026-05-30 — Phase 1 polish (post-migration)

Two small commits landed after the Phase 1 entry above, both in the
Markdown-over-JSX direction.

- `d8b13af` (`refactor(builtins): drop List/Table built-ins in favor of
  Markdown`): removed `<List>` / `<Item>` / `<Table>` / `<Row>` and
  their wrapper content types. Markdown's `- item`, `1. item`, and
  `| a | b |` already cover these and read better. `<Definitions>` /
  `<Definition>` stay because CommonMark has no definition-list form.
  `docs/lit/baselit.md` drops the JSX entries and points at the
  Markdown forms; `getting-started.md`, `syntax.md`, `plugins.md` use
  Markdown lists directly.

- `6333475` (`fix: pass raw HTML through to the renderer instead of
  escaping`): `marklit/convert.go` was emitting `ast.String(text)` for
  inline raw HTML and HTML blocks, so any `<dl>` / `<em>` etc. in a
  Markdown doc came out as `&lt;dl&gt;`. Now those routes emit
  `\raw-html` / `\raw-html-block` invokes that wrap as
  `booklit.Styled{Style: "raw-html"}` (with `Block=true` for blocks);
  the existing `raw-html.tmpl` (`{{. | rawHTML}}`) does the pass-
  through. Method names are `RawHtml` / `RawHtmlBlock`, not
  `RawHTML` — invokes resolve via dash-to-CamelCase
  (`raw-html` → `RawHtml`), which doesn't align with Go's initialism
  convention. Markdown formatting inside a raw HTML block still isn't
  re-parsed, matching standard CommonMark behavior. Two regression
  tests in `tests/jsx_test.go`.

### 2026-05-30 — Phase 3 MVP: Dang expressions in JSX

Embedded the Dang interpreter so `{expr}` interpolations inside JSX
evaluate end-to-end. Booklit now imports `github.com/vito/dang` as a
Go dependency.

**Done.**

- New package `dangeval/`:
  - `dangeval.Evaluator` holds a long-lived type env (`hm.Env`) +
    value env (`dang.EvalEnv`) + service registry. One Evaluator per
    build session.
  - `New(ctx, projectDir)` discovers `dang.toml` via
    `dang.FindProjectConfig` (walks upward), resolves imports via
    `dang.ResolveImportConfigs`, auto-detects `dagger.json` via
    `dang.ResolveDaggerImport`, then bootstraps via
    `dang.BuildEnvFromImports`. Without any config the result is a
    Prelude-only env, which still handles primitive expressions.
  - `Eval(raw)` mirrors the REPL's `startEval` loop
    (`cmd/dang/repl_tuist.go`): `ParseWithRecovery` → take
    `*ModuleBlock.Forms` → `InferFormsWithPhases` against the held
    type env → `EvalNode` per form against the held eval env.
  - `Close()` calls `services.StopAll()` for any subprocesses started
    by dang.toml imports (e.g. `dagger session`).

- `dangeval/bridge.go` maps Dang values to `booklit.Content`:
  `StringValue` → `booklit.String`, scalars stringified, `ListValue`
  → `booklit.Sequence` of bridged elements, `NullValue` → empty.
  Records / functions / modules are an error for v1.

- `stages.Evaluate` grew a `Dang *dangeval.Evaluator` field.
  `VisitJSXExpression` now calls `Dang.Eval(node.Raw)` and appends
  the bridged content. `evalArg` propagates the field to sub-
  evaluators so JSX prop expressions and child expressions go through
  the same path.

- `load.Processor` grew a matching `Dang` field;
  `evaluateSection` plumbs it into the `stages.Evaluate`.

- `booklitcmd.Command.Build` constructs an Evaluator rooted at the
  input file's directory and defers `Close()`. `Serve` is unchanged
  for now — server-mode lifecycle is its own design call.

- `tests.Example.Run` constructs an Evaluator per-test rooted at the
  temp dir, so every test exercises the Dang path. Seven new cases in
  `tests/dang_test.go` cover int, string, arithmetic, boolean, list,
  expression-in-prop (via the existing `Card.tmpl` fallback), and a
  parse-error case. Full `go test ./...` green; the docs site
  (`go run ./cmd/booklit-docs -i docs/lit/index.md -o ... --html-templates docs/html`)
  builds clean.

**Decisions worth knowing.**

- *Discovery: dang.toml + dagger.json only, no `dang/` directory*.
  The original plan mentioned a `dang/` source-file directory at the
  project root. dang.toml is for GraphQL imports (Dagger sessions,
  SDL files, endpoints), not `.dang` source files. A directory-of-
  `.dang`-files convention is its own design call — defer.

- *Per-section scope: none yet*. Every snippet evaluates against the
  same global env. The `Section` type has no `Vars`/`Locals` field to
  hang per-section bindings off, and adding that interacts with
  built-ins' control flow (e.g. `<Section>` recursing via
  `Processor.EvaluateNode`). The simpler global model works for
  primitive expressions and for any data exposed by the project's
  Dagger module; richer scoping comes when we have a real use case.

- *Bridging is one-way*. Booklit consumes Dang values; Dang has no
  way to construct Booklit content (`<Foo/>` inside `{...}` doesn't
  work). The full `{items.map(i => <Item>{i.name}</Item>)}` story
  needs a `BooklitContent` opaque type in Dang. Deferred.

- *Threading is fine because Booklit is single-threaded*. Dang's
  `Env`/`EvalEnv` aren't safe for concurrent use (no mutexes on the
  module's `Values`/`Pending` maps). Booklit's evaluator processes
  sections sequentially, so a single shared Evaluator is correct.
  If parallelism is added later, we'll need per-goroutine envs.

- *Error surfacing is plain wrapping for now*. Errors include the
  raw expression text (`evaluating {1 +}: ...`) but don't map to
  source line/col within the document. Dang carries its own
  `SourceError` shape; integrating that into Booklit's
  `ErrorWithFile` framework is a follow-up.

- *Local replace directive*. `go.mod` has
  `replace github.com/vito/dang => /home/vito/src/dang` so dev-time
  changes flow through without a tag. When upstreaming, swap for a
  pinned version.

**Known limitations / follow-ups.**

- `{expr}` only works *inside* a JSX element (e.g. `<Foo>{1+2}</Foo>`
  or `<Foo bar={1+2}/>`). A bare `{1+2}` in a paragraph is literal
  text — the JSX inline parser triggers on `<`, not `{`. Adding a
  second trigger character means more brace-vs-literal-text disambig
  work; defer until the docs actually want it.

- No file/line context in expression errors yet (see decision above).

- `Serve` mode doesn't construct an Evaluator. Live-reload + a
  long-lived Dang env interact in non-obvious ways (re-loading
  dang.toml? Restarting Dagger session?); skipped for now.

- The `RawHtml` / `RawHtmlBlock` method-naming wart from the
  preceding polish entry still stands; orthogonal to Phase 3.

**Reasonable next steps.** (a) Decide on `dang/` source-file directory
support — when does a project want global `.dang` modules vs. just
pulling functions from a Dagger module via dang.toml? (b) Per-section
scope — what's the first use case that needs in-scope props/locals
inside `{...}`? (c) Source-mapping for Dang errors into Booklit's
file:line error framework. (d) JSX-inside-Dang
(`{items.map(i => <Item>{i.name}</Item>)}`) — needs design work on
the Dang side (`BooklitContent` opaque type) before Booklit can wire
it.

### 2026-05-31 — Phase 3a: tier-3 dispatch + .dang auto-discovery

Adds a third dispatch tier for JSX: when `<Foo>` isn't a built-in and
there's no `Foo.tmpl`, look up `Foo` in Dang scope and call it as a
component. Two patterns are supported:

- **Body-less** (`pub Foo(prop: T!): U! { … }`) — return value is
  coerced to content via the existing bridge. Useful for helpers that
  compute a string from props.
- **Body-ful** (`pub Foo(prop: T!, &body(name: U!): R!): R! { body(name: …) }`)
  — the JSX children compile into a Go-backed block; each `body(…)`
  call from inside the Dang function pushes the named args into Dang
  scope and re-evaluates the children, accumulating their content.
  The function's return value is ignored; content flows out as a side
  effect of body invocations.

**Done.**

- Dang: added `dang.ContextWithBlock(ctx, val)` (small upstream patch
  in `pkg/dang/eval.go`). Lets embedders set the block-arg context key
  programmatically — the surface syntax `Foo(...) { ... }` already
  sets the same key during parsing, but the key itself was unexported.

- `dangeval.LookupCallable(name)` returns a `dang.Callable` from the
  held eval env if one exists with that name.

- `dangeval.CallComponent(callable, props, body)` invokes a Dang
  function as a JSX component. The body closure is wrapped in a
  `componentBlock` (a custom `dang.Value` / `dang.Callable`
  implementation) and attached to the context via `ContextWithBlock`
  before calling. Returns the function's value (used for body-less
  components; ignored for body-ful).

- `dangeval.WithBindings(args, fn)` pushes the named bindings into
  both the eval env (via `Derive(true) + Bind`) and the type env
  (via `Clone + Add` with each value's `Type()`). The type-env
  extension is non-obvious but essential — without it,
  `InferFormsWithPhases` on a `{name}` snippet fails with "name not
  found" because Dang inference runs against the type env.

- `dangeval.New` now also auto-discovers `*.dang` files in the
  project directory (non-recursive, alphabetical), parses them, and
  merges all forms into the shared type and value envs via
  `InferFormsWithPhases` + `EvaluateFormsWithPhases`. Treats the
  directory as one module — no per-file import scoping for v1.

- `stages.Evaluate.VisitJSXElement` consults Dang scope between the
  built-in tier and the template fallback. A new `dispatchDang` helper
  bridges props (literal strings → `StringValue`; `{expr}` props →
  raw Dang values via `Eval`; anything else → stringified content) and
  drives the call. Accumulator pattern: body-ful components populate
  a per-element accumulator; body-less components contribute their
  return value bridged to content.

- `tests/dang_dispatch_test.go` covers five cases: body-less, body-ful
  with named bindings, iteration via `items.each { body(…) }`,
  missing function falls through to template fallback, and
  non-callable Dang binding (e.g. `pub Pal = "world"`) also falls
  through. `tests.Example.Run` was reordered to construct the
  evaluator *after* writing test inputs so any test-provided `.dang`
  files are picked up.

**Decisions worth knowing.**

- *Lookup rule: any callable Dang binding matching the JSX name is
  dispatched as a component.* No PascalCase check; lowercase Dang
  names just can't be triggered by JSX (which only matches uppercase
  tags). Non-callable bindings (lists, records, strings) fall through
  to the template fallback — they aren't an error.

- *Dispatch order: built-in → Dang → template.* Built-ins take
  precedence so language primitives (`<Section>`, `<Reference>`, …)
  can't be shadowed. Dang takes precedence over templates so a user
  who explicitly declares `pub Foo` always wins over an accidental
  `Foo.tmpl` file.

- *Block return type quirk.* The Go-backed block returns `NullValue`,
  but if the user declares it as `&body(…): T!` (non-null) and the
  function's body is just `body(…)`, Dang's runtime null-check fires
  ("null is not allowed for T!"). Workaround: the user function
  explicitly returns a truthy value at the end (`body(…); true`).
  Long-term we could make the block declare its return as nullable,
  or have `componentBlock.Call` return a value matching the declared
  type. For v1, the explicit-return idiom is fine.

- *Auto-discovery scope: project dir only, non-recursive.* Matches
  what `dang.toml` and `dagger.json` already do (lookup walks up from
  the input file's directory). No `dang/` subdirectory convention —
  the `.dang` files sit alongside the input docs.

- *Bridging props.* Literal `<Foo bar="x"/>` → `StringValue{"x"}`.
  Expression `<Foo bar={expr}/>` → raw Dang value via `Eval`. Other
  AST nodes → stringify (rare). This preserves type fidelity for the
  two common cases.

- *Single-threaded scope swap.* `WithBindings` mutates
  `Evaluator.evalEnv` / `typeEnv` for the duration of the body call
  and restores on return. Works because Booklit's evaluator is
  sequential; a future parallel evaluator would need per-goroutine
  scope frames.

**Known limitations / follow-ups.**

- *Helpers in `docs/booklitdoc/` not yet ported.* This phase only
  adds the capability; the actual cleanup (move OutputFrame /
  SyntaxHl / ColumnHeader / TemplateLink / Godoc to templates and/or
  Dang functions, then shrink `cmd/booklit-docs`) is its own slice.
  `<Define>` / `<LitSyntax>` / `<Columns>` and the chroma palette
  override stay as Go for now — they need either Dagger (per the
  agreed punts) or an API reshape.

- *Block return type quirk* (see above) — should revisit once we have
  a real use case that wants a non-trivial block return.

- *Templates are still Go html/template format.* Switching them to
  `.md` (MDX-style) is the next bigger piece (was Q7 → option (b):
  tier-3 first, mdx templates as a follow-up). Tier-3 dispatch is now
  in place; mdx-as-template is the next substantive phase.

**Reasonable next steps.** (a) Port the easy half of booklitdoc to
templates (the existing `.tmpl` files probably already cover most),
then shrink `cmd/booklit-docs` to just the helpers that still need Go.
(b) Phase 3b: mdx-as-template — uses tier-3 dispatch and the
`<Children/>` built-in to be added. Eliminates the Go html/template
dialect for user-facing templates. (c) Source-mapped errors (already
deferred). (d) JSX-inside-Dang (still deferred).

### 2026-05-31 — Port easy booklitdoc helpers to templates

`docs/booklitdoc/` shrinks from 376 → 307 lines as the helpers with
no Go logic move to pure templates that the tier-2 fallback handles.

**Done.**

- `<OutputFrame url="…"/>`: removed Go; renamed
  `output-frame.tmpl` → `OutputFrame.tmpl` and updated
  `{{.Partial "URL"}}` → `{{.Partial "url"}}` (camelCase as per the
  established convention). The original Go also constructed an
  unused `booklit.Link` for the Content; the template never read it,
  so it was dropped.

- `<TemplateLink tmpl="…"/>`: removed Go; new `TemplateLink.tmpl`
  interpolates the partial into the github.com URL and renders the
  filename as `<code>`. Uses `rawHTML` for the href interpolation
  (filenames are simple ASCII so no escaping wrinkle).

- `<SyntaxHl>…</SyntaxHl>`: removed Go AND template. Both were dead
  code — no caller in the docs or templates.

- `<ColumnHeader>` / `<Column>`: removed Go AND `column-header.tmpl`
  (also dead). `<Columns>` reaches into its raw `ast.JSXElement`
  children by name to extract title vs columns, so neither dispatcher
  is ever reached. The doc comment on `columnsFunc` now mentions
  this explicitly.

- `stages/evaluate.go` template fallback now defaults `body` to
  `booklit.Empty` when the JSX element has no children. Without
  this, a self-closing component (e.g. `<OutputFrame url="…"/>`)
  produced `Styled{Content: nil}` which panicked in
  `Styled.IsFlow` → `Content.IsFlow` during the collect stage.

The full docs site builds clean (`go run ./cmd/booklit-docs -i
docs/lit/index.md -o ... --html-templates docs/html`); spot-checked
rendered OutputFrame and TemplateLink output. `go test ./...` green.

**What's still in `cmd/booklit-docs` / `booklitdoc.go`.** Four
helpers and one init() side-effect, all blocked on either Dagger or
an API reshape:

- `<Define tag=… sig=…>desc</Define>` — uses chroma to syntax-
  highlight the signature and constructs a `<Target>` for the tag.
- `<LitSyntax>code</LitSyntax>` — chroma highlighting + regex sweep
  over the highlighted output to linkify `\fn` references.
- `<Godoc ref=pkg.Symbol/>` — `strings.SplitN` + Link construction.
  Could become a template if we add a string-split helper to
  `HTMLFuncs`, but punted for now.
- `<Columns>` — AST child introspection (matches `<ColumnHeader>` /
  `<Column>` by name). Stays a built-in until either the API
  reshapes or AST inspection becomes a first-class capability.
- chroma `styles.Fallback` override — init()-time global state for
  the booklitdoc palette.

**Reasonable next steps.** (a) Phase 3b: mdx-as-template — would
make the easy-helper templates themselves expressible in the same
syntax authors use for content. Larger, but unifies the dialects.
(b) Shrink the remaining Go further by porting Godoc to a template
once a string-split helper exists, or by making it a Dang function
that returns raw HTML. (c) Plan the Dagger story for LitSyntax /
Define so the binary can eventually disappear.
