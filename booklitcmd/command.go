package booklitcmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/vito/booklit"
	"github.com/vito/booklit/baselit"
	"github.com/vito/booklit/load"
	"github.com/vito/booklit/render"
)

type Command struct {
	Version func() `short:"v" long:"version" description:"Print the version of Boooklit and exit."`

	In  string `long:"in"  short:"i" required:"true" description:"Input .lit file to load."`
	Out string `long:"out" short:"o" description:"Directory into which sections will be rendered."`

	SectionTag  string `long:"section-tag"  description:"Section tag to render."`
	SectionPath string `long:"section-path" description:"Section path to load and render with --in as its parent."`

	SaveSearchIndex bool `long:"save-search-index" description:"Save a search index JSON file in the destination."`

	ServerPort int `long:"serve" short:"s" description:"Start an HTTP server on the given port."`

	Plugins []string `long:"plugin" short:"p" description:"Package to import, providing a plugin."`

	Debug bool `long:"debug" short:"d" description:"Log at debug level."`

	AllowBrokenReferences bool `long:"allow-broken-references" description:"Replace broken references with a bogus tag."`

	HTMLEngine struct {
		Templates string `long:"templates" description:"Directory containing .tmpl files to load."`
	} `group:"HTML Rendering Engine" namespace:"html"`

	TextEngine struct {
		FileExtension string `long:"file-extension" description:"File extension to use for generated files."`
		Templates     string `long:"templates"      description:"Directory containing .tmpl files to load."`
	} `group:"Text Rendering Engine" namespace:"text"`
}

func (cmd *Command) Execute(args []string) error {
	if cmd.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	isReexec := os.Getenv("BOOKLIT_REEXEC") != ""
	if !isReexec && len(cmd.Plugins) > 0 {
		logrus.Debug("plugins configured; reexecing")

		exitCode, err := cmd.reexec()
		if err != nil {
			return err
		}

		os.Exit(exitCode)

		return nil
	}

	if cmd.ServerPort != 0 {
		return cmd.Serve()
	}

	return cmd.Build()
}

func (cmd *Command) Serve() error {
	http.Handle("/", &Server{
		In: cmd.In,
		Processor: &load.Processor{
			AllowBrokenReferences: cmd.AllowBrokenReferences,
		},

		Templates:  cmd.HTMLEngine.Templates,
		Engine:     render.NewHTMLEngine(),
		FileServer: http.FileServer(http.Dir(cmd.Out)),
	})

	logrus.WithField("port", cmd.ServerPort).Info("listening")

	return http.ListenAndServe(fmt.Sprintf(":%d", cmd.ServerPort), nil)
}

var basePluginFactories = []booklit.PluginFactory{
	baselit.NewPlugin,
}

func (cmd *Command) Build() error {
	processor := &load.Processor{
		AllowBrokenReferences: cmd.AllowBrokenReferences,
	}

	var engine render.Engine
	if cmd.TextEngine.FileExtension != "" {
		textEngine := render.NewTextEngine(cmd.TextEngine.FileExtension)

		if cmd.TextEngine.Templates != "" {
			err := textEngine.LoadTemplates(cmd.TextEngine.Templates)
			if err != nil {
				return err
			}
		}

		engine = textEngine
	} else {
		htmlEngine := render.NewHTMLEngine()

		if cmd.HTMLEngine.Templates != "" {
			err := htmlEngine.LoadTemplates(cmd.HTMLEngine.Templates)
			if err != nil {
				return err
			}
		}

		engine = htmlEngine
	}

	section, err := processor.LoadFile(cmd.In, basePluginFactories)
	if err != nil {
		return err
	}

	sectionToRender := section
	if cmd.SectionTag != "" {
		tags := section.FindTag(cmd.SectionTag)
		if len(tags) == 0 {
			return fmt.Errorf("unknown tag: %s", cmd.SectionTag)
		}

		sectionToRender = tags[0].Section
	} else if cmd.SectionPath != "" {
		sectionToRender, err = processor.LoadFileIn(section, cmd.SectionPath, basePluginFactories)
		if err != nil {
			return err
		}
	}

	if cmd.Out == "" {
		return engine.RenderSection(os.Stdout, sectionToRender)
	}

	err = os.MkdirAll(cmd.Out, 0755)
	if err != nil {
		return err
	}

	writer := render.Writer{
		Engine:      engine,
		Destination: cmd.Out,
	}

	err = writer.WriteSection(sectionToRender)
	if err != nil {
		return err
	}

	if cmd.SaveSearchIndex {
		err = writer.WriteSearchIndex(section, "search_index.json")
		if err != nil {
			return err
		}
	}

	return nil
}

func (cmd *Command) reexec() (int, error) {
	tmpdir, err := ioutil.TempDir("", "booklit-reexec")
	if err != nil {
		return 0, err
	}

	defer func() {
		_ = os.RemoveAll(tmpdir)
	}()

	src := filepath.Join(tmpdir, "main.go")
	bin := filepath.Join(tmpdir, "main")

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
		return 0, err
	}

	build := exec.Command("go", "install", src)
	build.Env = append(os.Environ(), "GOBIN="+tmpdir)
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr

	logrus.Debug("building reexec binary")

	err = build.Run()
	if err != nil {
		return 0, fmt.Errorf("build failed: %w", err)
	}

	run := exec.Command(bin, os.Args[1:]...)
	run.Env = append(os.Environ(), "BOOKLIT_REEXEC=1")
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr

	logrus.Debug("reexecing")

	err = run.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode(), nil
		}

		return 0, fmt.Errorf("reexec failed: %w", err)
	}

	return 0, nil
}
