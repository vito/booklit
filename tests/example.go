package tests

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/vito/booklit"
	"github.com/vito/booklit/baselit"
	"github.com/vito/booklit/load"
	"github.com/vito/booklit/render"
)

type Example struct {
	Input       string
	Inputs      Files
	Outputs     Files
	SearchIndex string
	LoadErr     any
	RenderErr   any
	Ext         string // file extension, defaults to ".md"
}

type Files map[string]string

func (example Example) Run() {
	processor := &load.Processor{}

	pluginFactories := []booklit.PluginFactory{
		baselit.NewPlugin,
	}

	dir, err := os.MkdirTemp("", "booklit-tests")
	Expect(err).ToNot(HaveOccurred())

	defer os.RemoveAll(dir)

	ext := example.Ext
	if ext == "" {
		ext = ".md"
	}

	sectionPath := filepath.Join(dir, ginkgo.CurrentSpecReport().LeafNodeText+ext)

	err = os.WriteFile(sectionPath, []byte(example.Input), 0644)
	Expect(err).ToNot(HaveOccurred())

	for file, contents := range example.Inputs {
		err := os.MkdirAll(filepath.Join(dir, filepath.Dir(file)), 0755)
		Expect(err).ToNot(HaveOccurred())

		err = os.WriteFile(filepath.Join(dir, file), []byte(contents), 0644)
		Expect(err).ToNot(HaveOccurred())
	}

	section, err := processor.LoadFile(sectionPath, pluginFactories)
	if example.LoadErr != nil {
		Expect(err).To(MatchError(example.LoadErr))
		return
	}

	Expect(err).ToNot(HaveOccurred())

	engine := render.NewHTMLEngine()

	err = engine.LoadTemplates("fixtures")
	Expect(err).ToNot(HaveOccurred())

	writer := render.Writer{
		Engine:      engine,
		Destination: dir,
	}

	err = writer.WriteSection(section)
	if example.RenderErr != nil {
		Expect(err).To(MatchError(example.RenderErr))
		return
	}

	Expect(err).ToNot(HaveOccurred())

	for file, contents := range example.Outputs {
		fileContents, err := os.ReadFile(filepath.Join(dir, file))
		Expect(err).ToNot(HaveOccurred())
		Expect(string(fileContents)).To(MatchXML(contents))
	}

	if example.SearchIndex != "" {
		err := writer.WriteSearchIndex(section, "search_index.json")
		Expect(err).ToNot(HaveOccurred())

		fileContents, err := os.ReadFile(filepath.Join(dir, "search_index.json"))
		Expect(err).ToNot(HaveOccurred())
		Expect(string(fileContents)).To(MatchJSON(example.SearchIndex))
	}

	Expect(stringifyEverything(section)).ToNot(BeEmpty())
}

// NB: this is really just to cut down on "missing" non-critical test
// coverage. this should recursively stringify all the content.
func stringifyEverything(section *booklit.Section) string {
	var str strings.Builder
	str.WriteString(section.String() + " " + section.Body.String())

	for _, sub := range section.Children {
		str.WriteString(stringifyEverything(sub))
	}

	return str.String()
}
