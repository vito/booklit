package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/vito/booklit/baselit"
	"github.com/vito/booklit/load"
	"github.com/vito/booklit/render"
)

type Command struct {
	In  string `long:"in"  short:"i" required:"true" description:"Input .lit file."`
	Out string `long:"out" short:"o" required:"true" description:"Output directory in which to render."`
}

func (cmd *Command) Execute(args []string) error {
	processor := &load.Processor{}
	booklitFactory := baselit.PluginFactory{processor}
	processor.PluginFactories = append(processor.PluginFactories, booklitFactory)

	section, err := processor.LoadFile(cmd.In)
	if err != nil {
		return err
	}

	err = os.MkdirAll(cmd.Out, 0755)
	if err != nil {
		return err
	}

	writer := render.Writer{
		Engine:      render.NewHTMLRenderingEngine(),
		Destination: cmd.Out,
	}

	err = writer.WriteSection(section)
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
