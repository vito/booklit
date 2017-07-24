package load

import (
	"os"

	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
	"github.com/vito/booklit/stages"
)

type Processor struct {
	AllowBrokenReferences bool

	PluginFactories []booklit.PluginFactory
}

func (processor *Processor) LoadFile(path string) (*booklit.Section, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	result, err := ast.ParseReader(path, file)
	if err != nil {
		return nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	return processor.loadNode(path, result.(ast.Node))
}

func (processor *Processor) LoadSource(path string, source []byte) (*booklit.Section, error) {
	result, err := ast.Parse(path, source)
	if err != nil {
		return nil, err
	}

	return processor.loadNode(path, result.(ast.Node))
}

func (processor *Processor) EvaluateSection(parent *booklit.Section, path string, node ast.Node) (*booklit.Section, error) {
	section := &booklit.Section{
		Parent: parent,

		Path: path,

		Title: booklit.Empty,
		Body:  booklit.Empty,
	}

	for _, pf := range processor.PluginFactories {
		section.UsePlugin(pf)
	}

	evaluator := &stages.Evaluate{
		Section: section,
	}

	err := node.Visit(evaluator)
	if err != nil {
		return nil, err
	}

	if evaluator.Result != nil {
		section.Body = evaluator.Result
	}

	return section, nil
}

func (processor *Processor) loadNode(path string, node ast.Node) (*booklit.Section, error) {
	section, err := processor.EvaluateSection(nil, path, node)
	if err != nil {
		return nil, err
	}

	collector := &stages.Collect{
		Section: section,
	}

	err = section.Visit(collector)
	if err != nil {
		return nil, err
	}

	resolver := &stages.Resolve{
		AllowBrokenReferences: processor.AllowBrokenReferences,

		Section: section,
	}

	err = section.Visit(resolver)
	if err != nil {
		return nil, err
	}

	return section, nil
}
