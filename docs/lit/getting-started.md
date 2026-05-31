# Getting Started

The best way to get started with Booklit is to install the CLI:

```sh
# add booklit to your toolchain
go install github.com/vito/booklit/cmd/booklit@latest

# add GOPATH/bin to your $PATH
export PATH=$(go env GOPATH)/bin:$PATH
```

It's also possible to download the `booklit` executable from the latest
[GitHub release](https://github.com/vito/booklit/releases/latest), but
tracking it via `go install` makes it easier to follow updates.

<TableOfContents/>

## Hello, world!

First, create a file called `hello.md` with the following content:

```markdown
# Hello, world! {#hello}

I'm a Booklit document!
```

This file can exist anywhere, but one common convention is to place
`.md` documents under `lit/`, HTML templates under `html/`,
and any custom Go code under `go/`.

Run the following to build and render the file to `./docs/hello.html`:

```bash
$ booklit -i hello.md -o docs
```

Each of the changes in the following sections will require re-building, which
can be done by running the above command again. Alternatively, you can run
`booklit` with the `-s` flag to start a HTTP server:

```
$ booklit -i hello.md -s 8000
INFO[0000] listening              port=8000
```

Once Booklit says 'listening', browse to
[http://localhost:8000/hello.html](http://localhost:8000/hello.html).
When you change anything, just refresh and your content will be rebuilt and
re-rendered.

## Organizing with Sections

Next, let's try adding a section within our document. Headings (`#`, `##`,
`###`) become sections automatically — top-level becomes the page title,
each deeper level nests a sub-section:

```markdown
# Hello, world! {#hello}

I'm a Booklit document!

## Hi there!

I'm so organized!
```

After building, you should see something like this:

<Inset>
<Larger><Larger><Larger>Hello, world!</Larger></Larger></Larger>

I'm a Booklit document!

<Larger><Larger>1 Hi there!</Larger></Larger>

I'm so organized!
</Inset>

That number "1" might look a bit weird at the moment, but it's the section
number, and it'll be something like "3.2" for a nested section. You can always
remove it by specifying your own template (more on that later), but for now
let's leave it there.

For non-heading sections, or sections constructed by a component, use
`<Section>` explicitly:

```markdown
<Section>
  ## Custom Sub-section
  Body.
</Section>
```

## Splitting Sections

To render each sub-section on its own page, simply call
[#split-sections] somewhere in the section.

```markdown
# Hello, world! {#hello}

<SplitSections/>

I'm a Booklit document!

## Hi there!

I'm so organized!
```

So far we've just made the section disappear, which isn't very helpful. Let's
at least make it so we can browse to it! This can be done with
[#table-of-contents]:

```markdown
# Hello, world! {#hello}

<SplitSections/>

I'm a Booklit document!

<TableOfContents/>

## Hi there!

I'm so organized!
```

Note that when viewing the sub-section, its header is now a `<h1>`
rather than the `<h2>` it was before, since it stands on its own page.

## References & Tagging

Having a [#table-of-contents] is great and all, but more often
you'll want to reference sections from each other directly and in context.
Use the `[#tag]` shorthand:

```markdown
# Hello, world! {#hello}

<SplitSections/>

I'm a Booklit document! To read further, see [#hi-there].

## Hi there!

I'm so organized!
```

The default tag is a slugified form of the heading; you can also give an
explicit tag via `{#tag}` after the heading text. To override the link's
display name, write `[Display Name](#tag)`:

```markdown
I'm a Booklit document! Consult [this section](#hi-there) for more.
```

## Next Steps

What we've gone over should carry you pretty far. But you'll likely want
to know a lot more.

- To change how your generated content looks, check out the
  [HTML renderer](#html-renderer).
- To learn the components that come with Booklit, check out [#baselit].
- To extend your documents with your own components, check out [#plugins].
