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

## Standard project layout

What a Booklit project actually looks like today, taken from the docs
site. None of these paths are hard-coded — they're conventions, and the
relevant ones can be moved with CLI flags. But this is the layout the
code is written *around*:

```
project/
├── lit/                      # Content: *.md (and helpers.dang next to them)
│   ├── index.md              # `--in` target; entry point
│   ├── other-page.md
│   └── helpers.dang          # Auto-loaded by dangeval (non-recursive scan)
├── html/                     # Templates dir, passed via --html-templates:
│   ├── page.tmpl             #   • Go html/template overrides for renderer-
│   ├── section.tmpl          #     internal templates (page, section,
│   ├── sidebar.tmpl          #     sidebar, …) — framework infrastructure
│   ├── Card.md               #   • mdx JSX-component templates (PascalCase)
│   └── Define.md
├── css/                      # Static assets (CSS, fonts, images);
├── favicon.ico               # the renderer copies these to the output dir
├── dagger.json               # Dagger module metadata (sdk = "dang")
├── .dagger/                  # The project's own Dagger module (CI, build,
│   └── main.dang             # test); its functions are NOT in {expr} scope —
│                             # only its *dependencies* + core API are.
└── dist/                     # `--out` target; generated HTML
```

A few load-bearing details and rough edges:

- **The `--in` flag's directory anchors everything Dang-side.**
  `dangeval.New` walks up from `filepath.Dir(--in)` looking for
  `dang.toml` and `dagger.json`, and scans that same directory
  (non-recursive) for `*.dang` files. That's why `helpers.dang` lives
  inside `lit/` next to the content, not at the project root. There is
  no `dang/` source dir convention.
- **`--html-templates` does double duty.** The same flag (and the same
  directory) feeds both the Go html/template engine's override loader
  (`render/html.go`'s `engine.LoadTemplates`) AND the mdx-template
  registry (`templates.New(dir)`). The engine picks `.tmpl` files; the
  registry picks `.md` files. They share the directory because the
  flag predates the mdx work.
- **`html/` is misnamed.** It contains JSX-component definitions
  (`Define.md`, `Card.md`) — neither HTML nor Go templates. Renaming to
  `components/` (or splitting `templates/` vs `components/`) would
  match what's in there.
- **Renderer-internal `.tmpl` files stay Go templates.** The 28 files
  in `render/html/` (page, section, styled, list, table, code-block,
  …) are embedded at compile time and form the rendering substrate
  built-ins emit `Styled{Style: …}` against. Phase 3b-4 considered
  converting them to mdx and recommended deferring; that recommendation
  still holds.

## Cleanup candidates

Things that look ready to remove now that the pivot has settled.
Listed in order of "safest first":

1. **The PEG `.lit` parser** (`ast/booklit.peg`, `ast/booklit.peg.go`,
   the pigeon Makefile target, and the `.lit` branch in
   `load/processor.go` lines 119–144). The repo contains zero `.lit`
   files on disk. The PEG parser only fires for three test cases in
   `tests/prose_test.go` (`Ext: ".lit"`) and a few `.lit` fixture
   strings in `tests/sections_test.go` that exercise
   `<IncludeSection>`. Rewrite those fixtures as `.md` and the whole
   ~2700 lines of PEG can go, along with the `pigeon` tooling
   dependency.
2. **The user-facing `\foo{}` syntax in Markdown.** `marklit` still
   registers `NewInvokeInlineParser` so an author can write `\italic{x}`
   inside a `.md` file. After the docs migration, no in-tree content
   uses that form — every `\foo{}` instance in `marklit/convert.go`
   is *internal* lowering (heading → `\title`, fence → `\code-block`,
   etc., emitted as `ast.Invoke` for the reflection-dispatch path).
   Removing the inline-parser registration would shrink the public
   surface to "JSX + Markdown only" without touching the internal
   lowering.
3. **`marklit/convert.go` lowering to `ast.Invoke` → lowering directly
   to `ast.JSXElement`.** Once the user-facing `\foo{}` parser is gone,
   `ast.Invoke` exists only as an internal intermediate. Lowering
   Markdown straight to JSX nodes would let us delete `VisitInvoke`'s
   reflection dispatch, `Section.PluginFactories` / `Section.Plugins` /
   `UsePlugin`, the `booklit.Plugin` / `PluginFactory` types, and
   reduce `baselit/` to a set of `builtins/`-style functions that
   register the language primitives in the same registry as
   `<Title>`/`<Section>`/etc. This is the biggest cleanup on the list
   — and the most clarifying: "plugin" stops being a misleading term
   for "internal language primitive."
4. **Stale planning docs at the project root.** `jsx-dang.md` (50K),
   `phase-3b.md` (30K), `decisions.md` (7K), and `dagger-content.md`
   (1.6K) are all append-only historical logs from the pivot work.
   Their load-bearing content is now in this file; their progress logs
   are recoverable from `git log`. Recommendation: delete them, keep
   `pivot.md` as the authoritative status doc, and trust git history
   for the journey. (Leaving this for explicit confirmation — the
   user-facing "Where to read more" in the previous revision linked to
   them, so they may still be wanted as a paper trail.)
5. **`docs/booklit-docs/` and `docs/booklitdoc/` references** in any
   remaining tooling or scripts (the directories themselves are already
   gone). The Makefile already builds via `go run ./cmd/booklit`, so
   nothing should still point at the deleted binary, but it's worth a
   sweep.

## Open questions (the come-to-Jesus list)

Things that could undermine the pivot if their answers turn out to be
bad. Listed in rough order of how load-bearing the resolution is.

### 1. Dagger-as-a-plugin-tier: what is the wire shape?

The plan always said "Dagger functions are the fourth dispatch tier"
(`<Foo from="github.com/.../mod"/>`). Today that tier doesn't exist.
Dagger functions are reachable only through `{moduleFn(...)}` inside
Dang `{expr}` — and only the introspected module's *dependencies +
core API* land on the session `Query`, not the module's own
functions. This means:

- A project-owned helper in `.dagger/main.dang` can't be called
  directly; you have to either expose it as a Dagger dependency or
  write a Dang-side wrapper.
- A one-off remote module needs a `dang.toml` import entry, not a
  Booklit-level syntax.

The deferred design question is the JSX-tag form: `<Foo from="..."/>`
maps to *what call shape*? Options the team has talked through but not
resolved:

- (a) Pure tag sugar over `{moduleFn(props)}`. Module gets a single
  JSON object of props, returns `JSON` content. Cheapest to implement
  but limited.
- (b) A dedicated tier that constructs the module's import lazily,
  dispatches by tag name, and bridges children somehow (see Q2).
- (c) Skip the tag form entirely and lean on `{expr}` + Dang helpers.

### 2. Body/children passing to Dagger functions.

Tier-2 (built-ins) and tier-3 (Dang) both give the component
re-invocable access to its JSX children:

- A Go built-in receives `node.Children` directly.
- A Dang component receives a `&body(...)` block that, when called,
  pushes its named args into Dang scope and re-evaluates the children;
  iteration via `<For>` works because each `body(item: x)` re-renders
  the JSX in scope of that `x`.

A Dagger function call is process-boundary-crossing. There is no
obvious analogue. Three sketched approaches, none yet validated:

- (a) **Eager pre-render** — render children once into `contentjson`,
  pass as a `body` arg. The Dagger function gets HTML-shaped content
  and can splice it into its own output. Loses the per-iteration
  re-evaluation that makes `<For>` work.
- (b) **Callback channel** — boot a tiny Booklit-side service for the
  duration of the call; the Dagger function calls back with
  `body(args)` invocations that re-render in Booklit. State-heavy,
  complicates session lifecycle.
- (c) **Source-text body** — pass the children's *unevaluated* source
  text and let the Dagger function treat it however it wants
  (templating, code highlighting, …). Cheap, but it punts evaluation
  to the module.

The user-asked phrasing — *"like a Dang block arg, `<Foo>{bar}</Foo>`,
that the Dagger function can invoke"* — is closest to (b). Whether
that's worth the plumbing depends on use cases. None are concrete
today; the `<LitSyntax>` experiment used the source-text-via-prop
pattern (which inspired (c)) before being deleted.

### 3. Is `contentjson` worth keeping with no consumer?

`contentjson/` + `contentjson/wire/` are ~440 lines of speculative
infrastructure: the only producer was `<LitSyntax>`, which is gone.
The bridge in `dangeval.Evaluator.ContentFromValue` still recognizes
the `JSON` / `JSONValue!` shapes and decodes them, so the runtime cost
is small — but the maintenance cost grows the longer the format sits
without a real user (the wire schema can't evolve confidently without
something exercising it).

Three options:

- Keep as-is; the next out-of-process producer (a Dagger plugin
  ecosystem) needs it.
- Delete `contentjson/` and `contentjson/wire/`, simplify the bridge
  to refuse `JSON`-typed returns. Restore when a real consumer appears.
- Keep but document the constraint ("Section/TableOfContents/Lazy
  can't cross") prominently, and write at least one example module so
  the API is exercised.

This question feeds into Q1: if Dagger-as-a-tier ships as option (a)
above, `contentjson` is the natural return type. If it ships as option
(c), `contentjson` might be unnecessary.

### 4. mdx templates have *no* Markdown support.

By design (per the phase-3b 2026-05-31 entry): templates are HTML
scaffolding + JSX + `{expr}`, no `#`/`-`/`|`/fenced-code. The
trade-off is on purpose — templates are layout, not prose. But the
ceiling shows up fast: anything beyond a trivial wrapper has to either
take rendered content as a `children` block or reach for raw HTML in
the template body. If template authors start writing `<ul><li>...</li>
</ul>` by hand to dodge Markdown's absence, we'll know it's wrong.

No urgent action; flag it.

### 5. `html/` mixes two unrelated systems under one flag.

`--html-templates docs/html` feeds the *same* directory to two
different loaders: `engine.LoadTemplates` (Go templates for renderer
overrides) and `templates.Registry` (mdx for JSX components). The
directory ends up holding `page.tmpl` next to `Define.md` next to
`Card.md`, and you can't tell from the layout which file does what.

Cleanups, in order of disruption:

- Split into two dirs (`overrides/` for `.tmpl`, `components/` for
  `.md`), introduce a `--components` flag, deprecate `--html-templates`
  for JSX.
- Or accept the conflation but document it explicitly.

### 6. `Section.PluginFactories` survives as language plumbing only.

The "Plugin" concept that the pivot was meant to eliminate
user-facing is gone (no more compiled plugins). What survives is an
internal dispatch path: Markdown lowering emits `ast.Invoke` →
reflection dispatch on `Section.Plugins` → `baselit.Plugin` methods.
External users never see it. Cleanup candidate 3 above is the
follow-through: once Markdown lowers directly to JSX, this whole
machinery — `Plugin`, `PluginFactory`, `Section.PluginFactories`,
`Section.Plugins`, `UsePlugin`, the reflection table in
`stages/evaluate.go::VisitInvoke` — can collapse into the
`builtins/` registry. Until then it's vestigial but harmless.

## Where to read more

- `git log --oneline master..HEAD` is now the most accurate journey;
  the pre-existing append-only planning docs (`jsx-dang.md`,
  `phase-3b.md`, `decisions.md`, `dagger-content.md`) are kept around
  as historical paper trail but their load-bearing claims have been
  folded into this file. If they're still present after the cleanup
  sweep, treat their progress-log sections as authoritative for "what
  was decided at the time" and this file as authoritative for "what
  is true now."
