package booklitcmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

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

	SaveSearchIndex bool `long:"save-search-index" description:"Save a search index JSON file in the destination."`

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

	if isReexec {
		logrus.SetFormatter(&logrus.JSONFormatter{
			DisableTimestamp: true,
		})
	}

	if cmd.ServerPort != 0 && !isReexec {
		return cmd.Serve()
	} else {
		paths, err := cmd.Build(isReexec)
		if err != nil {
			if reexecErr, ok := err.(ReexecError); ok {
				os.Exit(reexecErr.ExitStatus)
			} else {
				return err
			}
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

	engine := render.NewHTMLRenderingEngine()

	writer := render.Writer{
		Engine:      engine,
		Destination: cmd.Out,
	}

	section, err := processor.LoadFile(cmd.In)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(cmd.Out, 0755)
	if err != nil {
		return nil, err
	}

	if cmd.HTMLEngine.Templates != "" {
		err := engine.LoadTemplates(cmd.HTMLEngine.Templates)
		if err != nil {
			return nil, err
		}
	}

	err = writer.WriteSection(section)
	if err != nil {
		return nil, err
	}

	if cmd.SaveSearchIndex {
		err = writer.WriteSearchIndex(section)
		if err != nil {
			return nil, err
		}
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
		return nil, err
	}

	build := exec.Command("go", "install", src)
	build.Env = append(os.Environ(), "GOBIN="+tmpdir)

	buildOutput, err := build.CombinedOutput()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("compilation failed:\n\n%s", string(buildOutput))
		} else {
			return nil, err
		}
	}

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	cmdProxyR, cmdProxyW := io.Pipe()
	errProxyR, errProxyW := io.Pipe()
	run := exec.Command(bin, os.Args[1:]...)
	run.Env = append(os.Environ(), "BOOKLIT_REEXEC=1")
	run.Stdout = outBuf
	run.Stderr = io.MultiWriter(cmdProxyW, errProxyW)

	cmdLogger := logrus.StandardLogger()

	errLogger := logrus.New()
	errLogger.Formatter = &logrus.TextFormatter{}
	errLogger.Out = errBuf
	errLogger.Level = cmdLogger.Level

	go cmd.proxyLogrus(cmdLogger, os.Stderr, cmdProxyR)
	go cmd.proxyLogrus(errLogger, errBuf, errProxyR)

	err = run.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, ReexecError{
				ExitStatus: exitErr.Sys().(syscall.WaitStatus).ExitStatus(),
				Output:     errBuf.String(),
			}
		} else {
			return nil, err
		}
	}

	var res reexecOutput
	err = json.Unmarshal(outBuf.Bytes(), &res)
	if err != nil {
		return nil, err
	}

	return res.Paths, nil
}

func (cmd *Command) proxyLogrus(logger *logrus.Logger, nonJSON io.Writer, from io.Reader) {
	buf := bufio.NewReader(from)

	var prefix []byte

	for {
		line, isPrefix, err := buf.ReadLine()
		if err == io.EOF {
			return
		}

		if isPrefix {
			prefix = append(prefix, line...)
			continue
		} else {
			line = append(prefix, line...)
			prefix = nil
		}

		fields := logrus.Fields{}
		err = json.Unmarshal(line, &fields)
		if err != nil {
			// not JSON; pass on through
			fmt.Fprintln(nonJSON, string(line))
			continue
		}

		msg := fields["msg"]
		delete(fields, "msg")

		level := fields["level"]
		delete(fields, "level")

		entry := logrus.WithFields(fields)
		entry.Logger = logger

		switch level {
		case "debug":
			entry.Debug(msg)
		case "info":
			entry.Info(msg)
		case "warning":
			entry.Warn(msg)
		case "error":
			entry.Error(msg)
		case "fatal":
			entry.Fatal(msg)
		case "panic":
			entry.Panic(msg)
		}
	}
}

type ReexecError struct {
	ExitStatus int
	Output     string
}

func (err ReexecError) Error() string {
	return fmt.Sprintf("reexec failed (exit status %d); logs:\n\n%s", err.ExitStatus, err.Output)
}
