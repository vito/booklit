# Transition Plan: Backwards-Compatible Parser Swap

## Current State

The `marklit` branch adds a Markdown-based parser (`@function{arg}` syntax)
built on goldmark alongside the existing PEG parser (`\function{arg}` syntax).
Both parsers produce the same `ast.Node` types, so everything downstream
(evaluation, plugins, rendering) works unchanged.

## Approach: File-Extension-Based Dispatch

- **`.md` files → marklit parser** (Markdown + `@invoke` syntax)
- **`.lit` files → PEG parser** (original `\invoke` syntax)

Dispatch happens in `load/processor.go` based on `filepath.Ext(path)`.

### Why This Works Well

- **Zero breakage** — existing `.lit` projects keep working unchanged.
- **Incremental adoption** — users migrate file-by-file by renaming
  `.lit` → `.md` and converting syntax.
- **Mixed projects** — `@include-section{child.md}` from a `.lit` parent
  (or vice versa) works naturally since both produce the same `ast.Node`
  types and everything downstream is unchanged.
- **Clean deprecation path** — the PEG parser can be removed in a future
  major version if desired.
- **Booklit's own docs** serve as the migration example (already
  converted).

## Completed Steps

- [x] File-extension dispatch in `load/processor.go`.
- [x] Integration tests use `.md` extension.
- [x] `docs/lit/*.lit` renamed to `docs/lit/*.md`.
- [x] `@include-section` references and build script updated.
- [x] README updated with both syntaxes.
- [x] PEG parser and pigeon dependency retained.

## Future Considerations

- Add a `booklit migrate` subcommand to convert `.lit` → `.md`.
- Consider deprecating `.lit` in a future major version.
