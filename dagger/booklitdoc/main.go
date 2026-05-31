// Booklitdoc provides docs-site helpers for the Booklit documentation as
// Dagger functions. Each function returns Booklit content serialized with the
// contentjson wire format (as JSON), which Booklit decodes back into native
// content when the function is called from a {expr} interpolation.
package main

import (
	"bytes"
	"context"
	"regexp"

	"dagger/booklitdoc/internal/dagger"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/vito/booklit/contentjson/wire"
)

type Booklitdoc struct{}

// litStyle is the booklitdoc chroma palette, ported from the old in-process
// docs built-in.
var litStyle = chroma.MustNewStyle("booklitdoc", chroma.StyleEntries{
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

// funcRef matches a Booklit \function-name reference inside source.
var funcRef = regexp.MustCompile(`\\([a-z][a-z0-9-]*)`)

// LitSyntax highlights Booklit source code and returns it as serialized
// Booklit content. `\function` references in the code are linkified into
// Booklit references, which resolve against the section the content is decoded
// into. The result is the JSON wire format; Booklit turns it back into native
// content.
func (m *Booklitdoc) LitSyntax(
	ctx context.Context,
	// Booklit source code to highlight.
	code string,
	// Chroma lexer to use.
	// +optional
	// +default="lit"
	language string,
) (dagger.JSON, error) {
	node, err := litSyntaxNode(code, language)
	if err != nil {
		return "", err
	}
	data, err := wire.Marshal(node)
	if err != nil {
		return "", err
	}
	return dagger.JSON(data), nil
}

// litSyntaxNode builds the content tree for highlighted Booklit source,
// mirroring the old <LitSyntax> built-in: a lit-block wrapping a code-block
// whose content is the highlighted HTML with linkified \function references.
func litSyntaxNode(code, language string) (*wire.Node, error) {
	if language == "" {
		language = "lit"
	}
	highlighted, err := highlight(language, code)
	if err != nil {
		return nil, err
	}
	codeBlock := wire.StyledBlock("code-block", wire.Seq(linkify(highlighted)...))
	codeBlock.Partials = map[string]*wire.Node{"Language": wire.String(language)}
	return wire.StyledBlock("lit-block", codeBlock), nil
}

// highlight runs chroma over code and returns the styled HTML.
func highlight(language, code string) (string, error) {
	lexer := lexers.Get(language)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := html.New().Format(&buf, litStyle, iterator); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// linkify splits highlighted HTML into raw-html runs interspersed with Booklit
// references for each \function token, so the references resolve to their
// definitions on the Booklit side.
func linkify(s string) []*wire.Node {
	var out []*wire.Node
	last := 0
	for _, m := range funcRef.FindAllStringSubmatchIndex(s, -1) {
		if m[0] > last {
			out = append(out, wire.RawHTML(s[last:m[0]]))
		}
		fn := s[m[2]:m[3]]
		out = append(out, wire.String(`\`), wire.OptionalRef(fn, wire.String(fn)))
		last = m[1]
	}
	if last < len(s) {
		out = append(out, wire.RawHTML(s[last:]))
	}
	return out
}
