\use-plugin{booklitdoc}
\use-plugin{chroma}

# Plugins {#plugins}

Plugins provide the functionality behind function calls like
\syntax{lit}{\\foo{bar}}.

Out of the box, Booklit comes with a plugin called
[`baselit`](#baselit) which provides basic functions like
[#title], [#section], [#italic], and
[#bold].

More functions can be added by writing plugins and using them in your
documents.

If you've skipped ahead, you may want to check out [#getting-started]
to see how to set up your Go module.

\table-of-contents

## Using Plugins {#using-plugins}

To use a plugin, pass its Go package's import path as `--plugin` to the
`booklit` command when building your docs.

For example, Booklit comes with a `chroma` plugin for syntax
highlighting. To use it, run:

\syntax{bash}{{{
booklit -i index.lit -o out \
  --plugin github.com/vito/booklit/chroma/plugin
}}}

The `--plugin` flag must be passed every time you build your docs,
so you may want to put it in a script:

```bash
#!/bin/bash

booklit -i lit/index.lit -o public \
  --plugin github.com/vito/booklit/chroma/plugin \
  "$@" # forward args from script to booklit
```


Booklit imports all specified plugins at build time, automatically adding
them to `go.mod`. When imported, plugins register themselves under a
certain name - typically guessable from the import path.

To use the plugin in your documents, call [#use-plugin] with its
registered name:

\lit-syntax{{{
\title{My Section}

\use-plugin{chroma}

\syntax{ruby}{{{
  def fib(n)
    fib(n - 2) + fib(n - 1)
  end
}}}
}}}

The `--plugin` flag can be specified multiple times, and
[#use-plugin] can be invoked multiple times.

Note: [inline sections](#section) inherit plugins from their parent
sections, but [included sections](#include-section) do not.

## Writing Plugins {#using-plugins}

Plugins are just Go packages that register a *plugin factory* with
Booklit when they're imported with the `--plugin` flag.

It's possible to use Booklit without writing any plugins of your own, but
being able to write a plugin help you get the most out of Booklit.

To create a new plugin, create a directory within your Go module (where
`go.mod` lives) - let's call it `example` for this example:

\syntax{bash}{{{
mkdir example
}}}

Then, we'll create the initial skeleton for our plugin at
`example/plugin.go`:

```go
package example

import (
  "github.com/vito/booklit"
)

func init() {
  booklit.RegisterPlugin("example", NewPlugin)
}

func NewPlugin(sec *booklit.Section) booklit.Plugin {
  return Plugin{
    section: sec,
  }
}

type Plugin struct {
  section *booklit.Section
}
```

This registers a plugin that does nothing. Let's define some document
functions!

Functions work by simply defining methods on the plugin struct. Let's define
a basic one with no arguments:

```go
func (plugin Plugin) HelloWorld() booklit.Content {
  return booklit.String("Hello, world!")
}
```

Now let's create a Booklit document that uses it as `hello-plugins.lit`:

\lit-syntax{{{
\title{Hello Plugins}

\use-plugin{example}

Zero args: \hello-world
}}}

To build this document, pass the package import path (including your module
name) as the `--plugin` flag. For example, if your `go.mod` says
`module foo`, the flag would be:

```bash
booklit -i hello-plugins.lit -o out \
    --plugin foo/example
```

This should result in a page showing:

\inset{
  Zero args: Hello, world!
}

### Argument Types

Functions can be invoked with any number of arguments, like so:

\lit-syntax{{{
\hello-world{arg1}{arg2}
}}}

See [#function-syntax] for more information.

Each argument to the function corresponds to an argument for the plugin's
method, which may be variadic.

The plugin's arguments must each be one of the following types:

\definitions{
  \definition{\godoc{booklit.Content}}{
    The evaluated content. This can be just about anything from a word to a
    sentence to a series of paragraphs, depending on how the function is
    invoked. It is typically used unmodified.
  }
}{
  \definition{`string`}{
    The evaluated content, converted into a string. This is useful when the
    content is expected to be something simple, like a word or line of
    text. The [#title] function, for example, uses this type for
    its variadic *tags* argument.
  }
}{
  \definition{\godoc{booklit/ast.Node}}{
    The unevaluated syntax tree for the content. This is useful when doing
    meta-level things like [#section] which need to control the
    evaluation context of the content.
  }
}

### Return Values

Plugin methods can then return one of the following:

\list{
  nothing
}{
  `error`
}{
  \godoc{booklit.Content}
}{
  `(`\godoc{booklit.Content}`, error)`
}

If a method returns a non-nil `error` value, it will bubble up and
the building will fail.

### A Full Example

Putting the pieces together, let's extend our `pluglit` plugin from
earlier write a real function that does something useful:

```go
func (plugin Plugin) DescribeFruit(
  name string,
  definition booklit.Content,
  tags ...string,
) (booklit.Content, error) {
  if name == "" {
    return nil, errors.New("name cannot be blank")
  }

  content := booklit.Sequence{}
  if len(tags) == 0 {
    tags = []string{name}
  }

  for _, tag := range tags {
    content = append(content, booklit.Target{
      TagName: tag,
      Display: booklit.String(name),
    })
  }

  content = append(content, booklit.Paragraph{
    booklit.Styled{
      Style: booklit.StyleBold,
      Content: booklit.String(name),
    },
  })

  content = append(content, definition)

  return content, nil
}
```

There are many things to note here:

\list{
  there are two required arguments; *name* is a `string` and
  *value* is a \godoc{booklit.Content}
}{
  there's a variadic argument, *tags*, which is of type
  `[]string`
}{
  this function generates content, and can raise an error when building
}{
  the \godoc{booklit.Target} elements will result in tags being registered
  in the section the function is called from
}{
  the function name, `describe-fruit`, corresponds to the method name
  `DescribeFruit`
}

This function would be called like so:

\lit-syntax{{{
\describe-fruit{banana}{
  A banana is a yellow fruit that only really tastes
  good in its original form. Banana flavored
  anything is a pit of dispair.
}{banana-opinion}
}}}

...and will result in something like the following:

\inset{
  \describe-fruit{banana}{
    A banana is a yellow fruit that only really tastes
    good in its original form. Banana flavored
    anything is a pit of dispair.
  }{banana-opinion}
}

...which can be referenced as \syntax{lit}{\\reference{banana-opinion}}, which
results in a link like this: [#banana-opinion].
