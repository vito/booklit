// Package templates is the mdx-template tier in Booklit's JSX dispatch.
// Each .md file in the configured directory becomes a JSX component:
// `<Foo prop="x">body</Foo>` looks up `Foo.md`, evaluates it with the
// props bound in Dang scope and `children` holding the JSX body, and
// emits whatever the template produces.
package templates

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/vito/booklit/ast"
	"github.com/vito/booklit/components"
	"github.com/vito/booklit/marklit"
)

// Registry loads and caches mdx template files from a list of search
// directories, falling back to the embedded stdlib (the components
// package).
//
// Lookups are by component name (PascalCase, matching the JSX tag):
// `<Foo/>` looks for `<dir>/Foo.md` in each directory in order, taking
// the first hit. If no on-disk directory has a match, the stdlib
// embed.FS is consulted last — so a project can override any stdlib
// component (Larger, Smaller, Strike, Inset, Aside, …) by dropping a
// same-named .md file into its components/ dir. Parsed ASTs are
// cached by mtime for on-disk templates (so edits in serve mode get
// picked up); stdlib templates are cached forever since their bytes
// are embedded into the binary at build time.
//
// A nil Registry still serves stdlib lookups (the stdlib FS is a
// package-level value, not held on the receiver). Pass nil only when
// the JSX evaluator should also skip the on-disk search.
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
// components directory found" case). Caller-supplied earlier dirs
// shadow later ones, which is how the test harness lets a per-test
// tempdir override shared fixtures.
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
// If no directory has a match, the embedded stdlib (components.Assets)
// is consulted last. Returns (nil, false, nil) when none of those
// sources has a matching template.
func (r *Registry) Load(name string) (ast.Node, bool, error) {
	if r == nil {
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
		return r.loadStdlib(name)
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

	r.cache[name] = cached{node: parseTemplate(source), modTime: info.ModTime()}
	return r.cache[name].node, true, nil
}

// loadStdlib resolves name from the embedded stdlib (components.Assets).
// Misses return (nil, false, nil) — they aren't errors, they're just
// "this component isn't in any source", which the JSX evaluator
// translates into an unknown-component build error one tier up.
//
// Stdlib bytes are baked into the binary, so the parsed AST is cached
// keyed by name alone; mtime isn't meaningful for embed.FS entries.
func (r *Registry) loadStdlib(name string) (ast.Node, bool, error) {
	source, err := fs.ReadFile(components.Assets, name+".md")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("read stdlib component %s: %w", name, err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if c, ok := r.cache[name]; ok {
		return c.node, true, nil
	}

	r.cache[name] = cached{node: parseTemplate(source)}
	return r.cache[name].node, true, nil
}

// parseTemplate strips the conventional trailing newline (so multiple
// invocations of the same template don't pile up blank lines) and
// parses via ParseInlineArg (not Parse) so a single-paragraph template
// doesn't pick up a stray `<p>` wrap. Block vs. inline is decided by
// what the template actually contains, not by goldmark's default
// paragraph wrapping. Interior whitespace is preserved verbatim
// because templates are HTML-significant.
func parseTemplate(source []byte) ast.Node {
	source = bytes.TrimRight(source, "\n")
	return marklit.ParseInlineArg(source)
}
