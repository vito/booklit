package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/vito/booklit"
	"github.com/vito/booklit/baselit"
	"github.com/vito/booklit/load"
	"github.com/vito/booklit/render"

	_ "github.com/vito/booklit/booklitdoc"
)

type Command struct {
	In  string `long:"in"  short:"i" required:"true" description:"Input .lit file."`
	Out string `long:"out" short:"o" required:"true" description:"Output directory in which to render."`

	HTMLEngine struct {
		Templates string `long:"templates" description:"Directory containing .tmpl files to load."`
	} `group:"HTML Rendering Engine" namespace:"html"`
}

func (cmd *Command) Execute(args []string) error {
	processor := &load.Processor{
		PluginFactories: []booklit.PluginFactory{
			booklit.PluginFactoryFunc(baselit.NewPlugin),
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

	engine := render.NewHTMLRenderingEngine()
	err = engine.LoadTemplates(cmd.HTMLEngine.Templates)
	if err != nil {
		return err
	}

	writer := render.Writer{
		Engine:      engine,
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
