# Transition Plan: Backwards-Compatible Parser Swap

## Current State

The `marklit` branch replaces the PEG-based parser (`\function{arg}` syntax)
with a Markdown-based parser (`@function{arg}` syntax) built on goldmark. The
old PEG parser (`ast.ParseReader`) is still present in the tree but
`load/processor.go` unconditionally calls `marklit.Parse`. Any existing user
with `\function{arg}` `.lit` files would be silently broken — their `\` calls
become literal text.

## Proposed Approach: File-Extension-Based Dispatch

- **`.md` files → marklit parser** (new Markdown + `@invoke` syntax)
- **`.lit` files → old PEG parser** (existing `\invoke` syntax, preserved
  as-is)

In `load/processor.go`, dispatch based on `filepath.Ext(path)`:

```go
switch filepath.Ext(path) {
case ".md":
    node = marklit.Parse(source)
case ".lit":
    node = parseLegacy(path, file) // existing ast.ParseReader path
}
```

### Why This Works Well

- **Zero breakage** — existing `.lit` projects keep working unchanged.
- **Incremental adoption** — users migrate file-by-file by renaming
  `.lit` → `.md` and converting syntax.
- **Mixed projects** — `@include-section{child.md}` from a `.lit` parent
  (or vice versa) works naturally since both produce the same `ast.Node`
  types and everything downstream is unchanged.
- **Clean deprecation path** — announce `.lit` as deprecated, remove the
  PEG parser in a future major version.
- **Booklit's own docs** serve as the migration example (already
  converted).

## Steps to Ship

1. Restore the PEG parser dispatch for `.lit` files.
2. Keep marklit as the parser for `.md` files.
3. Rename `docs/lit/*.lit` → `docs/lit/*.md` (they're already in the new
   syntax).
4. Keep the old PEG code + pigeon dependency (remove in a later release).
5. Update README to describe both syntaxes and the migration path.
6. Optionally: add a `booklit migrate` subcommand that converts
   `.lit` → `.md`.

## What NOT to Do

- Don't add a CLI flag for parser selection — file extension is simpler
  and self-documenting.
- Don't try to auto-detect syntax — `@` vs `\` heuristics would be
  fragile.
- Don't remove the PEG parser yet — that's a separate major version
  decision.
