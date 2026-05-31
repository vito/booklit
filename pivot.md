# Pivot: what changed since origin/main

A short tour for someone who only knows pre-pivot Booklit. The branch
is one big rewrite of the extension model. The section tree, renderer,
and Markdown prose support are unchanged; everything else around them
moved.

## Invocation syntax: `\foo{}` → JSX

`\title{Hi}` became `<Title>Hi</Title>`. Multi-arg invokes mapped
positional args to named props, with one "content" arg (typically the
display text or body) becoming JSX children:

- `\link{display text}{url}` → `<Link target="url">display text</Link>`
- `\image{path}{description}` → `<Image path="path" description="..."/>`
- `\target{tag}{title}` → `<Target tag="tag">rich title</Target>` (children
  become Title; the third Content slot from baselit is dropped, and a
  `title="..."` prop is available as a plain-string shorthand)
- `\reference{tag}{display}` → `<Reference tag="tag">display</Reference>`

Prop names are camelCase. Lowercase tags pass through as raw HTML
(React rule), so a template like `<div class="card">{title}</div>`
mixes literal HTML with `{expr}` prop interpolation. Markdown features
(`#` headings, `- lists`, `| tables |`, fenced code) still work;
headings still auto-create sections.

`marklit/jsx_parser.go` + `jsx_block_parser.go` are the goldmark
extensions; `ast.JSXElement` and `ast.JSXExpression` are the new AST
nodes, layered *alongside* the existing pipeline rather than replacing
it. `ast.Invoke` and the `\foo{}` parser are NOT gone: Markdown prose
still lowers to `ast.Invoke` for every built-in feature (`#` headings →
`\section` / `\title`, links, images, code, code blocks, tables, lists,
insets, references, and raw HTML → `raw-html` / `raw-html-block`), and
`marklit` still registers the `\foo{}` inline parser, so backslash
invokes keep working in `.md`. `.lit` files still go through the
original PEG parser (`ast.ParseReader`). These `ast.Invoke` nodes are
evaluated by `VisitInvoke`, which dispatches by reflection against the
section's plugins — i.e. `baselit` — exactly as before the pivot. The
JSX layer (`VisitJSXElement`) is the new, parallel path.

## Plugin system: gone

Deleted: `--plugin` CLI flag, the `BOOKLIT_REEXEC` re-exec dance,
`booklit.RegisterPlugin` / `LookupPlugin`, `\use-plugin`, every example
Go plugin (`chroma/`, `docs/hello/`, `docs/go/`,
`tests/fixtures/*-plugin/`). New components are added by dropping a
file, not by recompiling Booklit.

What's gone is the *user-facing, load-a-compiled-plugin* mechanism. The
internal `booklit.Plugin` / `PluginFactory` types and the reflection
dispatch in `VisitInvoke` remain — that's how `baselit` is still wired
in to evaluate Markdown-emitted `ast.Invoke` nodes (see below).

## Three-tier JSX dispatch

`stages/evaluate.go::VisitJSXElement` resolves a tag through three
tiers; an unknown name is an error (hard cutover, no Styled fallback):

1. **Built-in** — Go function registered in the `builtins/` package.
   Handles the language primitives (`<Title>`, `<Section>`,
   `<Reference>`, `<Target>`, `<TableOfContents>`, `<Code>`,
   `<CodeBlock>`, `<Syntax>`, …) and a few quality-of-life additions
   (`<For>`, `<If>`, `<Unless>`, `<Children/>`, `<RawHTML>`). ~34 live
   in `builtins/`, most reimplementing the behaviour natively (a few,
   e.g. `<Reference>` / `<Code>`, still call into `baselit`). `baselit`
   itself has not been retired: it is still registered as the base
   `PluginFactory` and remains the reflection target for the
   Markdown-emitted `ast.Invoke` nodes described above (some of which,
   like tables and lists, have no JSX counterpart at all).
2. **Dang function** — a `pub PascalCase(...)` callable in scope.
   `dangeval.LookupCallable` + `CallComponent` bridge props as named
   args and wrap the JSX children as a `&body` block; each `body(...)`
   call from inside the Dang function pushes its named args into scope
   and re-evaluates the children. Used for closures over per-project
   data and for parametric control flow that needs `&body`.
3. **mdx template** — `<dir>/<Name>.md` parsed by a small custom
   tokenizer (`templates/parse.go`) that recognizes raw HTML, JSX
   elements, and `{expr}` and nothing else (no Markdown). The
   evaluator binds props by name in Dang scope and binds the JSX
   children's rendered content as `children` (a `dangeval.ContentValue`
   so nested styling survives). `{children}` and `<Children/>` both
   emit it.

Renderer-internal templates (page, section, sidebar, etc.) stay as Go
`html/template` — they are framework infrastructure, not user
extension. Built-ins that emit `booklit.Styled{Style: …}` still
resolve through the renderer's `.tmpl` lookup as before; only the JSX
auto-wrap is gone.

## Embedded Dang interpreter

`dangeval/` wraps `github.com/vito/dang` so `{expr}` interpolations in
JSX evaluate as real Dang code. `dangeval.New` walks up from the input
file looking for `dang.toml` (Dang's GraphQL imports + Dagger session
config) and `dagger.json` (a local Dagger module). It also scans the
project directory for `*.dang` files, treats them as one module, and
merges the forms into the held type + value envs. Booklit's docs use
this for `docs/lit/helpers.dang` (godoc URL composition, the
componentName kebab→Pascal converter).

The bridge (`dangeval/bridge.go`) maps Dang values to `booklit.Content`.
`ToContent` handles the primitives: `StringValue` → `String`;
`IntValue`/`FloatValue`/`BoolValue` → stringified `String`; `ListValue`
→ `Sequence`; `NullValue` → empty; `ContentValue` → its carried
`Content` verbatim. `Evaluator.ContentFromValue` is the richer,
section-aware path used for JSX/`{expr}` results: on top of the above it
decodes content returned by a Dagger module (see "Content from a Dagger
module") and rehydrates `Reference`/`Target` nodes against the current
section. The Evaluator is one per build session, single-threaded.

A Dagger session is implicit: `dangeval.New` finds the `dagger.json`
walking up from the input file, introspects the module's schema for type
checking, and **serves** it into the session (with its dependencies) so
its functions are callable at runtime — the serve is the piece that
turns a type-checking-only import into a runnable one. There is no
separate Booklit "Dagger dispatch tier"; calls go through `{expr}` like
any other Dang code. Note that what an introspected *module* exposes on
the session `Query` is its **dependencies + core API**, not the
module's own functions — so docs reach the highlighter as
`booklitdoc.litSyntax(...)` (the dependency), not as a root-module
function.

## Content from a Dagger module

A Dagger module can return Booklit content. The wire format is
`contentjson` — a tagged-union JSON encoding of the serializable subset
of `booklit.Content`, with the dependency-free node schema + builder
constructors split into `contentjson/wire` so a producer needn't import
all of `booklit`. A module builds the tree with `wire`'s constructors
and returns it as Dagger's `JSON` scalar; Booklit recognizes the
`JSON`-typed return in `ContentFromValue` and decodes it back into
native content. `JSONValue!` returns work too (the bridge forces
`.contents`), letting a module compose results lazily before Booklit
materializes them.

In-process-only content — `Section`, `TableOfContents`, `Lazy` — can't
serialize and errors from `Marshal`. Stateful-but-nameable content
(`Reference`, `Target`) crosses carrying only a tag name and is re-bound
to the live section on decode, so cross-references survive the round
trip.

The first real use is `<LitSyntax>`: the `docs/html/LitSyntax.md`
template calls `{booklitdoc.litSyntax(code: code)}`, which runs the
`dagger/booklitdoc` Go-SDK module (chroma highlighting + `\function`
linkification) and returns `contentjson`. The module depends on the
local `booklit` module via a `go.mod` replace + `dagger.json` includes
(the standard Dagger monorepo pattern) and is installed as a dependency
of the docs module so it lands on the introspected `Query`.

## Doc helpers (`docs/booklitdoc/`) collapsed

What used to be `~376 lines` is now `~44` — just the chroma
`styles.Fallback` palette override, still applied in-process because
`baselit` highlights fenced code blocks site-wide. Everything else
became a template or moved out of process: `<Define>` / `<Godoc>` →
`docs/html/Define.md` / `Godoc.md`; `<Columns>` / `<Column>` /
`<ColumnHeader>` → plain `<div>`-wrapper mdx templates with the layout
driven by CSS (`.columns > .column:first-child` is the narrow
description column); `<OutputFrame>` / `<TemplateLink>` / `<SyntaxHl>` →
templates or dead code; and `<LitSyntax>` → the Dagger module above. The
binary `cmd/booklit-docs` now exists only to install that palette and
bundle the docs' Dagger dependency.

## File map (new code)

- `marklit/jsx_*.go` — the JSX parser; produces `ast.JSXElement` /
  `ast.JSXExpression`
- `builtins/` — the first JSX dispatch tier (many former baselit
  methods now have a JSX-shaped Func counterpart taking a Context arg;
  `baselit` still exists and still handles Markdown-emitted invokes)
- `dangeval/` — Dang interpreter wrapper + the Content↔Value bridge +
  the tier-3 component-call plumbing
- `templates/` — tier-3-and-a-half template registry + the custom
  template parser
- `contentjson/` — JSON wire format for `booklit.Content`, plus the
  dependency-free `wire` subpackage producers build with
- `dagger/booklitdoc/` — Go-SDK Dagger module: highlights Booklit
  source and returns it as `contentjson`
- `cmd/booklit-docs/` — separate binary that installs the docs palette
  (`docs/booklitdoc`) on top of `cmd/booklit`; the docs' highlighter is
  a Dagger dependency, not compiled in

## Design notes

- **camelCase props end to end.** Templates reach for a prop by name
  (`{title}`), not via PascalCase Partials lookup.
- **Single-line vs multi-line children.** `<Title>x</Title>` on one
  line parses children as inline (matching the old `\title{x}`); a
  multi-line `<Section>` parses children as block. Nested elements use
  their own line span, so an inline `<Title>` inside a multi-line
  `<Section>` keeps inline semantics.
- **`<Target>` semantics.** Children become Title content (the
  display text references fall back to). Matches baselit's legacy
  variadic-position shape and lets Define use
  `<Target tag={tag}><Syntax language="html"><{componentName(tag: tag)}>
  </Syntax></Target>` to keep the highlighted-tag title.
- **Hard-cutover posture.** No backwards compatibility shims; unknown
  JSX components error rather than wrap in `Styled`.

## What's not implemented

- Per-document remote-Dagger-module imports. Local Dagger module
  functions (including content-returning ones — see "Content from a
  Dagger module") work via `{expr}`; pulling in an out-of-tree module
  for one document would mean reaching for `dang.toml`'s import config,
  not a Booklit-side syntax.
- `<LitSyntax>`'s `JSONValue!` (lazy) variant. The bridge supports it,
  but the module returns the `JSON` scalar because its pinned Go SDK
  predates the engine's `JSONValue` constructor; needs an SDK bump.
- Passing raw source *text* (not rendered content) as template/JSX
  children, so `<LitSyntax>` takes `code` as an attribute rather than a
  fenced body. Multi-line snippets are awkward until then.
- JSX literals inside Dang expressions (`{items.map(t => <Foo>…
  </Foo>)}`). Iteration and conditionals are covered by `<For>` and
  `<If>` / `<Unless>` built-ins instead.
- Source-mapped errors for template/Dang failures.
- Per-section `{expr}` scope (every snippet evaluates against the same
  global env).
- `booklit-init` scaffolding for new projects.

## Where to read more

- `jsx-dang.md` — original pivot plan and append-only progress log.
- `phase-3b.md` — mdx-as-template phase plan and its progress log.
- `decisions.md` — fork-in-the-road notes from the late autonomous
  session (e.g. why JSX-in-Dang and Phase 4 were left for a co-design
  pass). Note: its assumption that a root-module function like
  `{build(...)}` is callable from docs is wrong — introspection exposes
  dependencies + core, not the module's own functions.
- `dagger-content.md` — current state of the Dagger-content work
  (`contentjson`, the bridge, the served module, `<LitSyntax>`) and the
  open follow-ups; the place to pick up next.
