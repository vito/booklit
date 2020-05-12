package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/russross/blackfriday/v2"
	"github.com/sirupsen/logrus"
)

func bfmain() {
	md := blackfriday.New(blackfriday.WithExtensions(blackfriday.Titleblock | blackfriday.HeadingIDs))

	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	out := os.Stdout

	node := md.Parse(content)
	node.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		// if entering {
		// 	fmt.Println("-->", node)
		// } else {
		// 	fmt.Println("<--", node)
		// }

		// return blackfriday.GoToNext

		switch node.Type {
		case blackfriday.Document:
		case blackfriday.BlockQuote:
			invoke(out, "quote", entering)

		case blackfriday.List:
			if entering {
				fmt.Fprint(out, `\list`)
			} else {
				fmt.Fprintln(out)
				fmt.Fprintln(out)
			}

		case blackfriday.Item:
			if entering {
				fmt.Fprint(out, `{`)
			} else {
				fmt.Fprint(out, `}`)
			}

		case blackfriday.Paragraph:
			if !entering {
				fmt.Fprintln(out)
				fmt.Fprintln(out)
			}

		case blackfriday.Heading:
			var function string
			if node.IsTitleblock {
				function = "title"
			} else {
				function = fmt.Sprintf("%sheader", strings.Repeat("sub", node.HeadingData.Level-1))
			}

			invoke(out, function, entering)

			if !entering {
				if node.HeadingID != "" {
					fmt.Fprintf(out, "{%s}", node.HeadingID)
				}

				fmt.Fprintln(out)
				fmt.Fprintln(out)
			}

		case blackfriday.HorizontalRule:
			invoke(out, "hr", entering)

		case blackfriday.Emph:
			invoke(out, "italic", entering)

		case blackfriday.Strong:
			invoke(out, "strong", entering)

		case blackfriday.Del:
			invoke(out, "del", entering)

		case blackfriday.Link:
			invoke(out, "link", entering)

			if !entering {
				fmt.Fprintf(out, `{%s}`, string(node.Destination))

				if len(node.Title) > 0 {
					logrus.Warn("link titles are not supported by Booklit")
				}
			}

		case blackfriday.Image:

		case blackfriday.Text:
			fmt.Fprint(out, string(node.Literal))

		case blackfriday.HTMLBlock:
		case blackfriday.CodeBlock:
		case blackfriday.Softbreak:
		case blackfriday.Hardbreak:
		case blackfriday.Code:
		case blackfriday.HTMLSpan:
		case blackfriday.Table:
		case blackfriday.TableCell:
		case blackfriday.TableHead:
		case blackfriday.TableBody:
		case blackfriday.TableRow:
		}

		return blackfriday.GoToNext
	})
}

func invoke2(out io.Writer, name string, entering bool) {
	if entering {
		fmt.Fprintf(out, `\%s{`, name)
	} else {
		fmt.Fprintf(out, `}`)
	}
}
