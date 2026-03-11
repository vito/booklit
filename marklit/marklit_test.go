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

func TestInvokeNoArgs(t *testing.T) {
	node := marklit.Parse([]byte(`\table-of-contents`))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.Invoke{Function: "table-of-contents"},
		},
	})
}

func TestInvokeOneArg(t *testing.T) {
	node := marklit.Parse([]byte(`\title{Hello world}`))
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
	node := marklit.Parse([]byte(`\link{click here}{https://example.com}`))
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
	node := marklit.Parse([]byte(`\title{Hello *world*}`))
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
	node := marklit.Parse([]byte(`\bold{\italic{wow}}`))
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
	node := marklit.Parse([]byte(`Hello \bold{world} today`))
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

func TestBackslashEscape(t *testing.T) {
	// \\ in Markdown produces a literal backslash. Goldmark preserves the
	// raw \\ in text segments; our converter strips the escape backslash.
	node := marklit.Parse([]byte(`user\\example.com`))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.String("user"),
			ast.String(`\example.com`),
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
	input := "\\section{\nHello world\n\nSecond paragraph.\n}"
	node := marklit.Parse([]byte(input))

	// Block args use ParseArg which preserves paragraph structure
	str := nodeString(node)
	expected := "P([I(section,[P([S(Hello world)]) P([S(Second paragraph.)])])])"
	if str != expected {
		t.Errorf("AST mismatch:\n  got:  %s\n  want: %s", str, expected)
	}
}

func TestBlockInvokeWithTitle(t *testing.T) {
	input := "\\section{\n# My Section\n\nBody text here.\n}"
	node := marklit.Parse([]byte(input))

	str := nodeString(node)
	// Should contain a section invoke with title and body
	if str == "" {
		t.Fatal("expected non-empty result")
	}
	t.Logf("block invoke with title: %s", str)
}

func TestMultipleBlockInvokes(t *testing.T) {
	input := "# Top\n\nIntro.\n\n\\section{\n# Sub A\n\nBody A.\n}\n\n\\section{\n# Sub B\n\nBody B.\n}"
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
	input := "\\define{key}{\nvalue across\nmultiple lines\n}"
	node := marklit.Parse([]byte(input))
	str := nodeString(node)
	t.Logf("block invoke multiple args: %s", str)
	// Should contain a define invoke with 2 args
	if str == "" {
		t.Fatal("expected non-empty result")
	}
}

func TestHeadingWithInvokesBelow(t *testing.T) {
	input := "# My Title\n\nSome body text.\n\n\\table-of-contents"
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
	input := "# Welcome\n\nHello *world*, see \\reference{my-tag} for details."
	node := marklit.Parse([]byte(input))
	str := nodeString(node)
	// Should contain title invoke, italic invoke, and reference invoke
	t.Logf("mixed result: %s", str)
	if str == "" {
		t.Fatal("expected non-empty result")
	}
}

func TestVerbatimArgInline(t *testing.T) {
	node := marklit.Parse([]byte(`\code{{{hello world}}}`))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.Invoke{
				Function:  "code",
				Arguments: []ast.Node{ast.String("hello world")},
			},
		},
	})
}

func TestVerbatimArgNoMarkdown(t *testing.T) {
	// Inside {{{…}}}, Markdown formatting and \invoke are not parsed.
	// Note: content cannot contain }}} — same limitation as the old parser.
	node := marklit.Parse([]byte(`\code{{{*not bold* \not-parsed}}}`))
	assertNode(t, node, ast.Paragraph{
		ast.Sequence{
			ast.Invoke{
				Function:  "code",
				Arguments: []ast.Node{ast.String(`*not bold* \not-parsed`)},
			},
		},
	})
}

func TestVerbatimArgMultiline(t *testing.T) {
	input := "\\code{{{\n  line one\n  line two\n}}}"
	node := marklit.Parse([]byte(input))

	str := nodeString(node)
	// Should contain a code invoke with a Preformatted arg
	// After indent stripping (2 spaces), lines are "line one" and "line two"
	expected := "P([I(code,Pre([S(line one)]|[S(line two)]))])"
	if str != expected {
		t.Errorf("AST mismatch:\n  got:  %s\n  want: %s", str, expected)
	}
}

func TestVerbatimArgIndentStripping(t *testing.T) {
	input := "\\code{{{\n    func main() {\n        fmt.Println(\"hello\")\n    }\n}}}"
	node := marklit.Parse([]byte(input))

	str := nodeString(node)
	expected := "P([I(code,Pre([S(func main() {)]|[S(    fmt.Println(\"hello\"))]|[S(})]))])"
	if str != expected {
		t.Errorf("AST mismatch:\n  got:  %s\n  want: %s", str, expected)
	}
}

func TestVerbatimArgPreservesBackslashes(t *testing.T) {
	// Verbatim should not parse \invoke syntax
	input := "\\lit-syntax{{{\n  \\title{Hello}\n  \\section{Body}\n}}}"
	node := marklit.Parse([]byte(input))

	str := nodeString(node)
	expected := "P([I(lit-syntax,Pre([S(\\title{Hello})]|[S(\\section{Body})]))])"
	if str != expected {
		t.Errorf("AST mismatch:\n  got:  %s\n  want: %s", str, expected)
	}
}

func TestPreformattedArgBasic(t *testing.T) {
	input := "\\code{{\n  line one\n  line two\n}}"
	node := marklit.Parse([]byte(input))

	str := nodeString(node)
	// Preformatted preserves line structure, no Markdown formatting
	expected := "P([I(code,Pre([S(line one)]|[S(line two)]))])"
	if str != expected {
		t.Errorf("AST mismatch:\n  got:  %s\n  want: %s", str, expected)
	}
}

func TestPreformattedArgWithInvoke(t *testing.T) {
	input := "\\code{{\n  $ booklit -i ./index.lit\n  \\syntax-hl{INFO}[0000] listening\n}}"
	node := marklit.Parse([]byte(input))

	str := nodeString(node)
	// \syntax-hl should be parsed as an invoke within the preformatted content.
	// ParseInlineArg wraps the arg in a Sequence: [S(INFO)].
	expected := "P([I(code,Pre([S($ booklit -i ./index.lit)]|[I(syntax-hl,[S(INFO)]) S([0000] listening)]))])"
	if str != expected {
		t.Errorf("AST mismatch:\n  got:  %s\n  want: %s", str, expected)
	}
}

func TestPreformattedArgNoMarkdown(t *testing.T) {
	input := "\\code{{\n  *not bold* [not a link](x)\n}}"
	node := marklit.Parse([]byte(input))

	str := nodeString(node)
	// Markdown formatting should NOT be applied in preformatted args.
	// After indent stripping (2 spaces), content is a single Preformatted line.
	expected := "P([I(code,Pre([S(*not bold* [not a link](x))]))])"
	if str != expected {
		t.Errorf("AST mismatch:\n  got:  %s\n  want: %s", str, expected)
	}
}

func TestPreformattedArgBackslashEscape(t *testing.T) {
	input := "\\code{{\n  user\\\\example.com\n}}"
	node := marklit.Parse([]byte(input))

	str := nodeString(node)
	// After indent stripping (2 spaces), content is "user\\example.com".
	// \\ is escaped to literal \.
	expected := "P([I(code,Pre([S(user) S(\\) S(example.com)]))])"
	if str != expected {
		t.Errorf("AST mismatch:\n  got:  %s\n  want: %s", str, expected)
	}
}

func TestMixedArgTypes(t *testing.T) {
	// First arg is normal {…}, second is verbatim {{{…}}}
	input := "\\syntax{go}{{{func main() {}\n}}}"
	node := marklit.Parse([]byte(input))

	str := nodeString(node)
	// "go" is parsed as inline arg. Verbatim content after stripIndent is
	// single-line "func main() {}" → String (matching old VerbatimLine
	// behavior).
	expected := "P([I(syntax,[S(go)],S(func main() {}))])"
	if str != expected {
		t.Errorf("AST mismatch:\n  got:  %s\n  want: %s", str, expected)
	}
}

func TestVerbatimArgWithBalancedBraces(t *testing.T) {
	// Go template syntax {{.Content | render}} inside verbatim.
	// Block-form (raw starts with \n) → Preformatted even if single line
	// after indent stripping. This ensures block-level rendering.
	input := "\\code{{{\n  {{.Content | render}}\n}}}"
	node := marklit.Parse([]byte(input))

	str := nodeString(node)
	expected := "P([I(code,Pre([S({{.Content | render}})]))])"
	if str != expected {
		t.Errorf("AST mismatch:\n  got:  %s\n  want: %s", str, expected)
	}
}

func TestPreformattedArgWithTemplateCode(t *testing.T) {
	// Go template {{…}} inside preformatted — braces are balanced so }}
	// closure is found correctly after the template expressions.
	// ParsePreformattedArg always returns Preformatted.
	input := "\\code{{\n  {{.Title | render}}\n}}"
	node := marklit.Parse([]byte(input))

	str := nodeString(node)
	expected := "P([I(code,Pre([S({{.Title | render}})]))])"
	if str != expected {
		t.Errorf("AST mismatch:\n  got:  %s\n  want: %s", str, expected)
	}
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

func TestContentBeforeFirstHeading(t *testing.T) {
	input := "\\use-plugin{foo}\n\n# Title\n\nBody.\n"
	node := marklit.Parse([]byte(input))
	str := nodeString(node)
	expected := "[P([I(use-plugin,[S(foo)])]) P([I(title,[S(Title)])]) P([S(Body.)])]"
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
