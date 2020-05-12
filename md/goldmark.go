package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func main() {
	md := goldmark.DefaultParser()

	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	out := os.Stdout

	node := md.Parse(text.NewReader(content))

	if os.Getenv("DUMP") != "" {
		node.Dump(content, 0)
	}

	err = ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		switch node := n.(type) {
		case *ast.Document:

		case *ast.Blockquote:
			invoke(out, "quote", entering)

			if !entering {
				fmt.Fprintln(out)
				fmt.Fprintln(out)
			}

		case *ast.List:
			if entering {
				fmt.Fprint(out, `\list`)
			} else {
				fmt.Fprintln(out)
				fmt.Fprintln(out)
			}

		case *ast.ListItem:
			if entering {
				fmt.Fprint(out, `{`)
			} else {
				fmt.Fprint(out, `}`)
			}

		case *ast.Paragraph:
			if entering {
				fmt.Fprintln(out)
			} else {
				fmt.Fprintln(out)
				fmt.Fprintln(out)
			}

		case *ast.Heading:
			var function string
			// if node.IsTitleblock {
			// 	function = "title"
			// } else {
			function = fmt.Sprintf("%sheader", strings.Repeat("sub", node.Level-1))
			// }

			invoke(out, function, entering)

			if !entering {
				// if node.HeadingID != "" {
				// 	fmt.Fprintf(out, "{%s}", node.HeadingID)
				// }

				fmt.Fprintln(out)
				fmt.Fprintln(out)
			}

		case *ast.ThematicBreak:
			invoke(out, "hr", entering)

		case *ast.Emphasis:
			// TODO: is strong level 2?
			invoke(out, "italic", entering)

		case *ast.Link:
			invoke(out, "link", entering)

			if !entering {
				if len(node.Title) != 0 {
					return ast.WalkStop, fmt.Errorf("link titles are not supported by Booklit: %s", string(node.Title))
				}

				fmt.Fprintf(out, `{%s}`, string(node.Destination))
			}

		case *ast.Image:
			invoke(out, fmt.Sprintf(`image{%s}`, string(node.Destination)), entering)

		case *ast.TextBlock:
			// TextBlocks are used in lists which do not contain paragraphs. There is
			// nothing to explicitly do here.

		case *ast.Text:
			if entering {
				fmt.Fprint(out, string(node.Text(content)))

				if node.SoftLineBreak() {
					fmt.Fprintln(out)
				} else if node.HardLineBreak() {
					fmt.Fprint(out, `\break`)
				}
			}

		default:
			return ast.WalkStop, fmt.Errorf("unhandled markdown type: %T", node)
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		panic(err)
	}
}

func invoke(out io.Writer, name string, entering bool) {
	if entering {
		fmt.Fprintf(out, `\%s{`, name)
	} else {
		fmt.Fprintf(out, `}`)
	}
}
