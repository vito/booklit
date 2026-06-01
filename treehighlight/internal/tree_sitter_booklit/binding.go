//go:build cgo

package tree_sitter_booklit

// #cgo CFLAGS: -std=c11 -fPIC
// #include "src/parser.c"
import "C"

import "unsafe"

// Language returns the tree-sitter Language for the Booklit/Markdown-ish grammar.
func Language() unsafe.Pointer {
	return unsafe.Pointer(C.tree_sitter_booklit())
}
