# Dagger content: state + where to pick up

This is the working doc for "Booklit content produced by a Dagger
module." It captures what's landed, how it fits together, how to run it,
and what's left. Start here in a new session. For the wider pivot, see
`pivot.md`.

## The goal

Replace the last in-process docs built-in (`<LitSyntax>`) with a Dagger
module, and in doing so answer the real question: **how does a Dagger
module return Booklit content?** The answer we built: a JSON wire format
(`contentjson`) the module marshals and Booklit decodes back into native
`booklit.Content`.

## What's done (all committed on `jsx-dang`)

End-to-end works: a `<LitSyntax code="…"/>` in a `.md` runs a Dagger
module that highlights the code with tree-sitter and returns it as JSON;
Booklit decodes it into native content and renders it through the
existing templates. Verified by rendering a real doc.

Commits, oldest first:

- `feat(contentjson): add JSON wire format for booklit.Content`
- `feat(dangeval): decode Dagger-returned content via contentjson`
- `refactor(contentjson): extract dependency-free wire package`
- `feat(dangeval): serve the Dagger module so dependencies are callable`
- `feat(docs): replace LitSyntax built-in with a Dagger module`

## How it fits together (the round trip)

1. **Author** writes `<LitSyntax code="\reference{x}"/>`.
2. **Template** `docs/html/LitSyntax.md` is the body `{booklitdoc.litSyntax(code: code)}`.
   `<LitSyntax>` resolves to this mdx template (tier 3) — there is no
   longer a `<LitSyntax>` built-in.
3. **Dang** evaluates the `{expr}`. `booklitdoc` is the installed Dagger
   dependency; `litSyntax` is its function. The call goes to the engine.
4. **Module** (`dagger/booklitdoc`, Go SDK) highlights with its vendored
   copy of `treehighlight`, turns tree-sitter query captures for
   `\function` names into `wire.Ref` nodes, builds a `*wire.Node` tree,
   and returns `wire.Marshal(...)` as the Dagger `JSON` scalar.
5. **Bridge** (`dangeval.Evaluator.ContentFromValue`) sees a `JSON`
   scalar, runs `contentjson.Unmarshal(data, section)` → native
   `booklit.Content`, rehydrating `Reference`/`Target` against the
   current section.
6. **Renderer** renders the `Styled{lit-block}` / `code-block` /
   `raw-html` tree through the existing `.tmpl`s.

## The wire format (`contentjson`)

- `contentjson/wire` — dependency-free (`encoding/json` only): the
  `Node` union struct + builder constructors (`String`, `Seq`, `Styled`,
  `RawHTML`, `Ref`, `Target`, …) + `Marshal`/`Unmarshal`. Producers
  (the module) import only this.
- `contentjson` — imports `booklit` + `wire`; `Marshal(Content)` /
  `Unmarshal(data, *Section)` translate between native content and
  `wire.Node`.
- Serializable: String, Sequence, Paragraph, Preformatted, Styled
  (+Partials), Link, Image, List, Table, Definitions, Aux, Reference,
  Target.
- **Not** serializable (errors from `Marshal`, even nested): `Section`,
  `TableOfContents`, `Lazy` — they're bound to live evaluator state.
- `Reference`/`Target` cross as just a tag name; `Unmarshal` re-attaches
  the passed `*Section`. That's how the `\function` links survive.

## The runtime-serve fix (the subtle part)

`dang.ResolveDaggerImport` only *introspects* a module's schema (enough
to type-check `booklitdoc.litSyntax(...)`), but the runtime client is a
bare `dagger session` that never *serves* the module — so the call
type-checks yet fails at execution with `Cannot query field "booklitdoc"
on type "Query"`. `dangeval.New` now calls `dang.DaggerServeModule(...,
includeDependencies: true)` so the dependency lands on the live `Query`.
This is reusable: any Dagger-dependency call from `{expr}` now works,
not just `litSyntax`.

## The module + its local dependency

`dagger/booklitdoc` is a Go-SDK module with local copies of
`contentjson/wire` and `treehighlight`. The copy is deliberate: the Go
SDK runtime builds with cgo disabled, while tree-sitter's Go bindings
need cgo, so `LitSyntax` shells out to `cmd/lit-syntax` in a
`golang:1.26` container with `CGO_ENABLED=1`.

The module is installed into the root module with
`dagger install ./dagger/booklitdoc` (root `dagger.json` has a
`dependencies` entry). Generated SDK code (`dagger.gen.go`, `internal/`)
is gitignored; **after a fresh checkout run `dagger develop` in
`dagger/booklitdoc`** before the module will build.

## How to run / iterate

- Call the module in isolation:
  `cd dagger/booklitdoc && dagger call lit-syntax --code='\title{x}'`
  (prints the contentjson).
- Regenerate after editing the module: `dagger develop` in
  `dagger/booklitdoc`; if you change the dependency wiring, also
  `dagger develop` at the repo root.
- Render a doc end-to-end (starts the engine, builds + serves the
  module):
  `go run ./cmd/booklit -i docs/lit/<file>.md -o /tmp/out --html-templates docs/html`
- Unit tests (no engine needed):
  `go test ./contentjson/... ./dangeval/ ./treehighlight/...`
- Full site: `make docs/outputs/index.html` (or `scripts/build-docs`).

## Open follow-ups (roughly prioritized)

1. **`JSONValue!` variant.** The bridge already handles a `JSONValue!`
   return (forces `.contents`), but the module's pinned Go SDK
   (`v0.20.7-…`) predates `dag.JSONValue()`, so the lazy variant was
   dropped. Bump the SDK, re-add `LitSyntaxValue`, and exercise the
   `GraphQLValue{TypeName:"JSONValue"}` path in `ContentFromValue`.
2. **Raw-source children.** `<LitSyntax>` takes `code` as an attribute
   because template/JSX children arrive as *rendered* content, not
   source text. Add a way to get the unrendered body so authors can
   write a fenced code block instead of a one-line attribute. This is
   the main ergonomics gap.
3. **Verify link rendering end to end.** The module emits `ref` nodes
   and the bridge rehydrates them (covered by unit tests + the module's
   raw output), but a doc that actually `\define`s the target tags is
   needed to confirm `<a href>` output. `<LitSyntax>` is currently
   unused by the real docs.
4. **Fix `decisions.md`.** Its claim that `{build(...)}` (a root-module
   function) is callable from docs is wrong: introspection exposes the
   module's dependencies + core API on `Query`, not its own functions.

## Latest update: tree-sitter highlighter

`cmd/booklit-docs` and `docs/booklitdoc` are gone. `baselit.Syntax` now
uses `treehighlight`, a tree-sitter renderer with inline styles, so
fenced blocks no longer depend on a process-wide chroma `styles.Fallback`.
`<LitSyntax>` uses a vendored copy through the Dagger module and emits
references from tree-sitter query captures rather than a regex over
highlighted HTML. If Booklit is built with `CGO_ENABLED=0`,
`treehighlight` compiles to an escaped plain-code fallback.

## Gotchas worth remembering

- `{expr}` is only parsed inside JSX elements and templates, never in
  bare Markdown prose — that's why the invocation lives in a template.
- Dispatch order is built-in → Dang → template. A leftover built-in
  registration will shadow a same-named template (this bit us: the old
  `LitSyntax` built-in hid the new template until it was removed).
- A `JSON`-typed scalar return is the signal "this is content, decode
  it"; ordinary `String` stays text. Detection is by the scalar's type
  name (`*dang.Module` `.Named == "JSON"`).
