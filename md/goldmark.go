package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	bast "github.com/vito/booklit/ast"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type stack struct {
	seqs []bast.Sequence
}

func (stack *stack) push() {
	stack.seqs = append(stack.seqs, bast.Sequence{})
}

func (stack *stack) pop() bast.Sequence {
	end := stack.seqs[stack.last()]
	stack.seqs = stack.seqs[0:stack.last()]
	return end
}

func (stack *stack) append(node bast.Node) {
	end := stack.seqs[stack.last()]
	end = append(end, node)
	stack.seqs[stack.last()] = end
}

func (stack *stack) last() int {
	return len(stack.seqs) - 1
}

func (stack *stack) invoke(fun string, entering bool) {
	if entering {
		stack.push()
	} else {
		stack.append(bast.Invoke{
			Function:  fun,
			Arguments: stack.pop(),
		})
	}
}

func NewInvokeParser() parser.InlineParser {
	return &invokeParser{}
}

type invokeParser struct{}

func (s *invokeParser) Trigger() []byte {
	return []byte{'\'', '{', '}'}
}

func (parser *invokeParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	return nil
}

func main() {
	md := goldmark.DefaultParser()
	md.AddOptions(parser.WithInlineParsers(util.Prioritized(NewInvokeParser(), 100)))

	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	node := md.Parse(text.NewReader(content))

	if os.Getenv("DUMP") != "" {
		node.Dump(content, 0)
	}

	stack := &stack{}

	var doc bast.Sequence

	err = ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
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

		default:
			return ast.WalkStop, fmt.Errorf("unhandled markdown type: %T", node)
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		panic(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("  ", "  ")
	enc.Encode(doc)
}
