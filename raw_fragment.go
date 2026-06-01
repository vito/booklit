package booklit

// RawFragment is a pre-rendered HTML chunk emitted verbatim by the
// renderer. Used for fragments that have no structural meaning to
// Booklit — the syntax highlighter's `<span class="...">` wrappers
// around tokens, for example, where the markup is the whole point and
// there is no inner Content tree to traverse.
//
// RawFragment is always flow. Block-level raw HTML is expressed by
// RawElement (which knows its tag and consults the htmltags block set)
// rather than by smuggling block markup through a "flow" carrier.
type RawFragment struct {
	// HTML is the literal markup, ready to write directly to the output.
	// Callers are responsible for any escaping.
	HTML string
}

// IsFlow always returns true.
func (con RawFragment) IsFlow() bool {
	return true
}

// String returns "" so that the markup bytes don't leak into
// plain-text consumers (search index, stringifyEverything). A
// RawFragment carries no user-visible text — it is pure decoration
// around content that lives in sibling nodes.
func (con RawFragment) String() string {
	return ""
}

// Visit calls VisitRawFragment.
func (con RawFragment) Visit(visitor Visitor) error {
	return visitor.VisitRawFragment(con)
}
