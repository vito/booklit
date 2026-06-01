package marklit_test

import (
	"sort"
	"testing"

	"github.com/vito/booklit/ast"
	"github.com/vito/booklit/marklit"
)

func TestPlainText(t *testing.T) {
	node := marklit.Parse([]byte("Hello world"))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{ast.String("Hello world")},
	})
}

func TestMultipleParagraphs(t *testing.T) {
	node := marklit.Parse([]byte("First paragraph.\n\nSecond paragraph."))
	assertNode(t, node, ast.Sequence{
		ast.Paragraph{ast.Sequence{ast.String("First paragraph.")}},
		ast.Paragraph{ast.Sequence{ast.String("Second paragraph.")}},
	})
}

func TestEmphasis(t *testing.T) {
	node := marklit.Parse([]byte("Hello *world*"))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.String("Hello "),
			ast.JSXElement{
				Name:     "em",
				Children: []ast.Node{ast.String("world")},
			},
		},
	})
}

func TestBold(t *testing.T) {
	node := marklit.Parse([]byte("Hello **world**"))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.String("Hello "),
			ast.JSXElement{
				Name:     "strong",
				Children: []ast.Node{ast.String("world")},
			},
		},
	})
}

func TestCodeSpan(t *testing.T) {
	node := marklit.Parse([]byte("Use `go fmt` please."))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.String("Use "),
			ast.JSXElement{
				Name:     "code",
				Children: []ast.Node{ast.String("go fmt")},
			},
			ast.String(" please."),
		},
	})
}

func TestLink(t *testing.T) {
	node := marklit.Parse([]byte("[click here](https://example.com)"))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.JSXElement{
				Name:     "a",
				Props:    map[string]ast.Node{"href": ast.String("https://example.com")},
				Children: []ast.Node{ast.String("click here")},
			},
		},
	})
}

func TestReferenceShorthand(t *testing.T) {
	node := marklit.Parse([]byte("[#foo]"))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.JSXElement{
				Name:  "Reference",
				Props: map[string]ast.Node{"tag": ast.String("foo")},
			},
		},
	})
}

func TestReferenceShorthandWithTitle(t *testing.T) {
	node := marklit.Parse([]byte("[Some title](#foo)"))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.JSXElement{
				Name:     "Reference",
				Props:    map[string]ast.Node{"tag": ast.String("foo")},
				Children: []ast.Node{ast.String("Some title")},
			},
		},
	})
}

func TestReferenceShorthandInline(t *testing.T) {
	node := marklit.Parse([]byte("See [#my-section] for details."))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.String("See "),
			ast.JSXElement{
				Name:  "Reference",
				Props: map[string]ast.Node{"tag": ast.String("my-section")},
			},
			ast.String(" for details."),
		},
	})
}

func TestLinkNotReference(t *testing.T) {
	// Regular links (non-# destinations) should still produce <a>
	node := marklit.Parse([]byte("[click here](https://example.com)"))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.JSXElement{
				Name:     "a",
				Props:    map[string]ast.Node{"href": ast.String("https://example.com")},
				Children: []ast.Node{ast.String("click here")},
			},
		},
	})
}

func TestImage(t *testing.T) {
	node := marklit.Parse([]byte("![alt text](image.png)"))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.JSXElement{
				Name: "img",
				Props: map[string]ast.Node{
					"src": ast.String("image.png"),
					"alt": ast.String("alt text"),
				},
			},
		},
	})
}

func TestHeading(t *testing.T) {
	node := marklit.Parse([]byte("# Hello World"))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{ast.JSXElement{
			Name:     "Title",
			Children: []ast.Node{ast.String("Hello World")},
		}},
	})
}

func TestBackslashEscape(t *testing.T) {
	// \\ in Markdown produces a literal backslash. Goldmark preserves the
	// raw \\ in text segments; our converter strips the escape backslash.
	node := marklit.Parse([]byte(`user\\example.com`))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{ast.String(`user\example.com`)},
	})
}

func TestFencedCodeBlock(t *testing.T) {
	input := "```go\nfmt.Println(\"hello\")\n```"
	node := marklit.Parse([]byte(input))

	elem, ok := node.(ast.JSXElement)
	if !ok {
		t.Fatalf("expected JSXElement, got: %T %s", node, nodeString(node))
	}
	if elem.Name != "CodeBlock" {
		t.Fatalf("expected CodeBlock element, got: %s", elem.Name)
	}
	lang, ok := elem.Props["language"].(ast.String)
	if !ok || string(lang) != "go" {
		t.Fatalf("expected language='go' prop, got: %v", elem.Props["language"])
	}
}

func TestFencedCodeBlockNoLang(t *testing.T) {
	input := "```\nhello world\n```"
	node := marklit.Parse([]byte(input))

	elem, ok := node.(ast.JSXElement)
	if !ok {
		t.Fatalf("expected JSXElement, got: %T %s", node, nodeString(node))
	}
	if elem.Name != "pre" {
		t.Fatalf("expected pre element, got: %s", elem.Name)
	}
}

func TestBlockquote(t *testing.T) {
	node := marklit.Parse([]byte("> Hello world"))

	elem, ok := node.(ast.JSXElement)
	if !ok {
		t.Fatalf("expected JSXElement, got: %T %s", node, nodeString(node))
	}
	if elem.Name != "Inset" {
		t.Fatalf("expected Inset element, got: %s", elem.Name)
	}
}

func TestUnorderedList(t *testing.T) {
	input := "- one\n- two\n- three"
	node := marklit.Parse([]byte(input))

	elem, ok := node.(ast.JSXElement)
	if !ok {
		t.Fatalf("expected JSXElement, got: %T %s", node, nodeString(node))
	}
	if elem.Name != "ul" {
		t.Fatalf("expected ul element, got: %s", elem.Name)
	}
	if len(elem.Children) != 3 {
		t.Fatalf("expected 3 items, got %d", len(elem.Children))
	}
}

func TestOrderedList(t *testing.T) {
	input := "1. one\n2. two\n3. three"
	node := marklit.Parse([]byte(input))

	elem, ok := node.(ast.JSXElement)
	if !ok {
		t.Fatalf("expected JSXElement, got: %T %s", node, nodeString(node))
	}
	if elem.Name != "ol" {
		t.Fatalf("expected ol element, got: %s", elem.Name)
	}
	if len(elem.Children) != 3 {
		t.Fatalf("expected 3 items, got %d", len(elem.Children))
	}
}

func TestSoftLineBreak(t *testing.T) {
	// Goldmark produces two Text nodes with a soft line break between them
	// Our converter turns the soft break into a space
	node := marklit.Parse([]byte("Hello\nworld"))
	// Just verify it produces a single paragraph with content
	if node == nil {
		t.Fatal("expected non-nil node")
	}
	str := nodeString(node)
	// Should contain both Hello and world in some form
	if str == "" {
		t.Fatal("expected non-empty result")
	}
	t.Logf("soft line break result: %s", str)
}

func TestHeadingSectionStructuring(t *testing.T) {
	input := "# Top\n\nIntro.\n\n## Sub One\n\nSub content.\n\n## Sub Two\n\nMore content.\n"
	node := marklit.Parse([]byte(input))
	str := nodeString(node)
	// Should produce: Title(Top), body, Section(Title(Sub One) + body), Section(Title(Sub Two) + body)
	expected := "[P([J(Title,S(Top))]) P([S(Intro.)]) P([J(Section,[P([J(Title,S(Sub One))]) P([S(Sub content.)])])]) P([J(Section,[P([J(Title,S(Sub Two))]) P([S(More content.)])])])]"
	if str != expected {
		t.Errorf("AST mismatch:\n  got:  %s\n  want: %s", str, expected)
	}
}

func TestHeadingWithTag(t *testing.T) {
	input := "# Hello World {#hello}\n\nBody.\n"
	node := marklit.Parse([]byte(input))
	str := nodeString(node)
	expected := "[P([J(Title tag=S(hello),S(Hello World))]) P([S(Body.)])]"
	if str != expected {
		t.Errorf("AST mismatch:\n  got:  %s\n  want: %s", str, expected)
	}
}

func TestSubsectionWithTag(t *testing.T) {
	input := "# Top\n\n## Sub {#my-sub}\n\nContent.\n"
	node := marklit.Parse([]byte(input))
	str := nodeString(node)
	expected := "[P([J(Title,S(Top))]) P([J(Section,[P([J(Title tag=S(my-sub),S(Sub))]) P([S(Content.)])])])]"
	if str != expected {
		t.Errorf("AST mismatch:\n  got:  %s\n  want: %s", str, expected)
	}
}

func TestThreeLevelSections(t *testing.T) {
	input := "# Top\n\n## Mid\n\n### Deep\n\nDeep content.\n"
	node := marklit.Parse([]byte(input))
	str := nodeString(node)
	expected := "[P([J(Title,S(Top))]) P([J(Section,[P([J(Title,S(Mid))]) P([J(Section,[P([J(Title,S(Deep))]) P([S(Deep content.)])])])])])]"
	if str != expected {
		t.Errorf("AST mismatch:\n  got:  %s\n  want: %s", str, expected)
	}
}

func TestNoHeadingsUnchanged(t *testing.T) {
	input := "Just a paragraph.\n"
	node := marklit.Parse([]byte(input))
	str := nodeString(node)
	expected := "P([S(Just a paragraph.)])"
	if str != expected {
		t.Errorf("AST mismatch:\n  got:  %s\n  want: %s", str, expected)
	}
}

func TestMarkdownTable(t *testing.T) {
	input := "| A | B |\n| --- | --- |\n| 1 | 2 |\n| 3 | 4 |\n"
	node := marklit.Parse([]byte(input))
	row := func(a, b string) ast.JSXElement {
		return ast.JSXElement{
			Name:      "tr",
			MultiLine: true,
			Children: []ast.Node{
				ast.JSXElement{Name: "td", MultiLine: true, Children: []ast.Node{ast.String(a)}},
				ast.JSXElement{Name: "td", MultiLine: true, Children: []ast.Node{ast.String(b)}},
			},
		}
	}
	assertNode(t, node, ast.JSXElement{
		Name:      "table",
		MultiLine: true,
		Children: []ast.Node{
			row("A", "B"),
			row("1", "2"),
			row("3", "4"),
		},
	})
}

func TestMarkdownTableWithFormatting(t *testing.T) {
	input := "| Name | Status |\n| --- | --- |\n| foo | **ok** |\n"
	node := marklit.Parse([]byte(input))
	str := nodeString(node)
	expected := "J(table,J(tr,J(td,S(Name)),J(td,S(Status))),J(tr,J(td,S(foo)),J(td,J(strong,S(ok)))))"
	if str != expected {
		t.Errorf("AST mismatch:\n  got:  %s\n  want: %s", str, expected)
	}
}

func TestCommentInline(t *testing.T) {
	node := marklit.Parse([]byte("Hello {- comment -}world"))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{ast.String("Hello world")},
	})
}

func TestCommentBlock(t *testing.T) {
	input := "First paragraph.\n\n{- this is\na block comment -}\n\nSecond paragraph.\n"
	node := marklit.Parse([]byte(input))
	str := nodeString(node)
	expected := "[P([S(First paragraph.)]) P([S(Second paragraph.)])]"
	if str != expected {
		t.Errorf("AST mismatch:\n  got:  %s\n  want: %s", str, expected)
	}
}

func TestCommentNested(t *testing.T) {
	node := marklit.Parse([]byte("Hello {- outer {- inner -} still comment -}world"))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{ast.String("Hello world")},
	})
}

func TestCommentUnmatched(t *testing.T) {
	// Unmatched {- is left as-is
	node := marklit.Parse([]byte("Hello {- no end"))
	str := nodeString(node)
	if str == "" {
		t.Fatal("expected non-empty result")
	}
	t.Logf("unmatched comment: %s", str)
}

// assertNode compares a Booklit AST node to an expected value using a
// recursive structural comparison via string representation.
func assertNode(t *testing.T, got, want ast.Node) {
	t.Helper()
	gotStr := nodeString(got)
	wantStr := nodeString(want)
	if gotStr != wantStr {
		t.Errorf("AST mismatch:\n  got:  %s\n  want: %s", gotStr, wantStr)
	}
}

func nodeString(n ast.Node) string {
	v := &stringVisitor{}
	if n == nil {
		return "<nil>"
	}
	_ = n.Visit(v)
	return v.result
}

type stringVisitor struct {
	result string
}

func (v *stringVisitor) VisitString(s ast.String) error {
	v.result += "S(" + string(s) + ")"
	return nil
}

func (v *stringVisitor) VisitSequence(s ast.Sequence) error {
	v.result += "["
	for i, n := range s {
		if i > 0 {
			v.result += " "
		}
		sub := &stringVisitor{}
		_ = n.Visit(sub)
		v.result += sub.result
	}
	v.result += "]"
	return nil
}

func (v *stringVisitor) VisitParagraph(p ast.Paragraph) error {
	v.result += "P("
	for i, line := range p {
		if i > 0 {
			v.result += "|"
		}
		sub := &stringVisitor{}
		_ = line.Visit(sub)
		v.result += sub.result
	}
	v.result += ")"
	return nil
}

func (v *stringVisitor) VisitPreformatted(p ast.Preformatted) error {
	v.result += "Pre("
	for i, line := range p {
		if i > 0 {
			v.result += "|"
		}
		sub := &stringVisitor{}
		_ = line.Visit(sub)
		v.result += sub.result
	}
	v.result += ")"
	return nil
}

func (v *stringVisitor) VisitJSXElement(j ast.JSXElement) error {
	v.result += "J(" + j.Name
	// Sort prop names for stable output across map-iteration order.
	names := make([]string, 0, len(j.Props))
	for k := range j.Props {
		names = append(names, k)
	}
	sort.Strings(names)
	for _, k := range names {
		v.result += " " + k + "="
		sub := &stringVisitor{}
		_ = j.Props[k].Visit(sub)
		v.result += sub.result
	}
	for _, child := range j.Children {
		v.result += ","
		sub := &stringVisitor{}
		_ = child.Visit(sub)
		v.result += sub.result
	}
	v.result += ")"
	return nil
}

func (v *stringVisitor) VisitJSXExpression(e ast.JSXExpression) error {
	v.result += "E(" + e.Raw + ")"
	return nil
}
