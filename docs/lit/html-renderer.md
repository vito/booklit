# <Aux>The </Aux>HTML Renderer {#html-renderer}

The presentation of your content is controlled by a *renderer*. At
present, there is only one, and it's probably the one you'll want to use: HTML,
for generating static websites.

Booklit comes with some extremely barebones templates that don't include any
styling at all. You'll probably want to change that.

The HTML renderer uses Go's built-in
[`html/template`](https://golang.org/pkg/html/template) package. To
override templates, drop your `.tmpl` files into a directory named
`html/` at your project root. Booklit walks up from `--in` looking for
a sibling `html/` directory and loads anything it finds — no flag is
needed.

A typical layout:

```bash
project/
├── lit/
│   └── index.md     # passed to --in
└── html/
    ├── page.tmpl    # override
    └── section.tmpl # override
```

```bash
booklit -i lit/index.md -o public
```

<TableOfContents/>

## Base Templates

The following template files will be executed if present in the HTML
templates directory, with the corresponding data type as `.`:

| template | type for `.` |
| --- | --- |
| <TemplateLink tmpl="page.tmpl"/> | <Godoc ref="booklit.Section"/> |
| <TemplateLink tmpl="section.tmpl"/> | <Godoc ref="booklit.Section"/> |
| <TemplateLink tmpl="link.tmpl"/> | <Godoc ref="booklit.Link"/> |
| <TemplateLink tmpl="list.tmpl"/> | <Godoc ref="booklit.List"/> |
| <TemplateLink tmpl="paragraph.tmpl"/> | <Godoc ref="booklit.Paragraph"/> |
| <TemplateLink tmpl="preformatted.tmpl"/> | <Godoc ref="booklit.Preformatted"/> |
| <TemplateLink tmpl="reference.tmpl"/> | <Godoc ref="booklit.Reference"/> |
| <TemplateLink tmpl="sequence.tmpl"/> | <Godoc ref="booklit.Sequence"/> |
| <TemplateLink tmpl="string.tmpl"/> | <Godoc ref="booklit.String"/> |
| <TemplateLink tmpl="target.tmpl"/> | <Godoc ref="booklit.Target"/> |
| <TemplateLink tmpl="toc.tmpl"/> | <Godoc ref="booklit.Section"/> |
| <TemplateLink tmpl="aside.tmpl"/> | <Godoc ref="booklit.Aside"/> |
| <TemplateLink tmpl="definitions.tmpl"/> | <Godoc ref="booklit.Definitions"/> |
| <TemplateLink tmpl="table.tmpl"/> | <Godoc ref="booklit.Table"/> |
| <TemplateLink tmpl="image.tmpl"/> | <Godoc ref="booklit.Image"/> |

The most impactful of these is `page.tmpl`, which is used for the
top-level section for each "page" rendered. This is where you would place
assets in `<head>`, for example.

## Template Functions

Booklit executes templates with the following functions available:

<Definitions>
<Definition term='{{tag | url}}'>
generate a URL for the tag
</Definition>
<Definition term='{{content | stripAux}}'>
strip [#aux] elements from the content
</Definition>
<Definition term='{{string | rawHTML}}'>
render the string as raw HTML, unescaped
</Definition>
<Definition term='{{string | rawURL}}'>
permit the rendered value to be placed in a `url=""` attribute
</Definition>
<Definition term='{{content | render}}'>
render the content
</Definition>
<Definition term='{{walkContext currentSection subSection}}'>
generate a convenience struct with fields `.Current` and `.Section`,
useful for traversing a tree of sections while retaining the "current"
section, e.g. so it can be marked as "active" in a navigation tree
</Definition>
<Definition term='{{section | headerDepth}}'>
return the number that should be used for the section's header, i.e. `<hN>`
</Definition>
</Definitions>

## Styled Content

Styled content, i.e. <Godoc ref="booklit.Styled"/>, instructs the HTML
renderer to use the `*.tmpl` template named after the style.

For example, [#bold] is implemented in the
[`baselit`](#baselit) plugin by returning:

```go
booklit.Styled{
  Style:   booklit.StyleBold, // "bold"
  Content: content,
}
```

Booklit includes a `bold.tmpl` template which is evaluated with `.`
as the `booklit.Styled` value:

```go-html-template
<strong>{{.Content | render}}</strong>
```

Thus, when content is styled with `"bold"`, it will render in
**strong tags**.

### Styles with Partials

Additional content can be propagated to the template by setting it
`Partials`:

```go
booklit.Styled{
  Style:   "my-wackadoo-style",
  Content: content,

  Partials: booklit.Partials{
    "Title": title,
  },
}
```

Then, with `my-wackadoo-style.tmpl` as the following:

```go-html-template
<div class="wack">
  <h1>{{.Partial "Title" | render}}</h1>

  {{.Content | render}}
</div>
```

This would result with `title` rendered in between the `<h1>`
tags, and `content` rendered below.

JSX components that don't have a built-in dispatch will be wrapped in
a `booklit.Styled` automatically, so dropping `html/MyComponent.tmpl`
into the templates directory is enough to introduce a new component —
no Go code required. Props are passed through as Partials keyed by the
JSX attribute name (camelCase).

## Styled Sections

Using [#styled] instructs the HTML renderer to use
`(name).tmpl` instead of `section.tmpl`, or `(name)-page.tmpl`
instead of `page.tmpl` (if it exists).

So, given the following example:

```markdown
# Fancy Section

<Styled name="fancy"/>

I'm a fancy section!

## Sub-section

I'm a normal sub-section!
```

...and the following as `fancy.tmpl` in your project's `html/`
directory:

```go-html-template
<div class="fancy">
  <em><strong>{{.Title | render}}</strong></em>

  {{.Body | render}}

  {{if not .SplitSections}}
    {{range .Children}}
      {{. | render}}
    {{end}}
  {{end}}
</div>
```

...the following will be the rendered HTML for the section:

```html
<div class="fancy">
  <em><strong>Fancy Section</strong></em>

  <p>I'm a fancy section!</p>

  <h2>Sub-section</h2>

  <p>I'm a normal sub-section!</p>
</div>
```

Note that the styling only applies to the section that declares it; it does
not propagate to its children.
