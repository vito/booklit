# Syntax {#syntax}

Booklit documents are Markdown files with embedded JSX-style component
invocations. Standard Markdown formatting (emphasis, bold, links, code,
lists, etc.) is supported natively, and everything else is either text or
a `<Component>` call.

<TableOfContents/>

## Prose Syntax {#prose-syntax}

Booklit builds on top of standard Markdown, so the prose rules are familiar:

<List>
<Item>
The top-level of a document is a series of *paragraphs*, separated
by one or more blank lines.
</Item>
<Item>
A *paragraph* is a series of *lines*, separated by linebreaks.
Adjacent lines within a paragraph are joined (soft line breaks become
spaces).
</Item>
<Item>
*Emphasis* is written with `*asterisks*` and **bold** with
`**double asterisks**`.
</Item>
<Item>
*Inline code* is written with `` `backticks` ``.
</Item>
<Item>
*Links* are written as `[display text](url)`.
</Item>
<Item>
*Images* are written as `![alt text](path)`.
</Item>
<Item>
*Headings* can be written with the `#` prefix, which maps to [#title].
A `{#tag}` after a heading sets an explicit anchor tag.
</Item>
<Item>
In addition to Markdown formatting, JSX-style [component
calls](#component-syntax) can be used inline or at the block level.
</Item>
</List>

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

```jsx
<TableOfContents/>

<Section>
  Hello, sub-section!
</Section>
```

### Props

Components accept named *props* as attributes. String values use double
quotes; expression values use curly braces (Dang expressions, evaluated
at build time):

```jsx
<Reference tag="getting-started"/>

<Image path="diagram.png"/>

<List items={primitiveTypes}/>
```

For props that map to keyword-style arguments, the convention is
camelCase: `<Link target="x">y</Link>`, not `<Link Target="x">`.

### Children

Content between the opening and closing tag becomes the component's
*children*. Single-line invocations parse children as inline Markdown;
multi-line invocations parse children as block Markdown (blank lines
yield paragraphs):

```jsx
<Title>Hello, *world*!</Title>

<Section>
  <Title>Sub-section</Title>

  Body paragraph with **bold** and *italic*.

  Another paragraph.
</Section>
```

### Lowercase tags

Lowercase tag names are treated as literal HTML and pass through to the
output unchanged:

```jsx
This is <br/> a line break.
```

If a component name collides with an HTML element name (rare, given the
PascalCase convention), the JSX parser still wins on `<UpperCase` —
because lowercase falls through to HTML, this only really comes up if
you try to name a component `<Div>` etc.
