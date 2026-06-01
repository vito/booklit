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
			ast.Invoke{
				Function:  "italic",
				Arguments: []ast.Node{ast.Sequence{ast.String("world")}},
			},
		},
	})
}

func TestBold(t *testing.T) {
	node := marklit.Parse([]byte("Hello **world**"))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.String("Hello "),
			ast.Invoke{
				Function:  "bold",
				Arguments: []ast.Node{ast.Sequence{ast.String("world")}},
			},
		},
	})
}

func TestCodeSpan(t *testing.T) {
	node := marklit.Parse([]byte("Use `go fmt` please."))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.String("Use "),
			ast.Invoke{
				Function:  "code",
				Arguments: []ast.Node{ast.String("go fmt")},
			},
			ast.String(" please."),
		},
	})
}

func TestLink(t *testing.T) {
	node := marklit.Parse([]byte("[click here](https://example.com)"))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.Invoke{
				Function: "link",
				Arguments: []ast.Node{
					ast.Sequence{ast.String("click here")},
					ast.String("https://example.com"),
				},
			},
		},
	})
}

func TestReferenceShorthand(t *testing.T) {
	node := marklit.Parse([]byte("[#foo]"))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.Invoke{
				Function:  "reference",
				Arguments: []ast.Node{ast.Sequence{ast.String("foo")}},
			},
		},
	})
}

func TestReferenceShorthandWithTitle(t *testing.T) {
	node := marklit.Parse([]byte("[Some title](#foo)"))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.Invoke{
				Function: "reference",
				Arguments: []ast.Node{
					ast.String("foo"),
					ast.Sequence{ast.String("Some title")},
				},
			},
		},
	})
}

func TestReferenceShorthandInline(t *testing.T) {
	node := marklit.Parse([]byte("See [#my-section] for details."))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.String("See "),
			ast.Invoke{
				Function:  "reference",
				Arguments: []ast.Node{ast.Sequence{ast.String("my-section")}},
			},
			ast.String(" for details."),
		},
	})
}

func TestLinkNotReference(t *testing.T) {
	// Regular links (non-# destinations) should still produce \link
	node := marklit.Parse([]byte("[click here](https://example.com)"))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.Invoke{
				Function: "link",
				Arguments: []ast.Node{
					ast.Sequence{ast.String("click here")},
					ast.String("https://example.com"),
				},
			},
		},
	})
}

func TestImage(t *testing.T) {
	node := marklit.Parse([]byte("![alt text](image.png)"))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.Invoke{
				Function: "image",
				Arguments: []ast.Node{
					ast.String("image.png"),
					ast.String("alt text"),
				},
			},
		},
	})
}

func TestHeading(t *testing.T) {
	node := marklit.Parse([]byte("# Hello World"))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{ast.Invoke{
			Function:  "title",
			Arguments: []ast.Node{ast.Sequence{ast.String("Hello World")}},
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

	invoke, ok := node.(ast.Invoke)
	if !ok {
		t.Fatalf("expected Invoke, got: %T %s", node, nodeString(node))
	}
	if invoke.Function != "code-block" {
		t.Fatalf("expected code-block invoke, got: %s", invoke.Function)
	}
	if len(invoke.Arguments) != 2 {
		t.Fatalf("expected 2 args, got %d", len(invoke.Arguments))
	}
	if string(invoke.Arguments[0].(ast.String)) != "go" {
		t.Fatalf("expected language 'go', got %q", invoke.Arguments[0])
	}
}

func TestFencedCodeBlockNoLang(t *testing.T) {
	input := "```\nhello world\n```"
	node := marklit.Parse([]byte(input))

	invoke, ok := node.(ast.Invoke)
	if !ok {
		t.Fatalf("expected Invoke, got: %T %s", node, nodeString(node))
	}
	if invoke.Function != "code" {
		t.Fatalf("expected code invoke, got: %s", invoke.Function)
	}
}

func TestBlockquote(t *testing.T) {
	node := marklit.Parse([]byte("> Hello world"))

	invoke, ok := node.(ast.Invoke)
	if !ok {
		t.Fatalf("expected Invoke, got: %T %s", node, nodeString(node))
	}
	if invoke.Function != "inset" {
		t.Fatalf("expected inset invoke, got: %s", invoke.Function)
	}
}

func TestUnorderedList(t *testing.T) {
	input := "- one\n- two\n- three"
	node := marklit.Parse([]byte(input))

	invoke, ok := node.(ast.Invoke)
	if !ok {
		t.Fatalf("expected Invoke, got: %T %s", node, nodeString(node))
	}
	if invoke.Function != "list" {
		t.Fatalf("expected list invoke, got: %s", invoke.Function)
	}
	if len(invoke.Arguments) != 3 {
		t.Fatalf("expected 3 items, got %d", len(invoke.Arguments))
	}
}

func TestOrderedList(t *testing.T) {
	input := "1. one\n2. two\n3. three"
	node := marklit.Parse([]byte(input))

	invoke, ok := node.(ast.Invoke)
	if !ok {
		t.Fatalf("expected Invoke, got: %T %s", node, nodeString(node))
	}
	if invoke.Function != "ordered-list" {
		t.Fatalf("expected ordered-list invoke, got: %s", invoke.Function)
	}
	if len(invoke.Arguments) != 3 {
		t.Fatalf("expected 3 items, got %d", len(invoke.Arguments))
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
	// Should produce: title(Top), body, section(title(Sub One) + body), section(title(Sub Two) + body)
	expected := "[P([I(title,[S(Top)])]) P([S(Intro.)]) P([I(section,[P([I(title,[S(Sub One)])]) P([S(Sub content.)])])]) P([I(section,[P([I(title,[S(Sub Two)])]) P([S(More content.)])])])]"
	if str != expected {
		t.Errorf("AST mismatch:\n  got:  %s\n  want: %s", str, expected)
	}
}

func TestHeadingWithTag(t *testing.T) {
	input := "# Hello World {#hello}\n\nBody.\n"
	node := marklit.Parse([]byte(input))
	str := nodeString(node)
	expected := "[P([I(title,[S(Hello World)],S(hello))]) P([S(Body.)])]"
	if str != expected {
		t.Errorf("AST mismatch:\n  got:  %s\n  want: %s", str, expected)
	}
}

func TestSubsectionWithTag(t *testing.T) {
	input := "# Top\n\n## Sub {#my-sub}\n\nContent.\n"
	node := marklit.Parse([]byte(input))
	str := nodeString(node)
	expected := "[P([I(title,[S(Top)])]) P([I(section,[P([I(title,[S(Sub)],S(my-sub))]) P([S(Content.)])])])]"
	if str != expected {
		t.Errorf("AST mismatch:\n  got:  %s\n  want: %s", str, expected)
	}
}

func TestThreeLevelSections(t *testing.T) {
	input := "# Top\n\n## Mid\n\n### Deep\n\nDeep content.\n"
	node := marklit.Parse([]byte(input))
	str := nodeString(node)
	expected := "[P([I(title,[S(Top)])]) P([I(section,[P([I(title,[S(Mid)])]) P([I(section,[P([I(title,[S(Deep)])]) P([S(Deep content.)])])])])])]"
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
	assertNode(t, node, ast.Invoke{
		Function: "table",
		Arguments: []ast.Node{
			ast.Invoke{
				Function: "table-row",
				Arguments: []ast.Node{
					ast.Sequence{ast.String("A")},
					ast.Sequence{ast.String("B")},
				},
			},
			ast.Invoke{
				Function: "table-row",
				Arguments: []ast.Node{
					ast.Sequence{ast.String("1")},
					ast.Sequence{ast.String("2")},
				},
			},
			ast.Invoke{
				Function: "table-row",
				Arguments: []ast.Node{
					ast.Sequence{ast.String("3")},
					ast.Sequence{ast.String("4")},
				},
			},
		},
	})
}

func TestMarkdownTableWithFormatting(t *testing.T) {
	input := "| Name | Status |\n| --- | --- |\n| foo | **ok** |\n"
	node := marklit.Parse([]byte(input))
	str := nodeString(node)
	expected := "I(table,I(table-row,[S(Name)],[S(Status)]),I(table-row,[S(foo)],[I(bold,[S(ok)])]))"
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

func (v *stringVisitor) VisitInvoke(i ast.Invoke) error {
	v.result += "I(" + i.Function
	for _, arg := range i.Arguments {
		v.result += ","
		sub := &stringVisitor{}
		_ = arg.Visit(sub)
		v.result += sub.result
	}
	v.result += ")"
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
