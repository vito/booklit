// Package components ships the stdlib MarkDangJSX component library
// embedded into the binary. Each *.md file under this directory is a
// component the JSX evaluator dispatches to as a fallback after the
// built-ins and the project-local components/ directory.
//
// The stdlib covers the small set of styling components Booklit
// previously implemented as `Styled{Style: "..."}` wrappers with
// matching Go html/template files in render/html/: Larger, Smaller,
// Strike, Inset, Aside. Authoring them as MarkDangJSX components
// instead of Go-side primitives means they obey the same parsing and
// dispatch rules as user components — no template-name string
// indirection, no special-case Styled handling at render time, and
// any project's components/ dir can replace one by simply dropping a
// same-named .md file in.
package components

import "embed"

// Assets holds the stdlib components as files. Filenames match the
// JSX tag PascalCase (`Larger` -> `Larger.md`). The templates.Registry
// consults this FS after the project-local dirs miss, so project
// overrides take precedence by name.
//
//go:embed *.md
var Assets embed.FS
