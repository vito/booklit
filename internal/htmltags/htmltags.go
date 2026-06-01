// Package htmltags lists the HTML element names treated as block-level
// for Booklit's layout decisions: paragraph wrapping and the IsFlow
// check on RawElement.
//
// The set is the union of CommonMark's HTML block-tag list (types 1
// and 6 from the spec) and HTML5's block-level elements. CommonMark
// uses these names to decide when an HTML element interrupts a
// paragraph; Booklit uses the same membership to decide whether a
// lowercase JSX element should be wrapped in a paragraph or rendered
// standalone.
//
// Goldmark maintains an equivalent table internally
// (parser/html_block.go) but does not export it, so it is restated
// here.
package htmltags

// Block holds the HTML tag names treated as block-level. Keys are
// lowercase; callers should lowercase the tag before lookup since
// HTML element names are case-insensitive.
var Block = map[string]bool{
	"address":    true,
	"article":    true,
	"aside":      true,
	"base":       true,
	"basefont":   true,
	"blockquote": true,
	"body":       true,
	"caption":    true,
	"center":     true,
	"col":        true,
	"colgroup":   true,
	"dd":         true,
	"details":    true,
	"dialog":     true,
	"dir":        true,
	"div":        true,
	"dl":         true,
	"dt":         true,
	"fieldset":   true,
	"figcaption": true,
	"figure":     true,
	"footer":     true,
	"form":       true,
	"frame":      true,
	"frameset":   true,
	"h1":         true,
	"h2":         true,
	"h3":         true,
	"h4":         true,
	"h5":         true,
	"h6":         true,
	"head":       true,
	"header":     true,
	"hr":         true,
	"html":       true,
	"iframe":     true,
	"legend":     true,
	"li":         true,
	"link":       true,
	"main":       true,
	"menu":       true,
	"menuitem":   true,
	"meta":       true,
	"nav":        true,
	"noframes":   true,
	"ol":         true,
	"optgroup":   true,
	"option":     true,
	"p":          true,
	"param":      true,
	"pre":        true,
	"script":     true,
	"search":     true,
	"section":    true,
	"style":      true,
	"summary":    true,
	"table":      true,
	"tbody":      true,
	"td":         true,
	"textarea":   true,
	"tfoot":      true,
	"th":         true,
	"thead":      true,
	"title":      true,
	"tr":         true,
	"track":      true,
	"ul":         true,
}
