# Syntax {#syntax}

Booklit documents are Markdown files with embedded JSX-style component
invocations. Standard Markdown formatting (emphasis, bold, links, code,
lists, etc.) is supported natively, and everything else is either text or
a `<Component>` call.

<TableOfContents/>

## Prose Syntax {#prose-syntax}

Booklit builds on top of standard Markdown, so the prose rules are familiar:

- The top-level of a document is a series of *paragraphs*, separated
  by one or more blank lines.
- A *paragraph* is a series of *lines*, separated by linebreaks.
  Adjacent lines within a paragraph are joined (soft line breaks become
  spaces).
- *Emphasis* is written with `*asterisks*` and **bold** with
  `**double asterisks**`.
- *Inline code* is written with `` `backticks` ``.
- *Links* are written as `[display text](url)`.
- *Images* are written as `![alt text](path)`.
- *Headings* can be written with the `#` prefix, which maps to [#title].
  A `{#tag}` after a heading sets an explicit anchor tag.
- Standard Markdown lists (`- item`, `1. item`) and tables
  (`| a | b |`) work as expected.
- In addition to Markdown formatting, JSX-style [component
  calls](#component-syntax) can be used inline or at the block level.

## Comment Syntax {#comment-syntax}

Comments are delimited by `{-` and `-}`. They can be multi-line,
appear in between words, and they can also be nested. This makes commenting
out large blocks of content easy:

```
Hi, I'm{- a comment -} in the middle of a sentence!

{-
  I'm hanging out at the top level,

  {- being nested and stuff -}

  with multiple lines.
-}
```

## Component Syntax {#component-syntax}

Components are written in JSX style. A component invocation begins with
`<` followed by an *uppercase* letter (lowercase tags are passed through
as raw HTML, matching React's convention).

A component can be self-closing or paired:

```markdown
<TableOfContents/>

<Section>
  Hello, sub-section!
</Section>
```

### Props

Components accept named *props* as attributes. String values use double
quotes; expression values use curly braces (Dang expressions, evaluated
at build time):

```markdown
<Reference tag="getting-started"/>

<Image path="diagram.png"/>

<Card title={greeting}/>
```

For props that map to keyword-style arguments, the convention is
camelCase: `<Link target="x">y</Link>`, not `<Link Target="x">`.

### Children

Content between the opening and closing tag becomes the component's
*children*. Single-line invocations parse children as inline Markdown;
multi-line invocations parse children as block Markdown (blank lines
yield paragraphs):

```markdown
<Title>Hello, *world*!</Title>

<Section>
  <Title>Sub-section</Title>

  Body paragraph with **bold** and *italic*.

  Another paragraph.
</Section>
```

### Expressions

Anywhere inside a component — as a prop value or as a child — you can
write `{expr}` to evaluate a [Dang](https://github.com/vito/dang)
expression at build time:

```markdown
The answer is <Italic>{2 + 2}</Italic>.

<Card title={"Greetings, " + name}>...</Card>
```

The result is converted to content: strings render as text, numbers
and booleans stringify, lists flatten into a sequence. Records and
functions are not yet renderable.

For instance, the answer to `<Italic>{2 + 2}</Italic>` is
<Italic>{2 + 2}</Italic>, computed at build time.

`{expr}` only triggers *inside* a JSX element. A bare `{1+2}` in a
paragraph is just literal text — the parser starts on `<`, not `{`.

Project-wide configuration (extra GraphQL imports, Dagger module
auto-import) lives in `dang.toml` and `dagger.json` alongside the
input file, following Dang's own discovery rules.

### Lowercase tags

Lowercase tag names are treated as literal HTML and pass through to the
output unchanged:

```markdown
This is <br/> a line break.
```

If a component name collides with an HTML element name (rare, given the
PascalCase convention), the JSX parser still wins on `<UpperCase` —
because lowercase falls through to HTML, this only really comes up if
you try to name a component `<Div>` etc.
