package load

import (
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
	"github.com/vito/booklit/stages"
)

type Processor struct {
	AllowBrokenReferences bool

	parsed  map[string]parsedNode
	parsedL sync.Mutex
}

type parsedNode struct {
	Node    ast.Node
	ModTime time.Time
}

func (processor *Processor) LoadFile(path string, pluginFactories []booklit.PluginFactory) (*booklit.Section, error) {
	return processor.LoadFileIn(nil, path, pluginFactories)
}

func (processor *Processor) LoadFileIn(parent *booklit.Section, path string, pluginFactories []booklit.PluginFactory) (*booklit.Section, error) {
	section, err := processor.EvaluateFile(parent, path, pluginFactories)
	if err != nil {
		return nil, err
	}

	return processor.runStages(section)
}

func (processor *Processor) EvaluateFile(parent *booklit.Section, path string, pluginFactories []booklit.PluginFactory) (*booklit.Section, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	modTime := info.ModTime()

	processor.parsedL.Lock()
	if processor.parsed == nil {
		processor.parsed = map[string]parsedNode{}
	}
	parsed, found := processor.parsed[path]
	processor.parsedL.Unlock()

	log := logrus.WithFields(logrus.Fields{
		"path": path,
	})

	var node ast.Node
	if found && !modTime.After(parsed.ModTime) {
		log.Debug("already parsed section")
		node = parsed.Node
	} else {
		log.Debug("parsing section")

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

		node = result.(ast.Node)
	}

	section := &booklit.Section{
		Parent: parent,

		Path: path,

		Title: booklit.Empty,
		Body:  booklit.Empty,

		Processor: processor,
	}

	err = processor.evaluateSection(section, node, pluginFactories)
	if err != nil {
		return nil, err
	}

	processor.parsedL.Lock()
	processor.parsed[path] = parsedNode{
		Node:    node,
		ModTime: modTime,
	}
	processor.parsedL.Unlock()

	return section, nil
}

func (processor *Processor) EvaluateNode(parent *booklit.Section, node ast.Node, pluginFactories []booklit.PluginFactory) (*booklit.Section, error) {
	section := &booklit.Section{
		Parent: parent,

		Title: booklit.Empty,
		Body:  booklit.Empty,

		Processor: processor,
	}

	err := processor.evaluateSection(section, node, pluginFactories)
	if err != nil {
		return nil, err
	}

	return section, nil
}

func (processor *Processor) evaluateSection(section *booklit.Section, node ast.Node, pluginFactories []booklit.PluginFactory) error {
	for _, pf := range pluginFactories {
		section.UsePlugin(pf)
	}

	evaluator := &stages.Evaluate{
		Section: section,
	}

	err := node.Visit(evaluator)
	if err != nil {
		return err
	}

	if evaluator.Result != nil {
		section.Body = evaluator.Result
	}

	return nil
}

func (processor *Processor) runStages(section *booklit.Section) (*booklit.Section, error) {
	collector := &stages.Collect{
		Section: section,
	}

	err := section.Visit(collector)
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
