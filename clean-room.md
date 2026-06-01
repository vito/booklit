# Booklit

A static site generator from semantic documents. Markdown for prose,
JSX for structured invocations, Dang for expressions, content authored
against a recursive `Section` tree, rendered to HTML.

## File format: MarkDangJSX

One format, one parser, used everywhere — content files in `lit/`,
component templates in `components/`, prose inside JSX bodies. Three
layers in one syntax:

- **Markdown prose** (CommonMark + GFM tables) — headings, lists,
  emphasis, links, fenced code, blockquotes.
- **JSX-style elements** for structured invocations and raw HTML —
  `<Section>`, `<Larger>`, `<div>`, `<span>`.
- **Dang expressions** in `{…}` for values from outside the static
  text.

Markdown features keep working as expected: `# heading` auto-creates a
section, `*em*` → `<em>`, fenced code goes through the syntax
highlighter, etc.

### JSX tags

**PascalCase** tags dispatch through the JSX tier system below.
Unknown PascalCase names are a build error — no implicit styling
fallback.

**lowercase** tags pass through as raw HTML wrappers. Both kinds of
tag can mix freely:

```html
<div class="card">
  <Title>Hello</Title>
  {description}
</div>
```

### Multi-line vs inline

The parser distinguishes JSX that spans multiple source lines from
single-line JSX. A JSX child on its own source line is treated as a
block standalone (its own item in the surrounding content). A JSX
child sharing a line with prose merges into the surrounding paragraph.
This matches an author's intuition: `<Larger>x</Larger>` inline within
prose stays inline; `<Larger>x</Larger>` on its own line becomes its
own block.

### Attributes

- `name="literal string"` — string attribute, spliced verbatim
  (HTML-escaped at render time).
- `name={Dang expression}` — expression attribute, evaluated against
  the Dang env and stringified.

camelCase is preserved.

### Dang interpolation

Wherever content can appear, `{expression}` evaluates against the Dang
env. The result is bridged to Booklit content via `dangeval/bridge.go`
(string → `String`, list → `Sequence`, null → empty, content-typed
return → its content verbatim, JSON-typed Dagger return → decoded
through `contentjson`).

### Markdown inside JSX

Markdown is parsed inside *both* lowercase and PascalCase JSX bodies.
`<div>**bold**</div>` produces `<div><strong>bold</strong></div>`.
This deviates from MDX, which treats lowercase-tag content as raw
HTML — Booklit keeps the format uniform.

## Project layout

```
project/
├── lit/                    # Content: *.md (and helpers.dang next to them)
│   ├── index.md            # --in target; entry point
│   ├── other-page.md
│   └── helpers.dang        # Auto-loaded by dangeval (non-recursive scan)
├── components/             # User JSX components: <Name/> → components/Name.md
│   ├── Card.md             # Same MarkDangJSX format as content files
│   └── Define.md
├── html/                   # Renderer overrides for framework templates
│   ├── page.tmpl           # (page, section, sidebar, …). Optional;
│   ├── section.tmpl        # most projects don't need this directory.
│   └── sidebar.tmpl
├── css/                    # Static assets (CSS, fonts, images);
├── favicon.ico             # copied to the output dir by the renderer.
├── dagger.json             # Dagger module metadata (sdk = "dang")
├── .dagger/                # The project's own Dagger module (CI, build,
│   └── main.dang           # test). Its functions are NOT in {expr} scope —
│                           # only its *dependencies* + core API are.
└── dist/                   # --out target; generated HTML
```

Most projects only need `lit/` plus `components/`. `html/` is for
overriding renderer-internal templates (page layout, section layout).
`dagger.json` opts into Dagger reachability from `{expr}`.

### Flags

```
booklit -i lit/index.md -o dist
```

`--in` (`-i`) and `--out` (`-o`) are the only flags most invocations
need. Everything else (components dir, html overrides, helpers
modules, Dagger session) is discovered by walking up from
`filepath.Dir(--in)`.

## Section tree

A document is a tree of `*booklit.Section`. The top section's title
comes from the document's `# H1` heading; deeper headings (`##`,
`###`, …) create sub-sections automatically. Sections carry their
`Title`, `Body`, `Tag`, `Style` (an optional renderer template
override), and `Children`.

`<Section>`, `<Title>`, `<IncludeSection>`, `<SplitSections>`,
`<SinglePage>`, `<TableOfContents>`,
`<OmitChildrenFromTableOfContents>`, `<Styled>` (a builtin that
overrides a section's renderer template) all manipulate the section
tree.

Headings can carry an explicit anchor: `## My Heading {#my-tag}` →
section with `Tag: "my-tag"`. Without one, the tag is the slugified
title.

## JSX dispatch

`<Foo>` resolves through three tiers in order. An unknown PascalCase
name is a build error.

### Tier 1: built-in

A Go function registered in `builtins/`. There are 23 of them:

```
Aux, Children, Code, CodeBlock, Definition, Definitions, For, If,
Image, IncludeSection, Link, OmitChildrenFromTableOfContents,
RawHTML, Reference, Section, SinglePage, SplitSections, Styled,
Syntax, TableOfContents, Target, Title, Unless
```

These are the language primitives and aren't user-extensible.
Built-ins receive `(ctx, props, children)` and return
`booklit.Content`.

### Tier 2: Dang function

A `pub PascalCase(...)` callable in the Dang env. The dispatcher
passes JSX props as named args and wraps the JSX children as a
`&body` block. Each call to `body(...)` from inside the Dang function
pushes its named args into scope and re-evaluates the children —
useful for closures over per-project data and for parametric control
flow that needs the children.

### Tier 3: MarkDangJSX component template

A `<dir>/Foo.md` file. Props are bound in Dang scope by name; the JSX
children's evaluated content is bound as `children` (a
`dangeval.ContentValue` so nested styling survives). Both `{children}`
and `<Children/>` emit it.

Project `components/` takes precedence by name over the embedded
stdlib. The stdlib (`components/`) ships five styling components:
`Larger`, `Smaller`, `Strike`, `Inset`, `Aside`. Each is a one-line
MarkDangJSX template, e.g. `Larger.md`:

```html
<span style="font-size: 120%">{children}</span>
```

A project can override any stdlib component by dropping a same-named
`.md` into its own `components/`.

## Lowercase HTML

Lowercase JSX produces `booklit.RawElement{Tag, Attrs, Content}`.
- String props splice in as `name="value"` (HTML-escaped).
- Expression props evaluate against the Dang env and stringify.
- Empty Content emits self-closing (`<tag attrs/>`).
- Attributes render in **alphabetical order** for determinism.

`<RawHTML>...</RawHTML>` is an escape hatch that emits its body's
text as a `RawFragment` — pre-rendered HTML bytes, untouched.

## Block vs flow

`internal/htmltags.Block` is the source of truth for which HTML
element names are block-level — 68 entries covering CommonMark HTML
block types 1 and 6, plus `pre`/`script`/`style`/`textarea`.

`Content.IsFlow() bool` is on every content type:

- `RawElement.IsFlow()` returns `!htmltags.Block[Tag]` *and* checks
  that any inner Content is itself flow. A flow-tagged element
  wrapping block content (`<span>` around `<p>`) is treated as
  effectively block so the surrounding paragraph layout doesn't
  sandwich it in another `<p>`.
- `RawFragment.IsFlow()` is always true.
- Block content types (`Paragraph`, `Section`, `List`, `Table`,
  `Preformatted`, `Definitions`, `TableOfContents`, `Lazy`) return
  false; flow types (`String`, `Link`, `Image`, `Reference`,
  `Target`, `Aux`) return true.

### Paragraph segmentation

`stages/evaluate.go::VisitParagraph` segments each evaluated line
into runs of flow content (each wrapped in a `Paragraph`) interleaved
with block content (emitted unwrapped):

- A Sequence line whose items are mixed flow + block (an inline JSX
  evaluating to block content mid-prose) splits inside the line at
  the block boundary. CommonMark-style block HTML in paragraph
  behavior.
- Flow content whose `.String()` is empty (Target's bare anchor, an
  empty Sequence, a self-closing void RawElement) emits unwrapped.
  It's a side-effect marker, not a paragraph's worth of content.
- Standalone block-claimed flow JSX (e.g. `<Larger>x</Larger>` on its
  own line) wraps in a `<p>` so it renders as
  `<p><span>x</span></p>`, mirroring the inverse of the unwrap rule.

## Embedded Dang interpreter

`dangeval/` wraps `github.com/vito/dang` so `{expr}` evaluates as
real Dang code.

`dangeval.New` walks up from `--in` looking for:

- `dang.toml` — Dang's GraphQL imports + Dagger session config.
- `dagger.json` — a local Dagger module (see below).

It also scans the project directory for `*.dang` files (non-
recursive), treats them as one module, and merges their forms into
the type + value envs. `lit/helpers.dang` is the canonical place for
project-specific helpers.

`Evaluator` is one per build session, single-threaded.

### Dang ↔ Content bridge

`dangeval/bridge.go` maps Dang values to `booklit.Content`:

| Dang value     | Content                                |
|----------------|----------------------------------------|
| `StringValue`  | `String`                               |
| `IntValue`     | stringified `String`                   |
| `FloatValue`   | stringified `String`                   |
| `BoolValue`    | stringified `String`                   |
| `ListValue`    | `Sequence`                             |
| `NullValue`    | empty                                  |
| `ContentValue` | its carried `Content` verbatim         |
| Dagger `JSON`  | decoded via `contentjson` (see below)  |

`Evaluator.ContentFromValue` is the richer, section-aware path used
for JSX `{expr}` results — on top of the above it rehydrates
`Reference`/`Target` nodes against the current section.

## Dagger integration

Any project with a `dagger.json` gets Dagger reachability from
`{expr}` for free. `dangeval.New` finds the `dagger.json`,
introspects the served module's schema for type checking, and serves
it into the session so its dependencies are callable at runtime.

What an introspected module exposes on the session `Query` is its
**dependencies + core API**, not the module's own functions. If your
`.dagger/` module depends on `booklitdoc`, you reach the highlighter
as `booklitdoc.litSyntax(...)`, not as a root-module function.

There is no `<Foo from="..."/>` JSX-from-Dagger-module syntax and
won't be. Dagger calls go through `{expr}` like any other Dang code.

### Content from a Dagger module

A Dagger module can return Booklit content by encoding it in the
`contentjson` wire format and returning it as Dagger's `JSON` scalar.
Booklit recognizes a `JSON`-typed return in `ContentFromValue` and
decodes it back into native content. `JSONValue!` returns also work
(the bridge forces `.contents`), letting a module compose results
lazily before Booklit materializes them.

The wire schema lives in `contentjson/wire/`, dependency-free, so a
producer doesn't need to import `booklit`:

```go
import "github.com/vito/booklit/contentjson/wire"

return wire.Para(
    wire.String("hello "),
    wire.Element("strong", "", wire.String("world")),
)
```

In-process-only content (`Section`, `TableOfContents`, `Lazy`)
errors from `Marshal`. Stateful-but-nameable content
(`Reference`, `Target`) crosses carrying only a tag name and is
re-bound to the live section on decode, so cross-references survive
the round trip.

## Syntax highlighting: treehighlight

`treehighlight/` is a thin wrapper over tree-sitter that ships its
own Booklit grammar.

`<CodeBlock language="X">` / `<Syntax language="X">` runs the source
through the tree-sitter parser, groups the resulting chunks into
spans, and emits real `booklit.Reference` nodes for captures whose
kebab-case form matches an existing tag — fenced examples of
`<IncludeSection/>` automatically linkify to their definitions.

Output shape is the Prism/highlight.js convention:

```html
<pre style="…"><code class="language-X">…spans…</code></pre>   (block)
<code style="…" class="language-X">…spans…</code>              (inline)
```

cgo is the production path. With `CGO_ENABLED=0`, `treehighlight`
compiles to an escaped plain-code fallback (no spans, no links).

## Rendering

`render.HTMLEngine` walks the content tree, calling Visit methods.
Most content types render via a Go `html/template` from the embedded
`render/html/*.tmpl` set:

```
page, section, sequence, paragraph, preformatted, string,
link, image, list, table, definitions, reference, target, toc
```

A project's `html/` directory (looked up next to `--in`) overrides
any of these by file name.

`RawElement` and `RawFragment` skip the template layer: the renderer
composes `<tag attrs>...</tag>` (or self-closing) directly, and
fragments write their bytes verbatim. Both go through an
`engine.direct []byte` escape hatch that `render()` checks first,
bypassing template execution for content that's already final HTML.

Section-level template overrides: a section with `Style: "foo"`
renders through `foo.tmpl` instead of `section.tmpl`. `<Styled
name="foo"/>` sets `Style` on the current section.

## Content tree

```go
type Content interface {
    String() string         // plain-text projection
    IsFlow() bool           // false → block, true → flow (inline)
    Visit(Visitor) error
}
```

The concrete types:

| Type                                  | Block/flow | Notes                                      |
|---------------------------------------|------------|--------------------------------------------|
| `String`                              | flow       | Plain text                                 |
| `Sequence []Content`                  | derived    | Flow iff all items flow                    |
| `Paragraph []Sequence`                | block      | Each Sequence is a "line"; rendered as `<p>` joining lines with spaces |
| `Preformatted []Sequence`             | block      | Like Paragraph but line-preserving         |
| `Section *Section`                    | block      | Recursive section tree node                |
| `RawElement{Tag, Attrs, Content}`     | derived    | Lowercase JSX + the code builtins          |
| `RawFragment{HTML string}`            | flow       | Pre-rendered HTML bytes                    |
| `Link{Target, Content}`               | flow       | Hyperlink                                  |
| `Image{Path, Description}`            | flow       | `<img>`                                    |
| `List{Ordered bool, Items}`           | block      | `<ul>`/`<ol>`                              |
| `Table{Rows [][]Content}`             | block      | GFM table                                  |
| `Definitions []Definition`            | block      | `<dl>`/`<dt>`/`<dd>`                       |
| `Reference{TagName, Content, …}`      | flow       | Hyperlink to a tag, resolved at render time |
| `Target{TagName, Title, Content}`     | flow       | Anchor that registers a tag; renders `<a id></a>` |
| `Aux{Content}`                        | flow       | Auxiliary content stripped in references/ToC |
| `TableOfContents{Section}`            | block      | Generated from the section it's in         |
| `Lazy`                                | block      | Defers evaluation until rendered           |

The `Visitor` interface has one method per concrete type. Four
implementations: `stages.Collect` (registers Target anchors on
sections), `stages.Evaluate` (the JSX evaluator), the HTML renderer,
the text renderer, and `aux_.go`'s `stripAuxVisitor`.

## Building

```sh
booklit -i lit/index.md -o dist
```

The renderer reads `lit/`, finds `components/` and `html/` next to
it, bootstraps Dang (honoring `dang.toml` + `dagger.json` walking up
from `lit/`), evaluates each `.md` file, and writes HTML pages to
`dist/`. Static assets next to `lit/` (CSS, images) copy through.
