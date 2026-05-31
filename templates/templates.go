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

// Registry loads and caches mdx template files from a list of search
// directories.
//
// Lookups are by component name (PascalCase, matching the JSX tag):
// `<Foo/>` looks for `<dir>/Foo.md` in each directory in order, taking
// the first hit. Parsed ASTs are cached by mtime so edits in serve mode
// are picked up automatically.
//
// A nil Registry is valid: it returns (nil, false, nil) for every lookup,
// matching the "no templates configured" case.
type Registry struct {
	dirs []string

	mu    sync.Mutex
	cache map[string]cached
}

type cached struct {
	node    ast.Node
	modTime time.Time
}

// New returns a Registry that searches dirs in order. Empty strings
// are skipped (so `New("")` is a no-op registry, matching the "no
// --html-templates flag" case). Caller-supplied earlier dirs shadow
// later ones, which is how the test harness lets a per-test tempdir
// override shared fixtures.
func New(dirs ...string) *Registry {
	filtered := make([]string, 0, len(dirs))
	for _, d := range dirs {
		if d != "" {
			filtered = append(filtered, d)
		}
	}
	return &Registry{dirs: filtered, cache: map[string]cached{}}
}

// Load returns the parsed template for name. Searches each configured
// directory in order; the first `<dir>/<name>.md` that exists wins.
// Returns (nil, false, nil) when no directory has a matching template.
func (r *Registry) Load(name string) (ast.Node, bool, error) {
	if r == nil || len(r.dirs) == 0 {
		return nil, false, nil
	}

	var path string
	var info os.FileInfo
	for _, dir := range r.dirs {
		candidate := filepath.Join(dir, name+".md")
		st, err := os.Stat(candidate)
		if err == nil {
			path = candidate
			info = st
			break
		}
		if !os.IsNotExist(err) {
			return nil, false, fmt.Errorf("stat %s: %w", candidate, err)
		}
	}
	if path == "" {
		return nil, false, nil
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
