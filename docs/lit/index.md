# Booklit {#index}

<Styled name="splash"/>

<Larger>
Booklit is a tool for building static websites from semantic documents.
</Larger>

Booklit enforces a tidy separation between **content**, **logic**, and
**presentation** that makes it easy to write and refactor websites for
anything from technical documentation to slide decks to blogs.

For live examples, check out [Bass](https://bass-lang.org),
[Concourse CI](https://concourse-ci.org), and the site you're currently
viewing ([src](https://github.com/vito/booklit/tree/master/docs/lit)).

You're probably wondering "why does the world need another static site
generator?" The truth is I built this for myself; I had a lot of technical
content to maintain, and I didn't like the state of the art. I wanted something
more like [Scribble](https://docs.racket-lang.org/scribble/index.html) so
I could write code to minimize toil.

Booklit has been serving me well across multiple projects for years with little
modification needed, so I think it's good enough to share.


<Columns>
<ColumnHeader>content in `*.md`</ColumnHeader>

Booklit documents are Markdown with embedded `<Component>` invocations that
generate content, ultimately resulting in a tree of sections.

Sections are easy to move around, allowing you to continuously refactor and
restructure your content without having to tweak header sizes and update
internal links.

<Column>

```markdown
# Hello {#index}

Hello, world! I'm a Booklit document.

Check out my favorite [#quotes]!

<IncludeSection path="quotes.md"/>
```

```markdown
# Quotes

<Quote source="Travis Scott">
  It's lit!
</Quote>
```

</Column>
</Columns>

<Columns>
<ColumnHeader>logic in `*.go` (or templates)</ColumnHeader>

Components dispatch via a tiered resolver: built-in &rarr; HTML template
&rarr; and eventually Dang or Dagger. Most "plugins" are just an HTML
template; you only write Go for primitives that need to manipulate the
section tree.

<Column>

```go-html-template
<blockquote class="quote">
  {{.Content | render}}

  <footer>
    {{.Partial "source" | render}}
  </footer>
</blockquote>
```

</Column>
</Columns>

<Columns>
<ColumnHeader>presentation in `*.tmpl`</ColumnHeader>

Booklit separates presentation into a final rendering phase which determines
the output format.

The [#html-renderer] is powered by Go's standard
[`html/template` package](https://golang.org/pkg/html/template/).
More renderers may be implemented in the future.

All [base templates](#base-templates) can be overridden, sections
can be individually [#styled], and components can return
<Godoc ref="booklit.Styled"/> content, giving the author full control over what
comes out.

<Column>

```go-html-template
<!DOCTYPE html>
<html>
  <head>
    <title>{{.Title.String}}</title>
  </head>
  <body>
    {{. | render}}
  </body>
</html>
```

</Column>
</Columns>

<Columns>
<ColumnHeader>build with `booklit`</ColumnHeader>

The `booklit` CLI is a single command which loads Booklit documents
and renders them.

When an error occurs, Booklit will show the location of the error and try to
suggest a fix.

<Column>

```
$ booklit -i ./index.md -o ./public/
INFO[0000] rendering
```

```
$ booklit -i ./to-err-is-human.md
to-err-is-human.md:5: unknown tag 'helo'

   5| Say [#helo]!
          ^^^^^^^
These tags seem similar:

- hello

Did you mean one of these?
```

</Column>
</Columns>

<Columns>
<ColumnHeader>serve with `booklit -s $PORT`</ColumnHeader>

In server mode, Booklit renders content on each request.
The feedback loop is *wicked fast*.

<Column>

```
$ booklit -i ./index.md -s 3000
INFO[0000] listening
```

<OutputFrame url="outputs/index.html"/>

</Column>
</Columns>

This website is [written with
Booklit](https://github.com/vito/booklit/tree/master/docs/lit). Want to write
your own? Let's [get started](#getting-started)!

<SplitSections/>

<IncludeSection path="getting-started.md"/>
<IncludeSection path="baselit.md"/>
<IncludeSection path="html-renderer.md"/>
<IncludeSection path="plugins.md"/>
<IncludeSection path="syntax.md"/>
<IncludeSection path="thanks.md"/>
