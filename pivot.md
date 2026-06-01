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
it. User-facing `\foo{}` in Markdown has been removed (decision 2),
but `ast.Invoke` is NOT gone: Markdown prose still lowers to
`ast.Invoke` for every built-in feature (`#` headings → `\section` /
`\title`, links, images, code, code blocks, tables, lists, insets,
references, and raw HTML → `raw-html` / `raw-html-block`). The
`[#tag]` reference shorthand also produces an `ast.Invoke`. The PEG
`.lit` parser is gone (decision 1); `marklit.Parse` is the only
source-text parser. These `ast.Invoke` nodes are evaluated by
`VisitInvoke`, which dispatches by reflection against the section's
plugins — i.e. `baselit` — exactly as before the pivot. The JSX
layer (`VisitJSXElement`) is the new, parallel path.

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

## Syntax highlighting: `treehighlight`

Fenced code blocks (and inline code) are highlighted by `treehighlight/`,
a thin wrapper over tree-sitter that ships its own Booklit grammar
(`treehighlight/internal/tree_sitter_booklit/`). `baselit.Syntax` /
`baselit.CodeBlock` call it, group the resulting chunks into spans, and
emit real `booklit.Reference` nodes for captures whose kebab-case form
matches an existing tag — so fenced examples of `<IncludeSection/>` or
`\section{...}` are automatically linkified to their definitions.

cgo is the production path. If Booklit is built with `CGO_ENABLED=0`,
`treehighlight` compiles to an escaped plain-code fallback (no spans,
no links).

This replaced the earlier chroma-based highlighter — and replaced the
short-lived `<LitSyntax>` Dagger experiment (the in-tree
`dagger/booklitdoc` Go-SDK module + the `docs/booklitdoc` palette + the
`cmd/booklit-docs` wrapper binary) which all got deleted once
treehighlight handled the linkification natively. The docs site now
runs straight from `cmd/booklit` with no docs-specific Go.

## Content from a Dagger module (infrastructure-only, no in-tree user)

A Dagger module *can* return Booklit content. The wire format is
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

There is currently **no in-tree user**: `<LitSyntax>` was the first and
only consumer, deleted when treehighlight subsumed it. `contentjson/`
+ `contentjson/wire/` (~440 lines) sit as speculative infrastructure
waiting for the next out-of-process content producer. Whether to keep
that infrastructure live or remove it until a real use case appears is
one of the open questions below.

## File map (new code)

- `marklit/jsx_*.go` — the JSX parser; produces `ast.JSXElement` /
  `ast.JSXExpression`
- `builtins/` — the first JSX dispatch tier (many former baselit
  methods now have a JSX-shaped Func counterpart taking a Context arg;
  `baselit` still exists and still handles Markdown-emitted invokes)
- `dangeval/` — Dang interpreter wrapper + the Content↔Value bridge +
  the tier-3 component-call plumbing
- `templates/` — tier-3 template registry + the custom template parser
- `treehighlight/` — tree-sitter–based syntax highlighter; emits spans
  and Booklit references for matched captures
- `contentjson/` — JSON wire format for `booklit.Content`, plus the
  dependency-free `wire` subpackage producers build with (no in-tree
  consumer at the moment)

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

- A Dagger-as-JSX-tier story. `<Foo from="github.com/.../mod"/>` syntax
  doesn't exist; calling a Dagger function happens only through
  `{moduleFn(...)}` Dang expressions. See "Open questions" below for the
  design gap.
- Passing raw source *text* (not rendered content) as template/JSX
  children. Multi-line snippets that want their literal source can only
  use a fenced code block today.
- JSX literals inside Dang expressions (`{items.map(t => <Foo>…
  </Foo>)}`). Iteration and conditionals are covered by `<For>` and
  `<If>` / `<Unless>` built-ins instead.
- Source-mapped errors for template/Dang failures.
- Per-section `{expr}` scope (every snippet evaluates against the same
  global env).
- `booklit-init` scaffolding for new projects.

## Standard project layout (target, after the cleanup checklist)

```
project/
├── lit/                      # Content: *.md (and helpers.dang next to them)
│   ├── index.md              # `--in` target; entry point
│   ├── other-page.md
│   └── helpers.dang          # Auto-loaded by dangeval (non-recursive scan)
├── components/               # User JSX components: <Name/> → components/Name.md
│   ├── Card.md               # Same MarkDangJSX format as content files
│   └── Define.md
├── html/                     # Renderer override templates (Go html/template);
│   ├── page.tmpl             # only used to override framework infrastructure
│   ├── section.tmpl          # like page/section/sidebar. Most projects do
│   └── sidebar.tmpl          # not need this directory at all.
├── css/                      # Static assets (CSS, fonts, images); copied to
├── favicon.ico               # the output dir by the renderer.
├── dagger.json               # Dagger module metadata (sdk = "dang")
├── .dagger/                  # The project's own Dagger module (CI, build,
│   └── main.dang             # test). Its functions are NOT in {expr} scope —
│                             # only its *dependencies* + core API are.
└── dist/                     # `--out` target; generated HTML
```

The format used by both content files and component files is
**MarkDangJSX**: Markdown prose + JSX elements + `{Dang expressions}`.
There is one parser (`marklit`), one set of semantics. The distinction
between "content" and "component" is operational, not syntactic:

- A **content file** lives in `lit/` and is loaded as a top-level
  section (its headings build the section tree).
- A **component file** lives in `components/` and is invoked from
  *other* files as `<Name/>`. Props bind in Dang scope; the caller's
  JSX children bind as `children`. The same `lit/` content file could
  be used as a component, and vice versa — the difference is "did
  someone call it as `<Name/>`?".

Inside a component file, Markdown applies the same way it does in
content (the React/MDX convention): text inside a *lowercase* element
(`<div>`, `<pre>`) is raw HTML — no Markdown — and text inside a
PascalCase JSX element is parsed as Markdown as usual. So
`<Card>**emph**</Card>` produces emphasis inside the rendered
component; `<div>**not-emph**</div>` does not.

A few load-bearing details:

- **The `--in` flag's directory anchors everything Dang-side.**
  `dangeval.New` walks up from `filepath.Dir(--in)` looking for
  `dang.toml` and `dagger.json`, and scans that same directory
  (non-recursive) for `*.dang` files. That's why `helpers.dang` lives
  inside `lit/` next to the content, not at the project root.
- **Dagger comes for free.** Any project with a `dagger.json` gets
  the served module's dependencies + core API available as `{expr}`
  calls. There is no `<Foo from="..."/>` syntax and there will not be
  one: that role is filled by Dagger-the-language (Dang), not a
  separate Booklit-level dispatch tier.
- **Renderer-internal `.tmpl` files stay Go templates.** The 28 files
  in `render/html/` (page, section, styled, list, table, code-block,
  …) are embedded at compile time and form the rendering substrate
  built-ins emit `Styled{Style: …}` against. Projects that need to
  override one drop a same-named `.tmpl` into `html/`.

## Decisions

Recap of the questions raised in the previous revision and the
decisions taken. Each one becomes a checklist item below.

1. **PEG `.lit` parser — remove.** Standardize on MarkDangJSX. The
   ~2700 lines of generated PEG + the pigeon Makefile rule + the `.lit`
   branch in `load/processor.go` all go. Tests using `Ext: ".lit"` get
   rewritten as `.md`; `.lit` fixture strings (used by
   `<IncludeSection>` tests) get renamed.
2. **User-facing `\foo{}` in Markdown — remove.** The
   `NewInvokeInlineParser` registration goes away so authors can't
   write `\italic{x}` inline. Internal lowering keeps emitting Invokes
   for now; that goes in step 3.
3. **`ast.Invoke` reflection dispatch — replace with direct JSX
   lowering.** `marklit/convert.go` lowers Markdown straight to
   `ast.JSXElement` nodes. `VisitInvoke`, `Section.PluginFactories`,
   `Section.Plugins`, `UsePlugin`, `booklit.Plugin`,
   `booklit.PluginFactory`, and `baselit/`'s Plugin-method shape all
   collapse: baselit's primitives become `builtins/`-registered
   functions alongside `<Title>`/`<Section>`/etc. The "plugin" word
   stops being used internally.
4. **Stale planning docs at the project root — remove.**
   `jsx-dang.md`, `phase-3b.md`, `decisions.md`, `dagger-content.md`
   all deleted; this file becomes the only top-level pivot doc.
5. **`html/` directory — split.** Component definitions move to
   `components/` (PascalCase `.md` files). The `html/` directory
   remains for Go html/template overrides of renderer-internal
   templates (page, section, sidebar), but most projects do not need
   it. The `--html-templates` flag goes away; both locations become
   conventions Booklit looks up automatically next to `--in`.
6. **`<Foo from="..."/>` — slop, removed from the roadmap.** Dagger
   modules are reachable through `{expr}` because dangeval already
   serves the local `dagger.json` module into the Dang env. A
   dedicated JSX-tag tier for remote modules is unnecessary; one-off
   imports remain a `dang.toml` concern, not a Booklit-side syntax.
7. **Body/children passing to Dagger functions — punted.** When
   needed, the simplest path is the eager-pre-render approach: render
   children once and pass them as `contentjson` to the Dagger
   function. The callback-channel and source-text variants are
   parked; revisit only when a concrete use case demands one.
8. **`contentjson` — kept.** It's the wire format Dagger-provided
   plugins return content through. No in-tree consumer right now, but
   the infrastructure stays in place rather than being torn down and
   rebuilt later.
9. **Template format unification — same MarkDangJSX as content.**
   The custom template tokenizer in `templates/parse.go` goes away;
   components are parsed by `marklit` like everything else. This gives
   Markdown-inside-`<Component>` for free (per the React/MDX
   lowercase-vs-PascalCase convention) and removes the second parser
   that the codebase has to maintain.
10. **Partials — remove.** `\set-partial{name}{...}` /
    `<SetPartial name="...">...</SetPartial>` plus the
    `{{.Partial "Foo" | render}}` template hook were a workaround for
    the missing component system: a way to thread a named blob of
    content into a renderer template. Components (templates + Dang
    functions, both with full prop binding and `<Children/>`) cover the
    same use case more cleanly. Remove `Section.SetPartial`,
    `Section.Partial`, `Section.Partials`, the `<SetPartial>` builtin
    and `baselit.SetPartial` method, the partial-template fixtures
    (`partial-template.tmpl` + the `set-partial-read-template` routing
    in `page.tmpl`/`full-styled-page.tmpl`), and the orphaned
    `custom-style.tmpl`/`inline-custom-style.tmpl` fixtures.
    `Styled.Partials` stays — it carries renderer metadata (e.g.
    `Language` on highlighted code blocks), not user-set content.
11. **Retire `Styled` in favor of structured content + stdlib
    components.** `Styled{Style, Block, Content, Partials}` is the
    last partial-abstraction in the content model: a string-keyed
    template lookup, a flag-based `IsFlow` override, and a
    single-use `Partials` metadata bag (`Language` on highlighted
    code blocks — set in one place, read by zero templates or
    renderer code, fully dead weight). Replace with two structured
    types and a `go:embed`-ed stdlib `components/` dir:

    - **`RawElement{Tag, Attrs, Content}`** for lowercase JSX,
      `<RawHTML>`, and the code builtins. `IsFlow()` returns
      `!htmltags.Block[Tag]` — the outermost element's tag decides
      block/flow. No flag, no string-parsing of
      `Styled{raw-html, "<div>"}` bytes at any later stage.
    - **`RawFragment{HTML}`** for pre-rendered HTML chunks (the
      syntax highlighter's `<span>` wrappers around tokens). Always
      flow.

    Code builtins collapse onto these. `<Code>x</Code>` emits
    `RawElement{Tag: "code", Content: x}` when its content is
    flow, or `RawElement{Tag: "pre", Content: RawElement{Tag:
    "code", Content: x}}` when block — same IsFlow-on-children
    introspection as today, no new type. `<CodeBlock language="X">`
    / `<Syntax language="X">` emit the same pre/code shape with a
    `class="language-X"` attribute on the `<code>` (the
    Prism/highlight.js convention, and the natural place for the
    metadata that `Partials{Language}` was carrying for no
    reader). `verbatim.tmpl` / `code-block.tmpl` / `code-flow.tmpl`
    are all redundant once `<Code>` returns nested `RawElement`s
    directly — they delete with no replacement Go visitor needed.

    Block/flow detection moves from a Block flag to an
    introspection rule: the outermost HTML element of the
    rendered content decides. Lowercase tags consult an embedded
    CommonMark/HTML5 block-tag set (~60 entries — goldmark
    doesn't export theirs, hand-maintained). PascalCase
    components introspect their evaluated output — `<Aside>` is
    block because its body's outermost element is `<blockquote>`,
    `<Larger>` is flow because its body is a `<span>`. The 5
    styling shims (`<Larger>`, `<Smaller>`, `<Strike>`, `<Inset>`,
    `<Aside>`) become MarkDangJSX components shipped in an
    embedded `components/` FS, with project `components/` taking
    precedence by name (consistent with the `--in`-anchored
    override pattern used elsewhere).

    Two new layout behaviors fall out:

    - **Standalone-flow JSX gets paragraph-wrapped.** When the
      block JSX parser claims an element and the evaluated
      content is flow, wrap it in `Paragraph` so
      `<Larger>x</Larger>` standalone still renders as
      `<p><span>x</span></p>`. Mirrors the existing unwrap rule
      in `VisitParagraph` in the inverse direction.
    - **Mid-paragraph block content breaks the paragraph.** When
      a `Sequence` line in a Markdown paragraph contains a block
      element (`!IsFlow()`), split the paragraph at that point —
      lines before are wrapped in a `<p>`, the block content
      renders unwrapped, lines after start a new paragraph.
      Matches CommonMark's raw-block-HTML-in-paragraph behavior;
      replaces today's invalid `<p>…<blockquote>…</p>`
      passthrough (which was previously masked by the
      `Block: node.MultiLine` flag on lowercase JSX).

    `Styled`, the `Style` type, all `Style*` constants, the
    `booklit.Partials` type, the matching `.tmpl` files
    (`larger`/`smaller`/`strike`/`inset`/`aside`/`verbatim`/
    `code-block`/`code-flow`/`raw-html`), and contentjson's
    `Style`/`Block`/`Partials` wire fields all delete.

## Cleanup checklist

Concrete tasks, in dependency order. Each line links back to a
"Decisions" item above.

- [x] **Sweep stale references** (decision 4 prep + decision 5
      prep). All remaining references to `cmd/booklit-docs`,
      `docs/booklitdoc/`, `dagger/booklitdoc/`, and `<LitSyntax>`
      live in the planning-doc files (next item) or in pivot.md
      itself. No code-side cleanup needed.
- [x] **Delete planning docs** (decision 4): `jsx-dang.md`,
      `phase-3b.md`, `decisions.md`, `dagger-content.md`. `pivot.md`
      is now the only top-level pivot doc.
- [x] **Migrate test fixtures off `\foo{}`** (decision 2 prep +
      decision 1 prep). Done in commit `da05f5f` (2026-06-01).
      Approach: prefer Markdown wherever it can express the same
      content (`# Title`, `## Sub`, `[text](url)`, `![alt](src)`,
      `*italic*`, `**bold**`, fenced code blocks, `> inset`, GFM
      tables and lists, `[#tag]` for tag-only references); fall back
      to JSX only where Markdown has no equivalent (`<Larger>`,
      `<Smaller>`, `<Strike>`, `<Superscript>`, `<Subscript>`,
      `<Aside>`, `<Aux>`, `<Definitions>`/`<Definition>`,
      `<Target tag="..."/>`, and the section ops
      `<Styled>`/`<SplitSections/>`/`<SinglePage/>`/`<IncludeSection/>`/
      `<TableOfContents/>`/`<OmitChildrenFromTableOfContents/>`).
      Three `Ext: ".lit"` prose cases stay (mid-word invokes, the
      three `\code{...}` flavors, indent tracking) — they exercise
      `.lit`-specific parser quirks and go away in step 5.
      `tests/partials_test.go` is deleted entirely; see decision 10.
- [x] **Remove user-facing `\foo{}` parsing in Markdown**
      (decision 2). `NewInvokeInlineParser` registration dropped
      from `marklit.go`'s `newParser` and `Extension.Extend`.
      `preprocess` reduced to CRLF normalization + `{- comment -}`
      stripping. Dead code removed in the same pass:
      `invoke_parser.go`, `verbatim.go`, `InvokeBlockNode` +
      `convertInvokeBlock`, `parseAllBracedArgs` /
      `parseBracedContent`, the placeholder machinery
      (`tryResolvePlaceholder` / `resolveEmbeddedPlaceholders` /
      `stripPlaceholders`), and the `ArgType` field on
      `InvokeNode` (only `[#tag]` produces `InvokeNode` now, and
      it uses a single normal arg). `splitInvokeOnlyParagraph`
      narrowed to JSX-only and renamed `splitElementOnlyParagraph`.
      Obsolete `\foo{}` / `{{…}}` / `{{{…}}}` tests in
      `marklit/marklit_test.go` deleted (~25 cases).
- [x] **Rewrite remaining `.lit` test fixtures as `.md`**
      (decision 1 prep). "invokes interspersed in words" ported to
      `.md` as "inline JSX interspersed in words" (`This<Italic>is
      </Italic>a test.`). The other two `Ext: ".lit"` cases —
      "inline code and code blocks" and "code block indent
      tracking" — exercised the three `\code{}`/`\code{{}}`/
      `\code{{{}}}` variants and `\code{{` indent-tracking; both
      were pure `.lit` parser quirks with no MarkDangJSX surface,
      and fenced code blocks are already covered in the same
      file. Deleted rather than ported. `tests/sections_test.go`
      has no `.lit` fixture strings.
- [x] **Delete the PEG parser** (decision 1). `ast/booklit.peg`,
      `ast/booklit.peg.go`, and `ast/errors.go` (PEG-only
      `UnpackError` / `parserError` plumbing) all gone. The
      Makefile `ast/booklit.peg.go: ast/booklit.peg` rule and the
      `pigeon` tool entry + indirect require dropped from
      `go.mod`; `go mod tidy` also removed the transitive
      `golang.org/x/mod` and `golang.org/x/tools` indirects.
      `load/processor.go` no longer branches on extension —
      `marklit.Parse` handles everything via `io.ReadAll`.
      `booklit.ParseError` is now unused but kept; no harm and
      it's part of the public API.
- [x] **Split templates → components** (decisions 5 + 9).
      `docs/html/*.md` moved to `docs/components/` (git-mv'd, history
      preserved). `tests/fixtures/Card.md` → `tests/components/Card.md`;
      the harness reads from `<tmpdir>/components/` + a shared
      `components/` dir. `booklitcmd` walks up from `filepath.Dir(--in)`
      looking for a sibling `components/` dir (`findComponentsDir`),
      so the `--html-templates` flag is no longer required for the
      standard layout. `templates/parse.go` + `templates/raw.go`
      deleted; the Registry now calls `marklit.ParseInlineArg` so a
      single-paragraph template doesn't pick up a stray `<p>` wrap.
      The bigger change is in marklit itself (see findings below):
      one unified JSX parser handles both lowercase HTML and
      PascalCase JSX, with case-based dispatch in
      `VisitJSXElement` — lowercase tags wrap their evaluated
      children in literal `<name attrs>...</name>` raw HTML;
      PascalCase routes through the three existing tiers.
      Markdown-inside-Component verified with a new test (`<Card>
      **bold**</Card>`).
- [x] **Rename / retire `--html-templates`** (decision 5). Flag and
      its struct field (`HTMLEngine.Templates`) gone from
      `booklitcmd/command.go`. The renderer now auto-discovers an
      `html/` directory the same way `components/` is found —
      `findProjectSubdir(--in, "html")` walks up from
      `filepath.Dir(--in)` looking for a sibling. Both auto-discoveries
      share one helper now. Makefile + `docs/lit/html-renderer.md` +
      `docs/lit/plugins.md` updated to reflect the convention-based
      lookup.
- [x] **Lower Markdown directly to JSX** (decision 3). `marklit/convert.go`
      and `marklit/reference_parser.go` now emit `ast.JSXElement` for every
      Markdown construct; `ast.Invoke` is no longer produced anywhere in
      marklit. The lowerings landed as planned:

      - `# heading` → `<Title>`; `## sub` → `<Section><Title>…</Title>…
        </Section>` (PascalCase)
      - `[text](#tag)` and `[#tag]` → `<Reference tag="…">…</Reference>`
        (PascalCase — `reference_parser.go` flipped)
      - ` ```lang ` fenced code → `<CodeBlock language="lang">…
        </CodeBlock>` (PascalCase, treehighlight)
      - ` ``` ` fenced code without language → `<pre>`
      - `*italic*` → `<em>`; `**bold**` → `<strong>`;
        `` `code` `` → `<code>`
      - `[text](url)` → `<a href="url">…</a>`;
        `![alt](src)` → `<img src="…" alt="…"/>`
      - `- list` / `1. list` → `<ul>`/`<ol>` + `<li>`
      - `| tables |` → `<table><tr><td>…`
      - `---` → `<hr/>`
      - `> quote` → `<Inset>` (kept PascalCase so the `class="inset"`
        wrapper survives; lowering to bare `<blockquote>` would have
        broken the existing CSS/test expectations)
      - Inline / block raw HTML fallback → `<RawHTML>` /
        `<RawHTML block="true">` (the latter is the new prop on the
        existing `RawHTML` builtin)

      Templates *not* yet deleted — every `Styled{Style: "X"}` shim
      (`italic.tmpl`, `bold.tmpl`, `link.tmpl`, `image.tmpl`, …) is
      still reached through the matching PascalCase styled builtin
      (`<Italic>`, `<Bold>`, `<Link>`, …). They become unreachable
      only when those builtins are retired, so the template sweep is
      deferred to the next checklist item.

      This unblocks the next three items (retire Plugin machinery,
      collapse baselit/ into builtins/, template sweep).
- [x] **Retire the Plugin machinery** (decisions 3 + 6 follow-up).
      `booklit.Plugin`, `booklit.PluginFactory`, `Section.PluginFactories`,
      `Section.Plugins`, `Section.UsePlugin`, `VisitInvoke`, and the
      `ast.Invoke` / `Invoke.Method()` types are all deleted. Threading
      collapsed accordingly: `Processor.LoadFile` / `LoadFileIn` /
      `EvaluateFile` / `EvaluateNode` (and the `SectionProcessor`
      interface) lost their `[]PluginFactory` parameter; the `Evaluate`
      stage lost its reflection helpers (`convert`, `complainL`), its
      `SlowInvokeThreshold` field, and the `--slow-invoke-threshold`
      CLI flag that fed it. `UndefinedFunctionError` /
      `FailedFunctionError` and their `undefined-function.tmpl` /
      `function-error.tmpl` templates are gone — they only existed to
      report VisitInvoke failures. `golang.org/x/text` (which only
      powered `Invoke.Method`'s kebab→PascalCase title casing) is now
      indirect-only; `go mod tidy` reflects that.
- [x] **Collapse `baselit/` into `builtins/`** (decision 3). The
      `baselit/` directory is gone. By the time it landed nothing
      called any `baselit.Plugin` method except `Syntax`, so collapsing
      meant moving the treehighlight glue
      (`Syntax`/`highlightedCodeContent`/`autoLink`/`codeWrapper`/
      `rawHTML`) into `builtins/code.go` as the `syntax` helper used by
      the existing `<CodeBlock>`/`<Syntax>` builtin and porting
      `TestSyntaxLinksBooklitComponentTags` to `builtins/code_test.go`.
      Every other `baselit.Plugin` method (Title, Italic, Bold, List,
      Table, Reference, Target, RawHtml, …) was already dead code after
      the Plugin-machinery retirement and dropped with the package.
- [x] **Delete unused render/html templates** (decision 3 follow-up).
      With Markdown converting straight to lowercase `<em>` / `<strong>`
      / `<sub>` / `<sup>`, the matching PascalCase styled builtins
      (`Italic`, `Bold`, `Subscript`, `Superscript`) and their
      `italic.tmpl` / `bold.tmpl` / `subscript.tmpl` /
      `superscript.tmpl` shims became pure duplication and are deleted.
      `<ThematicBreak/>` (which had no template at all — would error if
      ever rendered) goes with them. `StyleItalic` /  `StyleBold` /
      `StyleSubscript` / `StyleSuperscript` constants are dropped;
      the two tests that used `StyleBold` as a generic round-trip sample
      switch to `StyleLarger`. Templates that still carry inline
      styling not expressible in lowercase HTML — `larger.tmpl`,
      `smaller.tmpl`, `strike.tmpl`, `inset.tmpl`, `aside.tmpl` —
      stay, as do `link.tmpl` / `image.tmpl` / `list.tmpl` /
      `table.tmpl` which are still reached through the `Link`, `Image`,
      `List` (Definition) and contentjson paths. Test inputs that
      wrapped JSX expressions in `<Italic>` / `<Bold>` now use Markdown
      `*…*` / `**…**` where the wrapped content is plain text and
      lowercase `<em>` / `<strong>` where it's a JSX expression
      (Markdown emphasis can't span across the parser's child
      boundaries — `*{x}*` splits into three text chunks and the `*`s
      become literal). Surprise from this work: tree-sitter-html hangs
      on inputs without HTML tags. Discovered when an edited Define
      `sig='*text*'` fed `<Syntax language="html">…</Syntax>` and
      `treehighlight.Chunks` never returned; reproducible with a one-
      line unit test. The docs' Define sigs are kept as HTML-tag form
      (`<em>text</em>`, `<strong>text</strong>`, `<sup>text</sup>`,
      `<sub>text</sub>`) to dodge the hang. Filed as a follow-up; not
      part of the pivot.
- [x] **Retire the partials machinery** (decision 10).
      `Section.SetPartial`, `Section.Partial`, and the `Section.Partials`
      field are deleted; the `<SetPartial>` builtin and its `<Define
      tag="set-partial">` doc entry go with them. `stages.Collect`'s
      VisitSection no longer iterates `con.Partials`. The four
      partial-routing fixtures (`partial-template.tmpl`,
      `custom-style.tmpl`, `inline-custom-style.tmpl`, and the
      `set-partial-read-template` branches in
      `page.tmpl`/`full-styled-page.tmpl`) are deleted; the surviving
      `page.tmpl`/`full-styled-page.tmpl` collapse to their non-routing
      bodies. `Styled.Partials` (and the `booklit.Partials` type) stay —
      they carry renderer metadata (e.g. `Language` on highlighted code
      blocks), which is a separate concern from per-section partials.
      The baselit `SetPartial` method called out in the checklist was
      already dropped with baselit/.
- [ ] **HTML block-tag registry** (decision 11 prep). New
      `htmltags.go` (top-level, or `internal/htmltags` if kept
      private) exporting a `Block` set / `Block(tag string) bool`
      function over the CommonMark/HTML5 block-level tag list
      (~60 entries: `address`, `article`, `aside`, `blockquote`,
      `body`, `details`, `div`, `dl`, `dt`, `fieldset`, `figcaption`,
      `figure`, `footer`, `form`, `h1`–`h6`, `header`, `hr`,
      `iframe`, `li`, `main`, `menu`, `nav`, `ol`, `p`, `pre`,
      `section`, `summary`, `table`, `tbody`, `td`, `tfoot`, `th`,
      `thead`, `tr`, `ul`, …). Source: HTML5 spec § 4.6 +
      CommonMark HTML block tag list. Goldmark uses an equivalent
      list internally but doesn't export it; hand-maintained.
- [ ] **`RawElement` / `RawFragment` content types + migration**
      (decision 11). New `raw_element.go` / `raw_fragment.go` in
      package `booklit`. `RawElement{Tag, Attrs, Content}` with
      `IsFlow()` = `!htmltags.Block[Tag]`; implements
      `VisitRawElement`. `RawFragment{HTML}` always flow; implements
      `VisitRawFragment`. Renderer methods in `render/html.go`
      replace `raw-html.tmpl` — `RawElement` renders
      `<tag attrs>content</tag>` (self-closing for void elements),
      `RawFragment` writes its bytes verbatim. Migrate every emit
      site:
      - `stages/evaluate.go::dispatchRawHTML` returns
        `RawElement{Tag, Attrs, Content}` (the
        `Sequence{Styled{open}, children, Styled{close}}` shape
        and the `Block: node.MultiLine` plumbing all go).
      - `builtins/raw_html.go` emits `RawElement` (the `block`
        prop goes away — tag determines block-ness).
      - `builtins/code.go`'s `rawHTML()` highlighter-span helper
        emits `RawFragment`.
      - `builtins/content.go`'s `<Code>` emits
        `RawElement{Tag: "code"}` when content is flow, or
        `RawElement{Tag: "pre", Content: RawElement{Tag: "code"}}`
        when block. Same `code.IsFlow()` introspection as today;
        no new type, no template.
      - `builtins/code.go`'s `<CodeBlock>` / `<Syntax>` emit the
        same `<pre><code>` shape with `class="language-X"` on the
        inner `<code>` (Prism/highlight.js convention). The
        `Partials{Language}` indirection — set in one place, read
        by zero — is dropped entirely.
      `verbatim.tmpl`, `code-block.tmpl`, and `code-flow.tmpl`
      delete with no replacement Go visitor needed.
- [ ] **Bake stdlib components in via `go:embed` + replace
      styling shims** (decision 11). New `components/` package at
      the repo root with `embed.FS` of `*.md` files (modeled on
      `render/html/embed.go`'s `//go:embed *.tmpl`). Initial
      contents (one-liners mirroring the current `.tmpl` bodies):
      `Larger.md` (`<span style="font-size: 120%">{children}</span>`),
      `Smaller.md`, `Strike.md`, `Inset.md`, `Aside.md`. The
      component dispatch tier (third tier in
      `VisitJSXElement`) gains a fallback to `components.Assets`;
      project `components/` overrides stdlib by name. Then delete
      `builtins/styled.go` (the `styled()` factory + 5
      `Register` calls), `render/html/{larger,smaller,strike,inset,
      aside}.tmpl`, and the `StyleLarger`/`StyleSmaller`/
      `StyleStrike`/`StyleInset`/`StyleAside` constants. Rendered
      HTML stays byte-identical, so integration fixtures in
      `tests/blocks_test.go` / `tests/prose_test.go` don't move;
      only content-tree-shape unit tests update.
- [ ] **Paragraph boundary handling for block/flow JSX**
      (decision 11). Two complementary changes in
      `stages/evaluate.go::VisitParagraph` (and the dispatch
      that handles standalone JSX):
      (a) When a `Sequence` line contains a block element
      (`!Sequence.IsFlow()` with `len(para) > 1`), split the
      paragraph at that point — emit preceding lines as
      `Paragraph`, emit the block content unwrapped, continue
      paragraph-building with trailing lines. CommonMark
      behavior; replaces the invalid `<p>…<blockquote>…</p>`
      passthrough that today's `Block: node.MultiLine` flag was
      masking.
      (b) When the block JSX parser claims a node and its
      evaluated content is flow (`content.IsFlow() == true` and
      not already a `Paragraph`), wrap in
      `booklit.Paragraph{content}` so standalone
      `<Larger>x</Larger>` renders as `<p><span>x</span></p>`.
      Mirrors the `len(para) == 1 && !para[0].IsFlow()` unwrap
      logic in the inverse direction.
- [ ] **Delete `Styled`, `Partials`, contentjson updates**
      (decision 11). `styled.go` (the type, `IsFlow` override,
      remaining `StyleVerbatim` / `StyleCodeBlock` /
      `StyleCodeFlow` constants, `Style` type alias) deleted.
      `booklit.Partials` type deleted from `section.go` — fully
      unused after the previous items.
      `contentjson/contentjson.go` loses its `case Styled`
      branches, gains `case RawElement` / `case RawFragment`.
      `contentjson/wire/wire.go` loses the `Style` / `Block` /
      `Partials` `Node` fields, gains `Tag` / `Attrs` / `HTML`. The
      `encodePartials` / `decodePartials` helpers go with the
      Partials type.

### Findings (2026-06-01)

- First two items (sweep stale references, delete planning docs)
  done in earlier commits.
- Attempted "remove the `\foo{}` parser" first; reverted because it
  broke ~80 cases whose `.md` `Input` fixtures still contain
  `\foo{}`. Bulk migration ordered ahead of parser removal.
- Test-fixture migration completed (commit `da05f5f`). Markdown
  preferred; JSX only where it has to be. The parser removal (next
  unticked item) is now unblocked.
- `\foo{}` Markdown parser removed (2026-06-01). One test change of
  note: `TestBackslashEscape` for `user\\example.com` now produces a
  single `ast.String("user\example.com")` instead of two split
  strings — goldmark used to split text segments at the `\` trigger
  character our (now-gone) invoke inline parser registered for.
  Output is identical; the change is internal AST shape only.
- Last three `Ext: ".lit"` tests resolved (2026-06-01). Mid-word
  inline parser case ported to MarkDangJSX (`<Italic>` survives
  cleanly mid-word, no quirk). The two `\code{}`-variant cases were
  pure `.lit` parser surface and got deleted — no MarkDangJSX
  equivalent and fenced-code coverage already exists. PEG parser
  removal is now unblocked.
- PEG parser deleted (2026-06-01). ~2700 lines of generated code +
  `ast/booklit.peg` source + `ast/errors.go` (PEG-only error
  unpacking) all gone. `go mod tidy` also removed transitive
  `golang.org/x/mod` and `golang.org/x/tools` indirects that pigeon
  pulled in. `booklit.ParseError` is now unreferenced but kept —
  public API, no cost. `marklit.Parse` is the only source-text
  parser now.
- Components split + unified JSX parser (2026-06-01). The big design
  shift here is in marklit: the JSX parser used to trigger only on
  `<Upper...>`, leaving lowercase tags to goldmark's raw-HTML
  pipeline. That meant nested PascalCase JSX inside `<div>` blocks
  (like `docs/components/Define.md`'s `<Target>...<Syntax>...`
  pattern) was eaten by goldmark's HTMLBlockParser before our
  parser saw it. Now the JSX parser claims ANY `<letter...>`. The
  evaluator distinguishes at dispatch time: lowercase wraps the
  evaluated children in literal raw HTML opening/closing tags
  (block flag piped through from `MultiLine`); PascalCase still
  routes through builtin → Dang → template tiers. Goldmark's
  HTMLBlockParser is stripped from the default set entirely (it
  would only race against us and occasionally win). Tradeoffs:
  - Void HTML elements (`<br>`, `<hr>`, `<img>`) require an
    explicit `/>` to be recognized as JSX; otherwise they fall
    through to goldmark's inline raw-HTML parser as literal text.
    Acceptable cost for the parser unification.
  - `docs/components/Define.md` needed one edit: bare `<` and `>`
    in JSX child position now get Markdown-parsed (the `>`
    triggered the blockquote parser), so the template uses a Dang
    string concat (`{"<" + componentName(tag: tag) + ">"}`)
    instead of literal angle brackets.
  - `convert.go` lost its `interpolateRawHTML` machinery (~130
    lines) — `{expr}` interpolation inside raw HTML is no longer
    needed; expressions are real JSX children of the lowercase
    element now.
- New finding: partials are dead weight now that components exist.
  `<SetPartial>` was solving "thread named content into a renderer
  template" by stuffing entries into a per-section map keyed by
  string; the same job is done more naturally by component children
  and props. Recorded as decision 10 + checklist item above.
- Markdown → JSX lowering landed (2026-06-01). `marklit/convert.go`
  and `marklit/reference_parser.go` emit `ast.JSXElement` for every
  construct; `convertInvoke`, the `KindInvoke` dispatch, and
  `marklit/invoke_node.go` are deleted. Two unexpected side-effects:
  - `Styled.String()` used to return its content's text verbatim;
    after lowering, the `Sequence{Styled{raw-html,"<em>"}, …,
    Styled{raw-html,"</em>"}}` wrapper meant plain-text consumers
    (search index, `stringifyEverything`) started seeing raw HTML
    bytes. Fix: short-circuit `Styled.String()` for `Style ==
    "raw-html"` so markup bytes don't leak into text contexts.
    Renderer path moves from `{{ . | rawHTML }}` to `{{ .Content |
    rawHTML }}` in `raw-html.tmpl` to compensate (it was previously
    relying on `Styled.String()` to return its content).
  - Search index text for lists/tables is now flatter
    (`Item 1\n\nItem 2\n\n` instead of `* Item 1\n\n* Item 2\n\n`,
    `a1b2` instead of `| a | 1 |`). Test expectation updated;
    revisit if the search UX wants the old prefixes back.

## Where to read more

- This file is the only one. Pre-existing planning docs are slated
  for deletion (checklist item 2); their progress logs are
  recoverable from `git log` if ever needed.
