// Package booklitdoc provides docs-site customization for Booklit's own
// documentation. It is imported by cmd/booklit-docs for its side effect of
// overriding the chroma syntax-highlighting palette; the main cmd/booklit
// binary does not include it.
//
// The doc-specific components are no longer Go built-ins: <Columns> / <Column>
// / <ColumnHeader> / <Define> / <Godoc> / <OutputFrame> / <TemplateLink> are
// mdx templates in docs/html/, and <LitSyntax> dispatches to the booklitdoc
// Dagger module (dagger/booklitdoc), which returns highlighted content as the
// contentjson wire format.
package booklitdoc

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/styles"
)

func init() {
	// The booklitdoc palette colors every fenced code block on the docs
	// site, which baselit still highlights in-process via styles.Fallback.
	// The Dagger module carries its own copy for <LitSyntax>; keep the two
	// in sync.
	styles.Fallback = chroma.MustNewStyle("booklitdoc", chroma.StyleEntries{
		chroma.Comment:               "#c29d7c italic",
		chroma.CommentPreproc:        "noitalic",
		chroma.Keyword:               "#ed6c30 bold",
		chroma.KeywordPseudo:         "nobold",
		chroma.KeywordType:           "nobold",
		chroma.OperatorWord:          "#fcc21b bold",
		chroma.NameClass:             "#fcc21b bold",
		chroma.NameNamespace:         "#fcc21b bold",
		chroma.NameException:         "#fcc21b bold",
		chroma.NameEntity:            "#fcc21b bold",
		chroma.NameTag:               "#fcc21b bold",
		chroma.LiteralString:         "#fcc21b",
		chroma.LiteralStringInterpol: "bold",
		chroma.GenericHeading:        "bold",
		chroma.GenericSubheading:     "bold",
		chroma.GenericEmph:           "italic",
		chroma.GenericStrong:         "bold",
		chroma.GenericPrompt:         "bold",
		chroma.Error:                 "border:#FF0000",
	})
}
