# booklit

[![Go Reference](https://pkg.go.dev/badge/github.com/vito/booklit.svg)](https://pkg.go.dev/github.com/vito/booklit)

Booklit is a tool for building static websites from semantic documents.

## documentation

[booklit.page](https://booklit.page)

## syntax

Booklit supports two document syntaxes, determined by file extension:

- **`.md`** — Markdown with `@invoke{arg}` function calls (recommended
  for new projects)
- **`.lit`** — the original Booklit syntax with `\invoke{arg}` function
  calls

Both syntaxes produce the same AST and work with the same plugins,
templates, and rendering pipeline. You can mix both in a single project
(e.g. `@include-section{child.md}` from a `.lit` parent).

### Markdown + @invoke (`.md`)

Standard Markdown plus Scribble-inspired function calls:

```markdown
# My Section

Hello, **world**! Here's a [link](https://example.com).

@include-section{child.md}

@section{
# Subsection

@italic{emphasized} text and `inline code`.
}
```

### Original syntax (`.lit`)

```
\title{My Section}

Hello, \bold{world}! Here's a \link{link}{https://example.com}.

\include-section{child.lit}

\section{
  \title{Subsection}

  \italic{emphasized} text and \code{inline code}.
}
```

## installation

Grab the latest [release](https://github.com/vito/booklit/releases), or
build from source:

```bash
go install github.com/vito/booklit/cmd/booklit@latest
```

## usage

```bash
booklit -i index.md -o ./out
```

## example

Clone this repo and build the Booklit website:

```bash
./scripts/build-docs
```

Then browse the generated docs from `./docs/index.html`.

Alternatively, run the docs in server mode:

```bash
./scripts/build-docs -s 3000
```

...and then browse to [localhost:3000](http://localhost:3000).
