package processor

import (
	"fmt"
	"reflect"

	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
)

type Processor struct {
	PluginFactories []PluginFactory
}

type PluginFactory interface {
	NewPlugin(*booklit.Section) Plugin
}

type Plugin interface {
	// methods are dynamically invoked
}

func (processor Processor) Load(path string) (*booklit.Section, error) {
	result, err := ast.ParseFile(path)
	if err != nil {
		return nil, err
	}

	node := result.(ast.Node)

	section := &booklit.Section{
		Title: booklit.Empty,
		Body:  booklit.Empty,
	}

	plugins := []Plugin{}
	for _, pf := range processor.PluginFactories {
		plugins = append(plugins, pf.NewPlugin(section))
	}

	evaluator := &Evaluator{
		Plugins: plugins,
		Section: section,
	}

	node.Visit(evaluator)

	section.Body = evaluator.Result

	return section, nil
}

type Evaluator struct {
	Plugins []Plugin

	Result booklit.Content

	Section *booklit.Section
}

func (eval *Evaluator) VisitString(str ast.String) {
	eval.Result = booklit.Append(eval.Result, booklit.String(str))
}

func (eval *Evaluator) VisitSequence(seq ast.Sequence) {
	for _, node := range seq {
		node.Visit(eval)
	}
}

func (eval *Evaluator) VisitInvoke(invoke ast.Invoke) {
	argContent := make([]booklit.Content, len(invoke.Arguments))
	for i, arg := range invoke.Arguments {
		eval := &Evaluator{
			Plugins: eval.Plugins,
			Section: eval.Section,
		}

		arg.Visit(eval)

		argContent[i] = eval.Result
	}

	var method reflect.Value
	for _, p := range eval.Plugins {
		value := reflect.ValueOf(p)
		method = value.MethodByName(invoke.Method)
		if method.IsValid() {
			break
		}
	}

	if !method.IsValid() {
		panic(fmt.Errorf("undefined method: %s", invoke.Method))
	}

	argc := method.Type().NumIn()
	if method.Type().IsVariadic() {
		argc--

		if len(argContent) < argc {
			panic(fmt.Errorf("argument count mismatch for %s: given %d, need at least %d", invoke.Method, len(argContent), argc))
		}
	} else {
		if len(argContent) != argc {
			panic(fmt.Errorf("argument count mismatch for %s: given %d, need %d", invoke.Method, argc, len(argContent)))
		}
	}

	argv := make([]reflect.Value, argc)
	for i := 0; i < argc; i++ {
		t := method.Type().In(i)
		argv[i] = eval.convert(t, argContent[i])
	}

	if method.Type().IsVariadic() {
		variadic := argContent[argc:]
		variadicType := method.Type().In(argc)

		subType := variadicType.Elem()
		for _, arg := range variadic {
			argv = append(argv, eval.convert(subType, arg))
		}
	}

	result := method.Call(argv)
	switch len(result) {
	case 0:
		return
	case 1:
		last := result[0]
		switch v := last.Interface().(type) {
		case error:
			if v != nil {
				panic(v)
			}
		case booklit.Content:
			eval.Result = booklit.Append(eval.Result, v)
		default:
			panic(fmt.Errorf("unknown return type: %T", v))
		}
	case 2:
		first := result[0]
		switch v := first.Interface().(type) {
		case booklit.Content:
			eval.Result = booklit.Append(eval.Result, v)
		default:
			panic(fmt.Errorf("unknown first return type: %T", v))
		}

		last := result[1]
		switch v := last.Interface().(type) {
		case error:
			if v != nil {
				panic(v)
			}
		default:
			panic(fmt.Errorf("unknown second return type: %T", v))
		}
	default:
		panic(fmt.Errorf("expected 0-2 return values from %s, got %d", invoke.Method, len(result)))
	}
}

func (eval Evaluator) convert(to reflect.Type, content booklit.Content) reflect.Value {
	switch reflect.New(to).Interface().(type) {
	case *string:
		switch v := content.(type) {
		case booklit.String:
			return reflect.ValueOf(v.String())
		default:
			panic(fmt.Errorf("cannot satisfy string argument with %T", content))
		}
	case *booklit.Content:
		return reflect.ValueOf(content)
	default:
		panic(fmt.Errorf("unsupported argument type: %s", to))
	}
}
