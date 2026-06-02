# Basic Components {#baselit}

Booklit comes with a default set of components — collectively, *baselit* —
that all sections can use without any setup. baselit provides the
structural primitives (sections, tables of contents, references) and the
common prose helpers (italic, bold, code blocks, lists, tables).

Most baselit components have a Markdown equivalent and you should reach
for the Markdown form first; the JSX components are there for the cases
Markdown can't express (custom tags, partials, sectioning beyond
headings).

<TableOfContents/>

## Defining Sections {#sections}

<Define tag="title" sig='<Title tag="optional">content</Title>'>
Set the title of the current section, optionally giving it an explicit
*tag* for cross-referencing. In Markdown form, `# heading {#tag}` does
the same thing — the JSX form is mainly useful when you need a title
that isn't derived from a `#` heading.

If no tag is specified, the section's tag defaults to a sanitized form of
the title (e.g. *I'm a fancy title!* becomes `im-a-fancy-title`).
</Define>

<Define tag="aux" sig='<Aux>content</Aux>'>
Denote auxiliary content which can be stripped in certain contexts without
losing meaning.

Used within a title declaration to provide content that will show up on the
section page itself but will be omitted when referencing the section. This
is handy for sub-titles that you don't care to show anywhere but the
section's page itself.
</Define>

<Define tag="section" sig='<Section>content</Section>'>
Introduce a sub-section. Each sub-section should conventionally begin
with a [#title] (or a `#` heading inside).

You usually don't need `<Section>` explicitly — Markdown headings already
create sections. Use it when you need a sub-section that isn't introduced
by a heading.
</Define>

<Define tag="include-section" sig='<IncludeSection path="..."/>'>
Load the Booklit document located at *path* (relative to the current
section's file) and set it as a child section.
</Define>

<Define tag="split-sections" sig='<SplitSections/>'>
Configure the renderer to generate a separate page for each immediate
sub-section rather than inlining them under smaller headings.
</Define>

<Define tag="single-page" sig='<SinglePage/>'>
When declared in a section, it overrules any [#split-sections] in
the section and any child sections (recursively) in order to force them all
on to one page. Each section's header sizing is preserved, however.
</Define>

<Define tag="table-of-contents" sig='<TableOfContents/>'>
Generate a block element that displays the table of contents from this
section downward upon rendering. Often used in combination with
[#split-sections].
</Define>

<Define tag="omit-children-from-table-of-contents" sig='<OmitChildrenFromTableOfContents/>'>
Configure the section to omit its children from table-of-contents listings.
Appropriate when the sub-sections within a section are brief and meant to
be consumed all at once.
</Define>

## Customizing Sections

<Define tag="styled" sig='<Styled name="..."/>'>
Set the section's template style to *name*. The renderer may then use
this to present the section in a different way. See [#styled-sections].
</Define>

## Links & References

<Define tag="link" sig='<Link target="...">display</Link>'>
Link to *target* (i.e. a URL), with *display* as the link's text.

For example, `<Link target="https://example.com">Example</Link>` becomes
<Link target="https://example.com">Example</Link>.

You can also use standard Markdown link syntax: `[Example](https://example.com)`.
</Define>

<Define tag="reference" sig='<Reference tag="...">display?</Reference>'>
Link to the target associated with *tag*. If children are provided they
become the link's content; otherwise the tag's configured display is
rendered.

You can also use the `[#tag]` shorthand in Markdown: `See [#some-tag]`.
</Define>

<Define tag="target" sig='<Target tag="..." title="display?"/>'>
Generate a target element that can be [#reference]d as *tag*. If *title*
is specified, references default to showing it as their link.
Otherwise, *tag* itself will be the default.

As an example, we'll create a target element in the following paragraph,
with the tag *some-tag* and title *I'm just some tag!*:

<Target tag="some-tag" title="I'm just some tag!"/> I'm a targetable
paragraph.

Then, we'll reference it with `<Reference tag="some-tag"/>`:

[#some-tag]

Clicking the above link should take you to the paragraph.
</Define>

## Flow Content

*Flow* content is anything that forms a *sentence*, i.e. a string of
words or inline elements.

<Define tag="code" sig='<Code>text</Code>'>
Present *text* in a monospace font upon rendering.

You can also use Markdown backticks for inline code: `` `code bits` ``.
For block-level code, use Markdown's triple-backtick fences with an
optional language tag:

````markdown
```go
fmt.Println("hello")
```
````
</Define>

<Define tag="italic" sig='<em>text</em>'>
Present *text* in *italics* upon rendering. Markdown: `*text*`.
</Define>

<Define tag="bold" sig='<strong>text</strong>'>
Present *text* in **bold** upon rendering. Markdown: `**text**`.
</Define>

<Define tag="larger" sig='<Larger>text</Larger>'>
Present *text* ~20% <Larger>larger</Larger> upon rendering.
</Define>

<Define tag="smaller" sig='<Smaller>text</Smaller>'>
Present *text* ~20% <Smaller>smaller</Smaller> upon rendering.
</Define>

<Define tag="strike" sig='<Strike>text</Strike>'>
Present *text* with a <Strike>strike through it</Strike> upon rendering.
</Define>

<Define tag="superscript" sig='<sup>text</sup>'>
Present *text* in <sup>superscript</sup> upon rendering.
</Define>

<Define tag="subscript" sig='<sub>text</sub>'>
Present *text* in <sub>subscript</sub> upon rendering.
</Define>

<Define tag="image" sig='<Image path="..."/>'>
Renders the image at *path* inline. Currently there is no "magic" that
will do anything with the file specified by *path* — if it's a local
path, make sure it's present in the output directory.

You can also use Markdown syntax: `![alt text](path)`.
</Define>

## Block Content

*Block* content is anything that forms a *paragraph*, i.e. a block of
text or any element that is standalone.

<Define tag="inset" sig='<Inset>content</Inset>'>
Render *content* indented a bit.

<Inset>
Like this!
</Inset>

You can also use Markdown blockquote syntax: `> content`.
</Define>

<Define tag="aside" sig='<Aside>content</Aside>'>
Render *content* in some way that conveys that it's a side-note.

<Aside>
Here I am!
</Aside>
</Define>

### Lists and Tables

Plain unordered and ordered lists, and tables, use standard Markdown
syntax — there's no JSX wrapper for them, because the Markdown forms
already cover the cases:

```
- one
- two
- three!

1. one
2. two
3. three!

| a | b | c |
| --- | --- | --- |
| 1 | 2 | 3 |
```

For *definition* lists, which Markdown doesn't have, use the lowercase
HTML `<dl>`/`<dt>`/`<dd>` tags directly:

<dl>
<dt>a</dt>
<dd>1</dd>
<dt>b</dt>
<dd>2</dd>
</dl>
