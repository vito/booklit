package booklitcmd

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/vito/booklit"
	"github.com/vito/booklit/baselit"
	"github.com/vito/booklit/load"
	"github.com/vito/booklit/render"
)

type Command struct {
	Version func() `short:"v" long:"version" description:"Print the version of Boooklit and exit."`

	In  string `long:"in"  short:"i" required:"true" description:"Input .lit file."`
	Out string `long:"out" short:"o" required:"true" description:"Output directory in which to render."`

	Plugins []string `long:"plugin" short:"p" description:"Package to import, providing a plugin."`

	AllowBrokenReferences bool `long:"allow-broken-references" description:"Replace broken references with a bogus tag."`

	HTMLEngine struct {
		Templates string `long:"templates" description:"Directory containing .tmpl files to load."`
	} `group:"HTML Rendering Engine" namespace:"html"`
}

func (cmd *Command) Execute(args []string) error {
	if len(cmd.Plugins) > 0 && os.Getenv("BOOKLIT_REEXEC") == "" {
		return cmd.reexec()
	}

	processor := &load.Processor{
		AllowBrokenReferences: cmd.AllowBrokenReferences,

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

	if cmd.HTMLEngine.Templates != "" {
		err := engine.LoadTemplates(cmd.HTMLEngine.Templates)
		if err != nil {
			return err
		}
	}

	writer := render.Writer{
		Engine:      engine,
		Destination: cmd.Out,
	}

	return writer.WriteSection(section)
}

func (cmd *Command) reexec() error {
	tmpdir, err := ioutil.TempDir("", "booklit-reexec")
	if err != nil {
		return err
	}

	defer os.RemoveAll(tmpdir)

	src := filepath.Join(tmpdir, "main.go")
	bin := filepath.Join(tmpdir, "booklit")

	goSrc := "package main\n"
	goSrc += "import \"github.com/vito/booklit/booklitcmd\"\n"
	for _, p := range cmd.Plugins {
		goSrc += "import _ \"" + p + "\"\n"
	}
	goSrc += "func main() {\n"
	goSrc += "	booklitcmd.Main()\n"
	goSrc += "}\n"

	err = ioutil.WriteFile(src, []byte(goSrc), 0644)
	if err != nil {
		return err
	}

	build := exec.Command("go", "build", "-o", bin, src)
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	err = build.Run()
	if err != nil {
		return err
	}

	run := exec.Command(bin, os.Args[1:]...)
	run.Env = append(os.Environ(), "BOOKLIT_REEXEC=1")
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr
	return run.Run()
}
