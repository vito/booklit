package load

import (
	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
	"github.com/vito/booklit/stages"
)

type Processor struct {
	PluginFactories []booklit.PluginFactory
}

func (processor Processor) LoadFile(path string) (*booklit.Section, error) {
	result, err := ast.ParseFile(path)
	if err != nil {
		return nil, err
	}

	return processor.process(result.(ast.Node))
}

func (processor Processor) LoadSource(name string, source []byte) (*booklit.Section, error) {
	result, err := ast.Parse(name, source)
	if err != nil {
		return nil, err
	}

	return processor.process(result.(ast.Node))
}

func (processor Processor) process(node ast.Node) (*booklit.Section, error) {
	section := &booklit.Section{
		Title: booklit.Empty,
		Body:  booklit.Empty,
	}

	plugins := []booklit.Plugin{}
	for _, pf := range processor.PluginFactories {
		plugins = append(plugins, pf.NewPlugin(section))
	}

	evaluator := &stages.Evaluate{
		Plugins: plugins,
	}

	node.Visit(evaluator)

	section.Body = evaluator.Result

	return section, nil
}
