package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/vito/booklit"
	"github.com/vito/booklit/baselit"
	"github.com/vito/booklit/load"
	"github.com/vito/booklit/render"
	"github.com/vito/booklit/stages"
)

type Command struct {
	In  string `long:"in"  short:"i" required:"true" description:"Input .lit file."`
	Out string `long:"out" short:"o" required:"true" description:"Output directory in which to render."`
}

func (cmd *Command) Execute(args []string) error {
	processor := load.Processor{
		PluginFactories: []booklit.PluginFactory{
			baselit.BaselitPluginFactory{},
		},
	}

	section, err := processor.LoadFile(cmd.In)
	if err != nil {
		return err
	}

	err = os.MkdirAll(cmd.Out, 0755)
	if err != nil {
		return err
	}

	write := stages.Write{
		Engine:      render.NewHTMLRenderingEngine(),
		Destination: cmd.Out,
	}

	err = section.Visit(write)
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
