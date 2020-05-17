package litmd

import (
	"fmt"
	"strings"

	bast "github.com/vito/booklit/ast"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func Parse(content []byte) (bast.Node, error) {
	md := goldmark.DefaultParser()

	inlineParser := parser.NewParser(
		parser.WithInlineParsers(
			parser.DefaultInlineParsers()...,
		),
	)
	inlineParser.AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewInvokeInlineParser(), 100),
	))

	md.AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(NewInvokeBlockParser(), 100),
		),
		parser.WithInlineParsers(
			util.Prioritized(NewInvokeInlineParser(), 100),
		),
	)

	node := md.Parse(text.NewReader(content))

	stack := &stack{}

	var doc bast.Sequence

	var lastInvoke *bast.Invoke

	depth := 0
	err := ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			depth++
		} else {
			depth--
		}

		// fmt.Printf("vvvvvvvvvvvvvvvvvvvvvvvvvvvv WALK %T %v\n", n, entering)
		// n.Dump(content, depth)
		// fmt.Printf("^^^^^^^^^^^^^^^^^^^^^^^^^^^^ WALK %T %v\n", n, entering)

		switch node := n.(type) {
		case *ast.Document:
			if entering {
				stack.push()
			} else {
				doc = stack.pop()
			}

		case *ast.Paragraph:
			if entering {
				stack.push()
			} else {
				stack.append(bast.Paragraph{stack.pop()})
			}

		case *ast.Text:
			if entering {
				stack.append(bast.String(node.Text(content)))

				if node.SoftLineBreak() {
					stack.append(bast.String("\n"))
				} else if node.HardLineBreak() {
					// fmt.Fprint(out, `\break`)
				}
			}

		case *ast.CodeBlock:
			stack.invoke("code", entering)

		case *ast.Blockquote:
			stack.invoke("quote", entering)

		case *ast.List:
			stack.invoke("list", entering)

		case *ast.ListItem:
			if entering {
				stack.push()
			} else {
				stack.append(stack.pop())
			}

		case *ast.Heading:
			if entering {
				stack.push()
			} else {
				stack.append(bast.Invoke{
					Function:  fmt.Sprintf("%sheader", strings.Repeat("sub", node.Level-1)),
					Arguments: stack.pop(),
				})
			}

		case *ast.ThematicBreak:
			// invoke(out, "hr", entering)

		case *ast.Emphasis:
			// TODO: is strong level 2?
			switch node.Level {
			case 1:
				stack.invoke("italic", entering)
			case 2:
				stack.invoke("bold", entering)
			default:
				return ast.WalkStop, fmt.Errorf("unknown emphasis level: %d", node.Level)
			}
			// invoke(out, "italic", entering)

		case *ast.Link:
			// invoke(out, "link", entering)

			// if !entering {
			// 	if len(node.Title) != 0 {
			// 		return ast.WalkStop, fmt.Errorf("link titles are not supported by Booklit: %s", string(node.Title))
			// 	}

			// 	fmt.Fprintf(out, `{%s}`, string(node.Destination))
			// }

		case *ast.Image:
			// invoke(out, fmt.Sprintf(`image{%s}`, string(node.Destination)), entering)

		case *ast.TextBlock:
			// TextBlocks are used in lists which do not contain paragraphs. There is
			// nothing to explicitly do here.

		// case *InvokeBlock:
		// 	if entering {
		// 		stack.push()
		// 	}

		// 	stack.invoke(node.Function, entering)

		// 	if !entering {
		// 		stack.append(bast.Paragraph{stack.pop()})
		// 	}

		case *Invoke:
			if entering {
				lastInvoke = &bast.Invoke{
					Function:  node.Function,
					Arguments: []bast.Node{},
				}

				stack.append(lastInvoke)
			}

		case *InvokeInlineArgument:
			if entering {
				stack.push()
			} else {
				arg := stack.pop()

				if lastInvoke != nil {
					lastInvoke.Arguments = append(lastInvoke.Arguments, arg)
				} else {
					panic("TODO: handle renegade inline arg")
				}
			}

		case *InvokeBlockArgument:
			if entering {
				stack.push()
			} else {
				arg := stack.pop()

				if lastInvoke != nil {
					lastInvoke.Arguments = append(lastInvoke.Arguments, arg)
				} else {
					panic("TODO: handle renegade block arg")
				}
			}

		default:
			return ast.WalkStop, fmt.Errorf("unhandled markdown type: %T", node)
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, err
	}

	return doc, nil
}
