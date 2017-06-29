package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	flags "github.com/jessevdk/go-flags"
	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
	"github.com/vito/booklit/render"
)

type Command struct {
	In  string `long:"in"  short:"i" required:"true" description:"Input .lit file."`
	Out string `long:"out" short:"o" required:"true" description:"Output directory in which to render."`
}

func (cmd *Command) Execute(args []string) error {
	source := `hi

im a thing

\title{Hello, \italic{world}!}{hello}

how are you?
`

	node, err := ast.Parse("test", []byte(source))
	if err != nil {
		return fmt.Errorf("failed to parse: %s", err)
	}

	json.NewEncoder(os.Stdout).Encode(node)

	section := &booklit.Section{
		Title: booklit.String("Hello, world!"),
		Body: booklit.Sequence([]booklit.Content{
			booklit.String("hi"),
		}),
	}

	// file, err := os.Open(cmd.In)
	// if err != nil {
	// 	return err
	// }

	err = os.MkdirAll(cmd.Out, 0755)
	if err != nil {
		return err
	}

	out, err := os.Create(filepath.Join(cmd.Out, "test.html"))
	if err != nil {
		return err
	}

	defer out.Close()

	engine := render.NewHTMLRenderingEngine()

	section.Visit(engine)

	err = engine.Render(out)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	cmd := &Command{}

	parser := flags.NewParser(cmd, flags.Default)
	parser.NamespaceDelimiter = "-"

	args, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}

	err = cmd.Execute(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
