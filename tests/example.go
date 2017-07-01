package tests

import (
	"io/ioutil"
	"path/filepath"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vito/booklit"
	"github.com/vito/booklit/baselit"
	"github.com/vito/booklit/load"
	"github.com/vito/booklit/render"
)

type Example struct {
	Input   string
	Outputs Outputs
}

type Outputs map[string]string

func (example Example) Run() {
	processor := load.Processor{
		PluginFactories: []booklit.PluginFactory{
			baselit.PluginFactory{},
		},
	}

	section, err := processor.LoadSource(ginkgo.CurrentGinkgoTestDescription().TestText, []byte(example.Input))
	Expect(err).ToNot(HaveOccurred())

	dir, err := ioutil.TempDir("", "booklit-tests")
	Expect(err).ToNot(HaveOccurred())

	writer := render.Writer{
		Engine:      render.NewHTMLRenderingEngine(),
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
