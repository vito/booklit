package tests

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/onsi/ginkgo"
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
	Err         error
}

type Files map[string]string

func (example Example) Run() {
	processor := &load.Processor{}

	pluginFactories := []booklit.PluginFactory{
		baselit.NewPlugin,
	}

	dir, err := ioutil.TempDir("", "booklit-tests")
	Expect(err).ToNot(HaveOccurred())

	defer os.RemoveAll(dir)

	sectionPath := filepath.Join(dir, ginkgo.CurrentGinkgoTestDescription().TestText+".lit")

	err = ioutil.WriteFile(sectionPath, []byte(example.Input), 0644)
	Expect(err).ToNot(HaveOccurred())

	for file, contents := range example.Inputs {
		err := ioutil.WriteFile(filepath.Join(dir, file), []byte(contents), 0644)
		Expect(err).ToNot(HaveOccurred())
	}

	section, err := processor.LoadFile(sectionPath, pluginFactories)
	if example.Err != nil {
		Expect(err).To(MatchError(err))
		return
	}

	Expect(err).ToNot(HaveOccurred())

	engine := render.NewHTMLRenderingEngine()

	err = engine.LoadTemplates("fixtures")
	Expect(err).ToNot(HaveOccurred())

	writer := render.Writer{
		Engine:      engine,
		Destination: dir,
	}

	err = writer.WriteSection(section)
	Expect(err).ToNot(HaveOccurred())

	for file, contents := range example.Outputs {
		fileContents, err := ioutil.ReadFile(filepath.Join(dir, file))
		Expect(err).ToNot(HaveOccurred())
		Expect(string(fileContents)).To(MatchXML(contents))
	}

	if example.SearchIndex != "" {
		err := writer.WriteSearchIndex(section, "search_index.json")
		Expect(err).ToNot(HaveOccurred())

		fileContents, err := ioutil.ReadFile(filepath.Join(dir, "search_index.json"))
		Expect(err).ToNot(HaveOccurred())
		Expect(string(fileContents)).To(MatchJSON(example.SearchIndex))
	}

	Expect(stringifyEverything(section)).ToNot(BeEmpty())
}

// NB: this is really just to cut down on "missing" non-critical test
// coverage. this should recursively stringify all the content.
func stringifyEverything(section *booklit.Section) string {
	str := section.String() + " " + section.Body.String()

	for _, sub := range section.Children {
		str += stringifyEverything(sub)
	}

	return str
}
