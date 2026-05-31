# Plugins {#plugins}

In the new Booklit model there is no separate "plugin" system: components
are resolved by a tiered dispatcher. When the evaluator hits a JSX
invocation like `<Foo bar="x">body</Foo>`, it looks for:

<OrderedList>
<Item>A *built-in* with name `Foo` registered in Go.</Item>
<Item>An HTML *template* named `Foo.tmpl` in the templates directory.</Item>
<Item>(later) A *Dang* function in scope.</Item>
<Item>(later) A *Dagger* function configured in `booklit.toml`.</Item>
</OrderedList>

The simplest "plugin" is a template — no Go code required. This page
walks through both approaches.

<TableOfContents/>

## Template-only Components {#template-components}

If a JSX invocation has no matching built-in, Booklit wraps it as a
<Godoc ref="booklit.Styled"/> with `Style` equal to the component name
and props passed through as Partials. The HTML renderer then looks up
`<Name>.tmpl` in the templates directory.

Say you want a `<Card>` component. Drop the following into `html/Card.tmpl`:

```go-html-template
<div class="card">
  <h3>{{.Partial "title" | render}}</h3>
  <div class="body">{{.Content | render}}</div>
</div>
```

Then use it in any document:

```markdown
<Card title="Greetings">
Welcome to the test.
</Card>
```

Build with `booklit --html-templates html ...` and the component renders.
No Go, no recompile. Props are passed through with their authored
(camelCase) names — `<Card title="..."/>` → `{{.Partial "title"}}`.

## Go Built-ins {#go-builtins}

When a component needs to do something a template can't express (touch
the section tree, evaluate AST sub-trees, error out), write it as a Go
built-in.

Built-ins are registered via the <Godoc ref="builtins.Register"/>
function. To add one to your own docs site, create a Go module that
imports `github.com/vito/booklit/builtins` and registers in `init()`:

```go
package mycomps

import (
  "github.com/vito/booklit"
  "github.com/vito/booklit/ast"
  "github.com/vito/booklit/builtins"
)

func init() {
  builtins.Register("HelloWorld", helloWorld)
}

func helloWorld(ctx *builtins.Context, _ map[string]ast.Node, _ []ast.Node) (booklit.Content, error) {
  return booklit.String("Hello, world!"), nil
}
```

Then create a binary that imports both your package and `booklitcmd`:

```go
package main

import (
  "github.com/vito/booklit/booklitcmd"
  _ "example.com/yourmodule/mycomps"
)

func main() {
  booklitcmd.Main()
}
```

Use that binary instead of the stock `booklit` CLI to build your docs.
This is exactly how the documentation site you're reading is built —
see `cmd/booklit-docs` and `docs/booklitdoc/` in the Booklit repo.

### Built-in Signature

A built-in is a function with this shape:

```go
type Func func(
  ctx *builtins.Context,
  props map[string]ast.Node,
  children []ast.Node,
) (booklit.Content, error)
```

<Definitions>
<Definition term="ctx.Section">
The current section. Built-ins that mutate the section tree (like
[#section] and [#title]) use this directly.
</Definition>
<Definition term="ctx.Evaluate(node)">
Evaluates an AST node and returns its booklit.Content result. Use this
to evaluate individual props or children on demand.
</Definition>
<Definition term="props">
The component's props as raw AST. Each value is either an
<Godoc ref="ast.String"/> (string-literal `attr="x"`) or an
<Godoc ref="ast.JSXExpression"/> (`attr={expr}`). Call
`ctx.Evaluate` on a prop to evaluate it.
</Definition>
<Definition term="children">
The component's children as raw AST nodes. For single-line invocations
these are inline content; for multi-line invocations they are block
content (paragraphs).
</Definition>
</Definitions>

### Return Values

A built-in returns a `(booklit.Content, error)` pair. Returning `(nil,
nil)` is fine — useful for pure-side-effect components like
[#split-sections] which just configure the section.

### A Full Example

Putting the pieces together, here's a real built-in that registers
multiple tags as targets and renders a bold name followed by a
description:

```go
func describeFruit(ctx *builtins.Context, props map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
  name, err := requireStringProp(ctx, props, "name")
  if err != nil {
    return nil, err
  }

  body, err := ctx.Evaluate(ast.Sequence(children))
  if err != nil {
    return nil, err
  }

  var content booklit.Sequence
  for _, tag := range strings.Fields(name) {
    content = append(content, booklit.Target{
      TagName:  tag,
      Location: ctx.Section.InvokeLocation,
      Title:    booklit.String(name),
      Content:  body,
    })
  }

  content = append(content,
    booklit.Paragraph{booklit.Styled{Style: booklit.StyleBold, Content: booklit.String(name)}},
    body,
  )
  return content, nil
}
```

Called as:

```markdown
<DescribeFruit name="banana">
A banana is a yellow fruit that only really tastes
good in its original form. Banana flavored
anything is a pit of despair.
</DescribeFruit>
```

The `<Target>` elements register `banana` as a reference target, so
[`<Reference tag="banana"/>`] anywhere in the section will link to it.
