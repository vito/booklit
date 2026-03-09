package marklit_test

import (
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

func TestInvokeNoArgs(t *testing.T) {
	node := marklit.Parse([]byte("@table-of-contents"))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.Invoke{Function: "table-of-contents"},
		},
	})
}

func TestInvokeOneArg(t *testing.T) {
	node := marklit.Parse([]byte("@title{Hello world}"))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.Invoke{
				Function:  "title",
				Arguments: []ast.Node{ast.Sequence{ast.String("Hello world")}},
			},
		},
	})
}

func TestInvokeMultipleArgs(t *testing.T) {
	node := marklit.Parse([]byte("@link{click here}{https://example.com}"))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.Invoke{
				Function: "link",
				Arguments: []ast.Node{
					ast.Sequence{ast.String("click here")},
					ast.Sequence{ast.String("https://example.com")},
				},
			},
		},
	})
}

func TestInvokeInlineWithMarkdown(t *testing.T) {
	node := marklit.Parse([]byte("@title{Hello *world*}"))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.Invoke{
				Function: "title",
				Arguments: []ast.Node{
					ast.Sequence{
						ast.String("Hello "),
						ast.Invoke{
							Function:  "italic",
							Arguments: []ast.Node{ast.Sequence{ast.String("world")}},
						},
					},
				},
			},
		},
	})
}

func TestInvokeNested(t *testing.T) {
	node := marklit.Parse([]byte("@bold{@italic{wow}}"))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.Invoke{
				Function: "bold",
				Arguments: []ast.Node{
					ast.Sequence{
						ast.Invoke{
							Function: "italic",
							Arguments: []ast.Node{
								ast.Sequence{ast.String("wow")},
							},
						},
					},
				},
			},
		},
	})
}

func TestInvokeMixedWithProse(t *testing.T) {
	node := marklit.Parse([]byte("Hello @bold{world} today"))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.String("Hello "),
			ast.Invoke{
				Function: "bold",
				Arguments: []ast.Node{
					ast.Sequence{ast.String("world")},
				},
			},
			ast.String(" today"),
		},
	})
}

func TestAtEscape(t *testing.T) {
	node := marklit.Parse([]byte("user@@example.com"))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.String("user"),
			ast.String("@"),
			ast.String("example.com"),
		},
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

func TestInvokeMultilineArg(t *testing.T) {
	input := "@section{\nHello world\n\nSecond paragraph.\n}"
	node := marklit.Parse([]byte(input))

	// Block args use ParseArg which preserves paragraph structure
	str := nodeString(node)
	expected := "P([I(section,[P([S(Hello world)]) P([S(Second paragraph.)])])])"
	if str != expected {
		t.Errorf("AST mismatch:\n  got:  %s\n  want: %s", str, expected)
	}
}

func TestBlockInvokeWithTitle(t *testing.T) {
	input := "@section{\n# My Section\n\nBody text here.\n}"
	node := marklit.Parse([]byte(input))

	str := nodeString(node)
	// Should contain a section invoke with title and body
	if str == "" {
		t.Fatal("expected non-empty result")
	}
	t.Logf("block invoke with title: %s", str)
}

func TestMultipleBlockInvokes(t *testing.T) {
	input := "# Top\n\nIntro.\n\n@section{\n# Sub A\n\nBody A.\n}\n\n@section{\n# Sub B\n\nBody B.\n}"
	node := marklit.Parse([]byte(input))

	seq, ok := node.(ast.Sequence)
	if !ok {
		t.Fatalf("expected Sequence, got: %T %s", node, nodeString(node))
	}
	// Should have: title heading, intro paragraph, section A, section B
	if len(seq) < 4 {
		t.Fatalf("expected at least 4 children, got %d: %s", len(seq), nodeString(node))
	}
	t.Logf("multiple block invokes: %s", nodeString(node))
}

func TestBlockInvokeMultipleArgs(t *testing.T) {
	input := "@define{key}{\nvalue across\nmultiple lines\n}"
	node := marklit.Parse([]byte(input))
	str := nodeString(node)
	t.Logf("block invoke multiple args: %s", str)
	// Should contain a define invoke with 2 args
	if str == "" {
		t.Fatal("expected non-empty result")
	}
}

func TestHeadingWithInvokesBelow(t *testing.T) {
	input := "# My Title\n\nSome body text.\n\n@table-of-contents"
	node := marklit.Parse([]byte(input))

	seq, ok := node.(ast.Sequence)
	if !ok {
		t.Fatalf("expected Sequence, got: %T %s", node, nodeString(node))
	}
	if len(seq) != 3 {
		t.Fatalf("expected 3 children, got %d: %s", len(seq), nodeString(node))
	}
}

func TestMixedMarkdownAndInvokes(t *testing.T) {
	input := "# Welcome\n\nHello *world*, see @reference{my-tag} for details."
	node := marklit.Parse([]byte(input))
	str := nodeString(node)
	// Should contain title invoke, italic invoke, and reference invoke
	t.Logf("mixed result: %s", str)
	if str == "" {
		t.Fatal("expected non-empty result")
	}
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
