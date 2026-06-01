//go:build !cgo

package tree_sitter_booklit

import "unsafe"

// Language is unavailable when Booklit is built without cgo.
func Language() unsafe.Pointer {
	panic("tree-sitter booklit grammar requires cgo")
}
