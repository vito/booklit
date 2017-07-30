package booklitcmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/build"
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

	In  string `long:"in"  short:"i" required:"true" description:"Input .lit file."`
	Out string `long:"out" short:"o" required:"true" description:"Output directory in which to render."`

	ServerPort int `long:"serve" short:"s" description:"Start an HTTP server on the given port."`

	Plugins []string `long:"plugin" short:"p" description:"Package to import, providing a plugin."`

	Debug bool `long:"debug" short:"d" description:"Log at debug level."`

	AllowBrokenReferences bool `long:"allow-broken-references" description:"Replace broken references with a bogus tag."`

	HTMLEngine struct {
		Templates string `long:"templates" description:"Directory containing .tmpl files to load."`
	} `group:"HTML Rendering Engine" namespace:"html"`
}

func (cmd *Command) Execute(args []string) error {
	if cmd.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	isReexec := os.Getenv("BOOKLIT_REEXEC") != ""

	if cmd.ServerPort != 0 && !isReexec {
		return cmd.Serve()
	} else {
		paths, err := cmd.Build(isReexec)
		if err != nil {
			return err
		}

		if isReexec {
			err := json.NewEncoder(os.Stdout).Encode(reexecOutput{
				Paths: paths,
			})
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func (cmd *Command) Serve() error {
	http.Handle("/", &Server{
		Command:    cmd,
		FileServer: http.FileServer(http.Dir(cmd.Out)),
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", cmd.ServerPort), nil)
}

type reexecOutput struct {
	Paths []string
}

func (cmd *Command) Build(isReexec bool) ([]string, error) {
	if len(cmd.Plugins) > 0 && !isReexec {
		return cmd.reexec()
	}

	processor := &load.Processor{
		AllowBrokenReferences: cmd.AllowBrokenReferences,

		PluginFactories: []booklit.PluginFactory{
			baselit.NewPlugin,
		},
	}

	section, err := processor.LoadFile(cmd.In)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(cmd.Out, 0755)
	if err != nil {
		return nil, err
	}

	engine := render.NewHTMLRenderingEngine()

	if cmd.HTMLEngine.Templates != "" {
		err := engine.LoadTemplates(cmd.HTMLEngine.Templates)
		if err != nil {
			return nil, err
		}
	}

	writer := render.Writer{
		Engine:      engine,
		Destination: cmd.Out,
	}

	err = writer.WriteSection(section)
	if err != nil {
		return nil, err
	}

	return cmd.pathsToWatch(section)
}

func (cmd *Command) pathsToWatch(section *booklit.Section) ([]string, error) {
	paths := cmd.sectionPaths(section)

	for _, plug := range cmd.Plugins {
		pkg, err := build.Import(plug, ".", 0)
		if err != nil {
			return nil, err
		}

		for _, file := range pkg.GoFiles {
			paths = append(paths, filepath.Join(pkg.Dir, file))
		}

		paths = append(paths, pkg.Dir)
	}

	templatesDir := cmd.HTMLEngine.Templates
	if templatesDir != "" {
		files, err := filepath.Glob(filepath.Join(templatesDir, "*.tmpl"))
		if err != nil {
			return nil, err
		}

		paths = append(paths, files...)

		paths = append(paths, templatesDir)
	}

	return paths, nil
}

func (cmd *Command) sectionPaths(section *booklit.Section) []string {
	pathsUniq := map[string]struct{}{section.Path: struct{}{}}

	for _, child := range section.Children {
		for _, path := range cmd.sectionPaths(child) {
			pathsUniq[path] = struct{}{}
		}
	}

	paths := []string{}
	for path, _ := range pathsUniq {
		paths = append(paths, path)
	}

	return paths
}

func (cmd *Command) reexec() ([]string, error) {
	tmpdir, err := ioutil.TempDir("", "booklit-reexec")
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = os.RemoveAll(tmpdir)
	}()

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
		return nil, err
	}

	build := exec.Command("go", "build", "-o", bin, src)

	buildOutput, err := build.CombinedOutput()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("compilation failed:\n\n%s", string(buildOutput))
		} else {
			return nil, err
		}
	}

	buf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	run := exec.Command(bin, os.Args[1:]...)
	run.Env = append(os.Environ(), "BOOKLIT_REEXEC=1")
	run.Stdout = buf
	run.Stderr = errBuf
	err = run.Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return nil, errors.New(errBuf.String())
		} else {
			return nil, err
		}
	}

	var res reexecOutput
	err = json.Unmarshal(buf.Bytes(), &res)
	if err != nil {
		return nil, err
	}

	return res.Paths, nil
}
