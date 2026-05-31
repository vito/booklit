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

The bridge (`dangeval/bridge.go`) maps Dang values to `booklit.Content`:
`StringValue` → `String`; `IntValue`/`FloatValue`/`BoolValue` →
stringified `String`; `ListValue` → `Sequence`; `NullValue` → empty;
`ContentValue` → its carried `Content` verbatim; anything richer is an
error. The Evaluator is one per build session, single-threaded.

A Dagger session is implicit: any project with a `dagger.json` gets
the module's functions in scope. Booklit doesn't need its own
"Dagger dispatch tier" — Dang already handles that, and a tag-level
JSX `from="..."` syntax (in the original Phase 4 sketch) would have
just duplicated the existing import machinery.

## Doc helpers (`docs/booklitdoc/`) collapsed

What used to be `~376 lines` is now `~165`. `<Define>` and `<Godoc>`
moved to `docs/html/Define.md` and `docs/html/Godoc.md`; `<Columns>`,
`<Column>`, and `<ColumnHeader>` are now plain `<div>`-wrapper mdx
templates (`docs/html/Columns.md`, `Column.md`, `ColumnHeader.md`) with
the layout driven by CSS (`.columns > .column:first-child` is the narrow
description column); `<OutputFrame>`, `<TemplateLink>`, `<SyntaxHl>`
either moved to templates or were dead code. What's still Go:
`<LitSyntax>` (chroma + regex) and the chroma `styles.Fallback` palette
override.

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
- `cmd/booklit-docs/` — separate binary that imports the docs-site
  built-ins on top of `cmd/booklit`

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

- Per-document remote-Dagger-module imports. Locally-bound Dagger
  functions work via `{expr}` from any project with a `dagger.json`;
  pulling in an out-of-tree module for one document would mean
  reaching for `dang.toml`'s import config, not a Booklit-side syntax.
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
  pass).
