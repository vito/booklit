// Package templates is the mdx-template tier in Booklit's JSX dispatch.
// Each .md file in the configured directory becomes a JSX component:
// `<Foo prop="x">body</Foo>` looks up `Foo.md`, evaluates it with the
// props bound in Dang scope and `children` holding the JSX body, and
// emits whatever the template produces.
package templates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/vito/booklit/ast"
)

// Registry loads and caches mdx template files from a directory.
//
// Lookups are by component name (PascalCase, matching the JSX tag):
// `<Foo/>` looks for `<dir>/Foo.md`. Parsed ASTs are cached by mtime so
// edits in serve mode are picked up automatically.
//
// A nil Registry is valid: it returns (nil, false, nil) for every lookup,
// matching the "no templates configured" case.
type Registry struct {
	dir string

	mu    sync.Mutex
	cache map[string]cached
}

type cached struct {
	node    ast.Node
	modTime time.Time
}

// New returns a Registry rooted at dir. An empty dir disables the
// registry entirely (Load always misses), which is what callers want
// when no --html-templates flag was passed.
func New(dir string) *Registry {
	return &Registry{dir: dir, cache: map[string]cached{}}
}

// Load returns the parsed template for name, or (nil, false, nil) if no
// `<dir>/<name>.md` file exists. Any parse error from a present file is
// surfaced as the third return.
func (r *Registry) Load(name string) (ast.Node, bool, error) {
	if r == nil || r.dir == "" {
		return nil, false, nil
	}

	path := filepath.Join(r.dir, name+".md")
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("stat %s: %w", path, err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if c, ok := r.cache[name]; ok && !info.ModTime().After(c.modTime) {
		return c.node, true, nil
	}

	source, err := os.ReadFile(path)
	if err != nil {
		return nil, false, fmt.Errorf("read %s: %w", path, err)
	}

	// Trim a single trailing newline (the conventional file terminator)
	// so multiple invocations of the same template don't pile up blank
	// lines in the rendered output. Interior whitespace is preserved
	// verbatim because templates are HTML-significant.
	source = bytes.TrimRight(source, "\n")

	node, err := Parse(source)
	if err != nil {
		return nil, false, fmt.Errorf("parse %s: %w", path, err)
	}

	r.cache[name] = cached{node: node, modTime: info.ModTime()}
	return node, true, nil
}
