package litmd

import (
	"fmt"
	"log"
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
		// parser.WithBlockParsers(util.Prioritized(NewBlockInvokeParser(), 100)),
		parser.WithInlineParsers(
			util.Prioritized(NewInvokeInlineParser(), 100),
		),
		parser.WithBlockParsers(
			util.Prioritized(NewInvokeBlockParser(md), 101),
		),
	)

	node := md.Parse(text.NewReader(content))

	node.Dump(content, 0)

	stack := &stack{}

	var doc bast.Sequence

	err := ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		// log.Printf("vvvvvvvvvvvvvvvvvvvvvvvvvvvv WALK %T %v\n", n, entering)
		// n.Dump(content, 0)
		// log.Printf("^^^^^^^^^^^^^^^^^^^^^^^^^^^^ WALK %T %v\n", n, entering)

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

		case *InvokeBlock:
			if entering {
				stack.push()
			}

			stack.invoke(node.Function, entering)

			if !entering {
				stack.append(bast.Paragraph{stack.pop()})
			}

		case *InvokeInline:
			stack.invoke(node.Function, entering)

		case *InvokeInlineArgument:
			log.Println("GEN INVOKE INLINE ARG", entering)
			node.Dump(content, 0)
			stack.dump()

			if entering {
				stack.push()
			} else {
				arg := stack.pop()

				log.Printf("ARG: %#v\n", arg)
				end := stack.seqs[stack.last()]

				if len(end) == 0 {
					stack.append(arg)
				} else if inv, ok := end[0].(bast.Invoke); ok {
					log.Println("ADDING TO ARGS")
					inv.Arguments = append(inv.Arguments, arg)
					end[0] = inv
				} else {
					stack.append(arg)
				}
			}

		case *InvokeBlockArgument:
			log.Println("GEN INVOKE BLOCK ARG", entering)
			node.Dump(content, 0)
			stack.dump()

			if entering {
				stack.push()
			} else {
				arg := stack.pop()

				log.Printf("ARG: %#v\n", arg)
				end := stack.seqs[stack.last()]

				if len(end) == 0 {
					stack.append(arg)
				} else if inv, ok := end[0].(bast.Invoke); ok {
					log.Println("ADDING TO ARGS")
					inv.Arguments = append(inv.Arguments, arg)
					end[0] = inv
				} else {
					stack.append(arg)
				}
			}

			// if entering {
			// 	stack.push()
			// } else {
			// 	stack.append(stack.pop())
			// }

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
