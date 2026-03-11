\use-plugin{booklitdoc}
\use-plugin{chroma}

# Booklit {#index}

\styled{splash}

\larger{
  Booklit is a tool for building static websites from semantic documents.
}

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


\columns{
  \column-header{content in `*.lit`}

  Booklit documents are text files which evaluate \syntax{lit}{\\functions} to
  generate content, ultimately resulting in a tree of sections.

  Sections are easy to move around, allowing you to continuously refactor and
  restructure your content without having to tweak header sizes and update
  internal links.
}{
  \lit-syntax{{{
  \title{Hello}{index}

  Hello, world! I'm a Booklit document.

  Check out my favorite \reference{quotes}!

  \include-section{quotes.lit}
  }}}

  \lit-syntax{{{
  \title{Quotes}
  \use-plugin{example}

  \quote{
    It's lit!
  }{Travis Scott}
  }}}
}

\columns{
  \column-header{logic in `*.go`}

  Sections use plugins to invoke \syntax{lit}{\\functions} written in
  [Go](https://golang.org). Go is a simple and fast language with
  [plenty of packages](https://pkg.go.dev/) around if you need them.

  Plugins allow your content to be semantic - saying what it means, decoupled
  from how it should be computed or displayed.
}{
  \syntax{go}{{{
  func (Example) Quote(
    quote, source booklit.Content,
  ) booklit.Content {
    return booklit.Styled{
      Style: "quote",
      Content: quote,
      Partials: booklit.Partials{
        "Source": source,
      },
    }
  }
  }}}
}

\columns{
  \column-header{presentation in `*.tmpl`}

  Booklit separates presentation into a final rendering phase which determines
  the output format.

  The [#html-renderer] is powered by Go's standard
  [`html/template` package](https://golang.org/pkg/html/template/).
  More renderers may be implemented in the future.

  All [base templates](#base-templates) can be overridden, sections
  can be individually [#styled], and plugins can return
  \godoc{booklit.Styled} content, giving the author full control over what
  comes out.
}{
  \syntax{go-html-template}{{{
  <!DOCTYPE html>
  <html>
    <head>
      <title>{{.Title.String}}</title>
    </head>
    <body>
      {{. | render}}
    </body>
  </html>
  }}}

  \syntax{go-html-template}{{{
  <blockquote class="quote">
    {{.Content | render}}

    <footer>
      {{.Partial "Source" | render}}
    </footer>
  </blockquote>
  }}}
}

\columns{
  \column-header{build with `booklit`}

  The `booklit` CLI is a single command which loads Booklit documents
  and renders them.

  When an error occurs, Booklit will show the location of the error and try to
  suggest a fix.
}{
  \code{{
  $ booklit -i ./index.lit -o ./public/
  \syntax-hl{INFO}[0000] rendering
  }}

  \code{{
  $ booklit -i ./to-err-is-human.lit
  to-err-is-human.lit:5: unknown tag 'helo'

     5| Say \\reference{helo}!
            \syntax-hl{^^^^^^^^^^}
  These tags seem similar:

  - hello

  Did you mean one of these?
  }}
}

\columns{
  \column-header{serve with `booklit -s $PORT`}

  In server mode, Booklit renders content on each request with only plugin
  changes requiring a server restart.

  The feedback loop is *wicked fast*.
}{
  \code{{
  $ booklit -i ./index.lit -s 3000
  \syntax-hl{INFO}[0000] listening
  }}

  \output-frame{outputs/index.html}
}

This website is [written with
Booklit](https://github.com/vito/booklit/tree/master/docs/lit). Want to write
your own? Let's [get started](#getting-started)!

\split-sections

\include-section{getting-started.md}
\include-section{baselit.md}
\include-section{html-renderer.md}
\include-section{plugins.md}
\include-section{syntax.md}
\include-section{thanks.md}
