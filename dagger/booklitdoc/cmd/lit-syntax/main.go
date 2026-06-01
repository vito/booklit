// Command lit-syntax is a tiny cgo-enabled helper used by the booklitdoc
// Dagger module. The Go SDK currently builds module runtimes with cgo disabled,
// while tree-sitter's Go bindings require cgo, so the module shells out to this
// helper inside a normal golang container.
package main

import (
	"flag"
	"fmt"
	"os"

	"dagger/booklitdoc/contentjson/wire"
	"dagger/booklitdoc/treehighlight"
)

func main() {
	code := flag.String("code", "", "Booklit source code to highlight")
	language := flag.String("language", "lit", "tree-sitter language to use")
	flag.Parse()

	node, err := litSyntaxNode(*code, *language)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	data, err := wire.Marshal(node)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Stdout.Write(data) //nolint:errcheck
}

func litSyntaxNode(code, language string) (*wire.Node, error) {
	if language == "" {
		language = "lit"
	}
	chunks, err := treehighlight.Chunks(language, code, treehighlight.Options{LinkReferences: true})
	if err != nil {
		return nil, err
	}

	nodes := []*wire.Node{wire.RawHTML(`<pre style=";-webkit-text-size-adjust:none;"><code>`)}
	for _, chunk := range chunks {
		switch {
		case chunk.HTML != "":
			nodes = append(nodes, wire.RawHTML(chunk.HTML))
		case chunk.LinkTag != "":
			nodes = append(nodes, wire.OptionalRef(chunk.LinkTag, wire.String(chunk.LinkText)))
		}
	}
	nodes = append(nodes, wire.RawHTML(`</code></pre>`))

	codeBlock := wire.StyledBlock("code-block", wire.Seq(nodes...))
	codeBlock.Partials = map[string]*wire.Node{"Language": wire.String(language)}
	return wire.StyledBlock("lit-block", codeBlock), nil
}
