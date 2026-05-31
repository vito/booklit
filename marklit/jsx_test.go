package marklit_test

import (
	"testing"

	"github.com/vito/booklit/ast"
	"github.com/vito/booklit/marklit"
)

func TestJSX(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want ast.Node
	}{
		{
			name: "self-closing, no attrs",
			in:   `<Foo/>`,
			want: jsx("Foo", nil, nil),
		},
		{
			name: "open/close, empty body",
			in:   `<Foo></Foo>`,
			want: jsx("Foo", nil, nil),
		},
		{
			name: "literal text body",
			in:   `<Foo>hello</Foo>`,
			want: jsx("Foo", nil, []ast.Node{ast.String("hello")}),
		},
		{
			name: "string attr",
			in:   `<Foo bar="x"/>`,
			want: jsx("Foo", map[string]ast.Node{"bar": ast.String("x")}, nil),
		},
		{
			name: "multiple attrs",
			in:   `<Foo bar="x" baz="y"/>`,
			want: jsx("Foo", map[string]ast.Node{
				"bar": ast.String("x"),
				"baz": ast.String("y"),
			}, nil),
		},
		{
			name: "camelCase attr preserved",
			in:   `<Card iconName="star"/>`,
			want: jsx("Card", map[string]ast.Node{"iconName": ast.String("star")}, nil),
		},
		{
			name: "expression attr",
			in:   `<Foo bar={x}/>`,
			want: jsx("Foo", map[string]ast.Node{
				"bar": ast.JSXExpression{Raw: "x"},
			}, nil),
		},
		{
			name: "expression child",
			in:   `<Foo>{y}</Foo>`,
			want: jsx("Foo", nil, []ast.Node{ast.JSXExpression{Raw: "y"}}),
		},
		{
			name: "nested JSX child",
			in:   `<Foo><Bar/></Foo>`,
			want: jsx("Foo", nil, []ast.Node{jsx("Bar", nil, nil)}),
		},
		{
			name: "mixed text and JSX",
			in:   `<Foo>hello <Bar/> world</Foo>`,
			// Text chunks between JSX get re-parsed as inline markdown,
			// which produces separate ast.String nodes for words and the
			// trailing/leading whitespace that ParseInlineArg restores.
			want: jsx("Foo", nil, []ast.Node{
				ast.String("hello"),
				ast.String(" "),
				jsx("Bar", nil, nil),
				ast.String(" "),
				ast.String("world"),
			}),
		},
		{
			name: "markdown inside children",
			in:   `<Foo>hello *world*</Foo>`,
			want: jsx("Foo", nil, []ast.Node{
				ast.String("hello "),
				ast.Invoke{
					Function:  "italic",
					Arguments: []ast.Node{ast.Sequence{ast.String("world")}},
				},
			}),
		},
		{
			name: "JSX inline within sentence",
			in:   `Hello <Foo/> world`,
			want: ast.Paragraph{ast.Sequence{
				ast.String("Hello "),
				jsx("Foo", nil, nil),
				ast.String(" world"),
			}},
		},
		{
			name: "multi-line element body",
			in:   "<Foo>\n  body\n</Foo>",
			// Soft-line-break text is collapsed by ParseInlineArg.
			want: jsx("Foo", nil, []ast.Node{ast.String("body")}),
		},
		{
			name: "multi-line attrs",
			in:   "<Foo\n  bar=\"x\"\n/>",
			want: jsx("Foo", map[string]ast.Node{"bar": ast.String("x")}, nil),
		},
		{
			name: "brace expression with quoted string",
			in:   `<Foo bar={"}"}/>`,
			want: jsx("Foo", map[string]ast.Node{
				"bar": ast.JSXExpression{Raw: `"}"`},
			}, nil),
		},
		{
			name: "nested braces in expression",
			in:   `<Foo bar={{a: 1}}/>`,
			want: jsx("Foo", map[string]ast.Node{
				"bar": ast.JSXExpression{Raw: "{a: 1}"},
			}, nil),
		},
		{
			name: "string attr with escape",
			in:   `<Foo bar="a\"b"/>`,
			want: jsx("Foo", map[string]ast.Node{"bar": ast.String(`a\"b`)}, nil),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := marklit.Parse([]byte(tc.in))
			assertNode(t, got, tc.want)
		})
	}
}

// TestJSXLowercaseFallsThrough verifies that lowercase tag names are NOT
// claimed by the JSX parser and remain available to goldmark's raw-HTML
// handling.
func TestJSXLowercaseFallsThrough(t *testing.T) {
	// Lowercase <div> should pass through; the resulting AST should NOT
	// contain a JSXElement.
	got := marklit.Parse([]byte(`<div>hi</div>`))
	if containsJSX(got) {
		t.Errorf("lowercase <div> was claimed by JSX parser; got %s", nodeString(got))
	}
}

// TestJSXUnmatchedLessThanFallsThrough verifies that a bare '<' (not followed
// by an uppercase letter) does not trigger JSX parsing.
func TestJSXUnmatchedLessThanFallsThrough(t *testing.T) {
	got := marklit.Parse([]byte(`a < b`))
	if containsJSX(got) {
		t.Errorf("'<' followed by space was claimed by JSX parser; got %s", nodeString(got))
	}
}

// TestJSXMalformedRollsBack verifies that a partial-looking JSX that fails
// to parse (e.g., missing close tag) rolls back so goldmark can handle the
// '<' as text.
func TestJSXMalformedRollsBack(t *testing.T) {
	// `<Foo` with no `>` and no close — should not become a JSXElement.
	got := marklit.Parse([]byte(`<Foo`))
	if containsJSX(got) {
		t.Errorf("malformed `<Foo` was claimed by JSX parser; got %s", nodeString(got))
	}
}

// jsx is a convenience for building expected ast.JSXElement values in tests.
func jsx(name string, props map[string]ast.Node, children []ast.Node) ast.JSXElement {
	if props == nil {
		props = map[string]ast.Node{}
	}
	return ast.JSXElement{
		Name:     name,
		Props:    props,
		Children: children,
	}
}

// containsJSX returns true if any JSXElement appears in the tree.
func containsJSX(n ast.Node) bool {
	v := &jsxDetector{}
	_ = n.Visit(v)
	return v.found
}

type jsxDetector struct{ found bool }

func (v *jsxDetector) VisitString(ast.String) error           { return nil }
func (v *jsxDetector) VisitInvoke(i ast.Invoke) error {
	for _, a := range i.Arguments {
		_ = a.Visit(v)
	}
	return nil
}
func (v *jsxDetector) VisitSequence(s ast.Sequence) error {
	for _, n := range s {
		_ = n.Visit(v)
	}
	return nil
}
func (v *jsxDetector) VisitParagraph(p ast.Paragraph) error {
	for _, line := range p {
		_ = line.Visit(v)
	}
	return nil
}
func (v *jsxDetector) VisitPreformatted(p ast.Preformatted) error {
	for _, line := range p {
		_ = line.Visit(v)
	}
	return nil
}
func (v *jsxDetector) VisitJSXElement(ast.JSXElement) error {
	v.found = true
	return nil
}
func (v *jsxDetector) VisitJSXExpression(ast.JSXExpression) error { return nil }
