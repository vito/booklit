package tests

import (
	"io/ioutil"
	"path/filepath"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vito/booklit/baselit"
	"github.com/vito/booklit/load"
	"github.com/vito/booklit/render"
)

type Example struct {
	Input   string
	Inputs  Files
	Outputs Files
}

type Files map[string]string

func (example Example) Run() {
	processor := &load.Processor{}
	baselitFactory := baselit.PluginFactory{processor}
	processor.PluginFactories = append(processor.PluginFactories, baselitFactory)

	dir, err := ioutil.TempDir("", "booklit-tests")
	Expect(err).ToNot(HaveOccurred())

	for file, contents := range example.Inputs {
		err := ioutil.WriteFile(filepath.Join(dir, file), []byte(contents), 0644)
		Expect(err).ToNot(HaveOccurred())
	}

	fakePath := filepath.Join(dir, ginkgo.CurrentGinkgoTestDescription().TestText+".lit")
	section, err := processor.LoadSource(fakePath, []byte(example.Input))
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
}
