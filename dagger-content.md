# Dagger content: state

This was the working doc for experimenting with Booklit content produced
by a Dagger module. The experiment landed the generic pieces, but the
original docs use case (`<LitSyntax>`) has since been removed.

## Current status

Booklit still has the generic bridge for Dagger-produced content:

- `contentjson/wire` is a dependency-free JSON representation of the
  serializable subset of `booklit.Content`.
- `contentjson` converts between that wire format and native
  `booklit.Content`.
- `dangeval.Evaluator.ContentFromValue` decodes Dagger `JSON` or
  `JSONValue` scalars as content and rehydrates references/targets
  against the current section.
- `dangeval.New` serves the discovered Dagger module into the live
  session so dependency calls can work at runtime, not just at type-check
  time.

There is no in-tree Dagger producer at the moment.

## What replaced `<LitSyntax>`

`<LitSyntax>` and `dagger/booklitdoc` were deleted. Fenced code blocks
now use `baselit.Syntax`, which renders with `treehighlight` directly and
emits real `booklit.Reference` nodes for linkable tree-sitter captures.
That means examples like this are just ordinary Markdown:

````markdown
```markdown
<IncludeSection path="quotes.md"/>
```
````

The `IncludeSection` tag name is captured by tree-sitter, converted to
`include-section`, and rendered as an optional Booklit reference if a tag
with that name exists.

## Notes

`treehighlight` uses cgo-backed tree-sitter bindings when cgo is
available. If Booklit is built with `CGO_ENABLED=0`, it compiles to an
escaped plain-code fallback.
