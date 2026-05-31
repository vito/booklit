# Pivot: what changed since origin/main

A short tour for someone who only knows pre-pivot Booklit. The branch
is one big rewrite of the extension model. The section tree, renderer,
and Markdown prose support are unchanged; everything else around them
moved.

## Invocation syntax: `\foo{}` → JSX

`\title{Hi}` became `<Title>Hi</Title>`. `\foo{a}{b}{c}` with positional
args became `<Foo prop="a">{b}{c}</Foo>` with camelCase named props.
Lowercase tags pass through as raw HTML (React rule), so a template
like `<div class="card">{title}</div>` mixes literal HTML with prop
interpolation. Markdown features (`#` headings, `- lists`, `| tables |`,
fenced code) still work; headings still auto-create sections.

`marklit/jsx_parser.go` + `jsx_block_parser.go` are the goldmark
extensions; `ast.JSXElement` and `ast.JSXExpression` are the new AST
nodes. The legacy `ast.Invoke` and its `\foo{}` parser are gone for
user content but still exist to support the auto-emitted
`<RawHTML>` / `<RawHTMLBlock>` invocations.

## Plugin system: gone

Deleted: `--plugin` CLI flag, the `BOOKLIT_REEXEC` re-exec dance,
`booklit.RegisterPlugin` / `LookupPlugin`, `\use-plugin`, every example
Go plugin (`chroma/`, `docs/hello/`, `docs/go/`,
`tests/fixtures/*-plugin/`). New components are added by dropping a
file, not by recompiling Booklit.

## Three-tier JSX dispatch

`stages/evaluate.go::VisitJSXElement` resolves a tag through three
tiers; an unknown name is an error (hard cutover, no Styled fallback):

1. **Built-in** — Go function registered in the `builtins/` package.
   Handles the language primitives (`<Title>`, `<Section>`,
   `<Reference>`, `<Target>`, `<TableOfContents>`, `<Code>`,
   `<CodeBlock>`, `<Syntax>`, …) and a few quality-of-life additions
   (`<For>`, `<If>`, `<Unless>`, `<Children/>`, `<RawHTML>`). All ~30
   live in `builtins/`. The old `baselit` package stays as a Go API
   the built-ins delegate into; it is no longer a "plugin".
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
the module's functions in scope. The `<Foo from="..."/>` anonymous-import
syntax from the original Phase 4 plan is **not** implemented — only the
local-module path works today.

## Doc helpers (`docs/booklitdoc/`) collapsed

What used to be `~376 lines` is now `~150`. `<Define>` and `<Godoc>`
moved to `docs/html/Define.md` and `docs/html/Godoc.md` (mdx
templates); `<OutputFrame>`, `<TemplateLink>`, `<SyntaxHl>`,
`<ColumnHeader>`, `<Column>` either moved to templates or were dead
code. What's still Go: `<Columns>` (AST child introspection),
`<LitSyntax>` (chroma + regex), and the chroma `styles.Fallback`
palette override.

## File map (new code)

- `marklit/jsx_*.go` — the JSX parser; produces `ast.JSXElement` /
  `ast.JSXExpression`
- `builtins/` — the first dispatch tier (everything that was a baselit
  method is now a JSX-shaped Func with a Context arg)
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

- Dagger anonymous-import syntax (`<Foo from="github.com/.../mod"/>`).
  Locally-bound Dagger functions work via `{expr}`; out-of-tree module
  references are not yet a thing.
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
