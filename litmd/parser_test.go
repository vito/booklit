package litmd_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/vito/booklit/ast"
	"github.com/vito/booklit/litmd"
)

type ParserSuite struct {
	suite.Suite
}

type Example struct {
	Title    string
	Markdown string
	Node     ast.Node
}

func (s *ParserSuite) TestParser() {
	for _, example := range []Example{
		{
			Title:    "basic invoke",
			Markdown: `\basic-invoke`,
			Node: ast.Sequence{
				ast.Paragraph{
					ast.Sequence{
						ast.Invoke{
							Function:  "basic-invoke",
							Arguments: []ast.Node{},
						},
					},
				},
			},
		},
		{
			Title:    "invoke in between words",
			Markdown: `in \between words`,
			Node: ast.Sequence{
				ast.Paragraph{
					ast.Sequence{
						ast.String("in "),
						ast.Invoke{
							Function:  "between",
							Arguments: []ast.Node{},
						},
						ast.String(" words"),
					},
				},
			},
		},
		{
			Title:    "inline argument",
			Markdown: `\inline-arg{Hello. *Goodbye.*}`,
			Node: ast.Sequence{
				ast.Paragraph{
					ast.Sequence{
						ast.Invoke{
							Function: "inline-arg",
							Arguments: []ast.Node{
								ast.Sequence{
									ast.String("Hello. "),
									ast.Invoke{
										Function: "italic",
										Arguments: []ast.Node{
											ast.String("Goodbye."),
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Title:    "multiple inline arguments",
			Markdown: `\inline-arg{Hello.}{*Goodbye.*}`,
			Node: ast.Sequence{
				ast.Paragraph{
					ast.Sequence{
						ast.Invoke{
							Function: "inline-arg",
							Arguments: []ast.Node{
								ast.Sequence{
									ast.String("Hello."),
								},
								ast.Sequence{
									ast.Invoke{
										Function: "italic",
										Arguments: []ast.Node{
											ast.String("Goodbye."),
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Title: "block argument",

			Markdown: `\block-arg{
  Hello.

  Goodbye.
}`,
			Node: ast.Sequence{
				ast.Paragraph{
					ast.Sequence{
						ast.Invoke{
							Function: "block-arg",
							Arguments: []ast.Node{
								ast.Sequence{
									ast.Paragraph{
										ast.Sequence{
											ast.String("Hello."),
										},
									},
									ast.Paragraph{
										ast.Sequence{
											ast.String("Goodbye."),
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Title: "block with multiple arguments",

			Markdown: `\block-arg{Hi!}{
  Hello.

  Goodbye.
}`,
			Node: ast.Sequence{
				ast.Paragraph{
					ast.Sequence{
						ast.Invoke{
							Function: "block-arg",
							Arguments: []ast.Node{
								ast.String("Hi!"),
								ast.Sequence{
									ast.Paragraph{
										ast.Sequence{
											ast.String("Hello."),
										},
									},
									ast.Paragraph{
										ast.Sequence{
											ast.String("Goodbye."),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	} {
		fmt.Println("------------------------------------------------------------------")
		fmt.Println("------------------------------------------------------------------")
		fmt.Println("------------------------------------------------------------------")
		fmt.Println("------------------------------------------------------------------")

		ok := s.Run(example.Title, func() {
			node, err := litmd.Parse([]byte(example.Markdown))
			s.NoError(err)

			s.Equal(example.Node, node)
		})
		if !ok {
			break
		}
	}
}

func TestLitmd(t *testing.T) {
	suite.Run(t, &ParserSuite{})
}
