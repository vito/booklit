// Package booklitdoc provides the JSX built-ins used by Booklit's own
// documentation site. It is imported by cmd/booklit-docs so the helpers
// register themselves; the main cmd/booklit binary does not include them.
package booklitdoc

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
	"github.com/vito/booklit/baselit"
	"github.com/vito/booklit/builtins"
)

func init() {
	builtins.Register("OutputFrame", outputFrameFunc)
	builtins.Register("SyntaxHl", syntaxHlFunc)
	builtins.Register("ColumnHeader", columnHeaderFunc)
	builtins.Register("Columns", columnsFunc)
	builtins.Register("Column", columnFunc)
	builtins.Register("LitSyntax", litSyntaxFunc)
	builtins.Register("Godoc", godocFunc)
	builtins.Register("TemplateLink", templateLinkFunc)
	builtins.Register("Define", defineFunc)

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

// outputFrameFunc — `<OutputFrame url="..."/>`. Renders an iframe with
// the given URL as both link and src (via the output-frame template).
func outputFrameFunc(ctx *builtins.Context, props map[string]ast.Node, _ []ast.Node) (booklit.Content, error) {
	url, err := requireStringProp(ctx, props, "url", "OutputFrame")
	if err != nil {
		return nil, err
	}
	return booklit.Styled{
		Style: "output-frame",
		Content: booklit.Link{
			Content: booklit.String(url),
			Target:  url,
		},
		Partials: booklit.Partials{
			"URL": booklit.String(url),
		},
	}, nil
}

// syntaxHlFunc — `<SyntaxHl>content</SyntaxHl>`. Wraps inline content in
// a syntax-highlighting marker; the syntax-hl template provides styling.
func syntaxHlFunc(ctx *builtins.Context, _ map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	content, err := evalChildren(ctx, children)
	if err != nil {
		return nil, err
	}
	return booklit.Styled{
		Style:   "syntax-hl",
		Content: content,
	}, nil
}

// columnHeaderFunc — `<ColumnHeader>content</ColumnHeader>`.
func columnHeaderFunc(ctx *builtins.Context, _ map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	content, err := evalChildren(ctx, children)
	if err != nil {
		return nil, err
	}
	return booklit.Styled{
		Style:   "column-header",
		Block:   true,
		Content: content,
	}, nil
}

// columnsFunc — `<Columns><ColumnHeader>title</ColumnHeader>
// <Column>a</Column><Column>b</Column></Columns>`. Title goes to
// Content; columns become a Partial sequence consumed by columns.tmpl.
// Children are recognized by their AST element name before evaluation
// so the rendered form of ColumnHeader/Column doesn't need to carry a
// wrapper type around.
func columnsFunc(ctx *builtins.Context, _ map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	var title booklit.Content
	var cols booklit.Sequence
	for _, child := range children {
		jsx, ok := child.(ast.JSXElement)
		if !ok {
			continue
		}
		val, err := ctx.Evaluate(ast.Sequence(jsx.Children))
		if err != nil {
			return nil, err
		}
		if val == nil {
			val = booklit.Empty
		}
		switch jsx.Name {
		case "ColumnHeader":
			title = val
		case "Column":
			cols = append(cols, val)
		}
	}
	if title == nil {
		title = booklit.Empty
	}
	return booklit.Styled{
		Style:   "columns",
		Block:   true,
		Content: title,
		Partials: booklit.Partials{
			"Columns": cols,
		},
	}, nil
}

// columnFunc — `<Column>content</Column>`. Outside of <Columns> this
// renders as a plain block; inside, the parent picks it out by name.
func columnFunc(ctx *builtins.Context, _ map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	content, err := evalChildren(ctx, children)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// linkTransformer rewrites `\function-name` occurrences inside
// syntax-highlighted Booklit code into references to the function's
// definition tag (set up by <Define>).
func linkTransformer(sec *booklit.Section) chromaTransformer {
	return chromaTransformer{
		Pattern: regexp.MustCompile(`\\([a-z][a-z0-9-]*)`),
		Transform: func(invoke string) booklit.Content {
			fn := strings.TrimPrefix(invoke, `\`)
			return booklit.Sequence{
				booklit.String(`\`),
				&booklit.Reference{
					Section:  sec,
					TagName:  fn,
					Content:  booklit.String(fn),
					Optional: true,
				},
			}
		},
	}
}

// chromaTransformer mirrors the old chroma.Transformer locally since
// the chroma plugin package is gone.
type chromaTransformer struct {
	Pattern   *regexp.Regexp
	Transform func(string) booklit.Content
}

func (t chromaTransformer) transformAll(str string) booklit.Sequence {
	matches := t.Pattern.FindAllStringIndex(str, -1)
	out := booklit.Sequence{}
	last := 0
	for _, match := range matches {
		if match[0] > last {
			out = append(out, booklit.String(str[last:match[0]]))
		}
		out = append(out, t.Transform(str[match[0]:match[1]]))
		last = match[1]
	}
	if len(str) > last {
		out = append(out, booklit.String(str[last:]))
	}
	return out
}

// syntaxTransform runs chroma syntax highlighting via baselit (which
// uses styles.Fallback — set in init() to the booklitdoc palette) and
// then applies a chain of textual transformers (e.g. linkify function
// names) over the highlighted output. Re-implemented locally so the
// docs aren't blocked on the deleted chroma plugin package.
func syntaxTransform(section *booklit.Section, language string, code booklit.Content, transformers ...chromaTransformer) (booklit.Content, error) {
	base := baselit.NewPlugin(section).(baselit.Plugin)
	highlighted, err := base.Syntax(language, code)
	if err != nil {
		return nil, err
	}

	styled, ok := highlighted.(booklit.Styled)
	if !ok {
		return highlighted, nil
	}
	inner, ok := styled.Content.(booklit.Styled)
	if !ok {
		return highlighted, nil
	}
	raw, ok := inner.Content.(booklit.String)
	if !ok {
		return highlighted, nil
	}

	transformed := booklit.Sequence{booklit.String(string(raw))}
	for _, t := range transformers {
		var next booklit.Sequence
		for _, con := range transformed {
			if s, ok := con.(booklit.String); ok {
				next = append(next, t.transformAll(s.String())...)
			} else {
				next = append(next, con)
			}
		}
		transformed = next
	}
	for i, con := range transformed {
		if _, ok := con.(booklit.String); ok {
			transformed[i] = booklit.Styled{Style: "raw-html", Content: con}
		}
	}
	styled.Content = transformed
	return styled, nil
}

// litSyntaxFunc — `<LitSyntax>code</LitSyntax>`. Highlights Booklit
// source as `lit`, then linkifies `\function` references.
func litSyntaxFunc(ctx *builtins.Context, _ map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	code, err := evalChildren(ctx, children)
	if err != nil {
		return nil, err
	}
	syntax, err := syntaxTransform(ctx.Section, "lit", code, linkTransformer(ctx.Section))
	if err != nil {
		return nil, err
	}
	var style booklit.Style = "lit-block"
	if code.IsFlow() {
		style = "lit-flow"
	}
	return booklit.Styled{Style: style, Content: syntax}, nil
}

// godocFunc — `<Godoc ref="booklit.Section"/>`. Renders a pkg.go.dev
// link styled as code.
func godocFunc(ctx *builtins.Context, props map[string]ast.Node, _ []ast.Node) (booklit.Content, error) {
	ref, err := requireStringProp(ctx, props, "ref", "Godoc")
	if err != nil {
		return nil, err
	}
	spl := strings.SplitN(ref, ".", 2)
	if len(spl) != 2 {
		return nil, fmt.Errorf("<Godoc> ref must be pkg.Symbol, got %q", ref)
	}
	pkg := strings.TrimLeft(spl[0], "*")
	return booklit.Link{
		Content: booklit.Styled{
			Style: booklit.StyleVerbatim,
			Content: booklit.Sequence{
				booklit.String(spl[0] + "."),
				booklit.Styled{Style: booklit.StyleBold, Content: booklit.String(spl[1])},
			},
		},
		Target: "https://pkg.go.dev/github.com/vito/" + pkg + "#" + spl[1],
	}, nil
}

// templateLinkFunc — `<TemplateLink tmpl="italic.tmpl"/>`. Renders a
// link to a bundled HTML template's source on GitHub.
func templateLinkFunc(ctx *builtins.Context, props map[string]ast.Node, _ []ast.Node) (booklit.Content, error) {
	tmpl, err := requireStringProp(ctx, props, "tmpl", "TemplateLink")
	if err != nil {
		return nil, err
	}
	return booklit.Link{
		Content: booklit.Styled{Style: booklit.StyleVerbatim, Content: booklit.String(tmpl)},
		Target:  "https://github.com/vito/booklit/blob/master/render/html/" + tmpl,
	}, nil
}

// defineFunc — `<Define tag="title" sig="<Title>x</Title>">description
// </Define>`. Renders a documentation entry for one component: registers
// the tag as a reference target, displays the syntax-highlighted
// signature, and shows the description.
func defineFunc(ctx *builtins.Context, props map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	tag, err := requireStringProp(ctx, props, "tag", "Define")
	if err != nil {
		return nil, err
	}
	sig, err := requireStringProp(ctx, props, "sig", "Define")
	if err != nil {
		return nil, err
	}
	content, err := evalChildren(ctx, children)
	if err != nil {
		return nil, err
	}

	title, err := syntaxTransform(ctx.Section, "html", booklit.String("<"+componentName(tag)+">"))
	if err != nil {
		return nil, err
	}

	thumb, err := syntaxTransform(ctx.Section, "html", booklit.String(sig))
	if err != nil {
		return nil, err
	}

	return booklit.Styled{
		Style:   "definition",
		Content: content,
		Partials: booklit.Partials{
			"Thumb": booklit.Sequence{
				booklit.Target{
					TagName:  tag,
					Location: ctx.Section.InvokeLocation,
					Title:    title,
					Content:  content,
				},
				thumb,
			},
		},
	}, nil
}

// componentName turns a kebab-case tag (matching the legacy \invoke
// naming, e.g. "table-of-contents") into the PascalCase JSX component
// name ("TableOfContents") used in the rendered signature.
func componentName(tag string) string {
	parts := strings.Split(tag, "-")
	var name strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		name.WriteString(strings.ToUpper(p[:1]))
		name.WriteString(p[1:])
	}
	return name.String()
}

// evalChildren wraps the children list in a Sequence and evaluates it.
func evalChildren(ctx *builtins.Context, children []ast.Node) (booklit.Content, error) {
	if len(children) == 0 {
		return booklit.Empty, nil
	}
	return ctx.Evaluate(ast.Sequence(children))
}

// requireStringProp is a local copy of builtins.requireStringProp;
// the builtins package keeps it unexported.
func requireStringProp(ctx *builtins.Context, props map[string]ast.Node, name, component string) (string, error) {
	node, ok := props[name]
	if !ok {
		return "", fmt.Errorf("<%s> requires prop %q", component, name)
	}
	val, err := ctx.Evaluate(node)
	if err != nil {
		return "", err
	}
	return val.String(), nil
}
