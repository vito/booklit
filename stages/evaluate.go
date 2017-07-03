package stages

import (
	"fmt"
	"reflect"

	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
)

type Evaluate struct {
	Plugins []booklit.Plugin

	Result booklit.Content
}

func (eval *Evaluate) VisitString(str ast.String) error {
	eval.Result = booklit.Append(eval.Result, booklit.String(str))
	return nil
}

func (eval *Evaluate) VisitSequence(seq ast.Sequence) error {
	for _, node := range seq {
		err := node.Visit(eval)
		if err != nil {
			return err
		}
	}

	return nil
}

func (eval *Evaluate) VisitParagraph(node ast.Paragraph) error {
	previous := eval.Result

	para := booklit.Paragraph{}
	for _, sentence := range node {
		eval.Result = nil

		err := sentence.Visit(eval)
		if err != nil {
			return err
		}

		if eval.Result != nil {
			para = append(para, eval.Result)
		}
	}

	eval.Result = nil

	if len(para) == 0 {
		// paragraph resulted in no content (e.g. an invoke with no return value)
		eval.Result = previous
		return nil
	}

	if len(para) == 1 && !para[0].IsSentence() {
		// paragraph resulted in block content (e.g. a section)
		eval.Result = booklit.Append(previous, para[0])
		return nil
	}

	eval.Result = booklit.Append(previous, para)

	return nil
}

func (eval *Evaluate) VisitPreformatted(node ast.Preformatted) error {
	previous := eval.Result

	pre := booklit.Preformatted{}
	for _, sentence := range node {
		eval.Result = nil

		err := sentence.Visit(eval)
		if err != nil {
			return err
		}

		if eval.Result != nil {
			pre = append(pre, eval.Result)
		}
	}

	eval.Result = booklit.Append(previous, pre)

	return nil
}

func (eval *Evaluate) VisitInvoke(invoke ast.Invoke) error {
	var method reflect.Value
	for _, p := range eval.Plugins {
		value := reflect.ValueOf(p)
		method = value.MethodByName(invoke.Method)
		if method.IsValid() {
			break
		}
	}

	if !method.IsValid() {
		return fmt.Errorf("undefined method: %s", invoke.Method)
	}

	rawArgs := invoke.Arguments

	argc := method.Type().NumIn()
	if method.Type().IsVariadic() {
		argc--

		if len(rawArgs) < argc {
			return fmt.Errorf("argument count mismatch for %s: given %d, need at least %d", invoke.Method, len(rawArgs), argc)
		}
	} else {
		if len(rawArgs) != argc {
			return fmt.Errorf("argument count mismatch for %s: given %d, need %d", invoke.Method, len(rawArgs), argc)
		}
	}

	argv := make([]reflect.Value, argc)
	for i := 0; i < argc; i++ {
		t := method.Type().In(i)
		arg, err := eval.convert(t, rawArgs[i])
		if err != nil {
			return err
		}

		argv[i] = arg
	}

	if method.Type().IsVariadic() {
		variadic := rawArgs[argc:]
		variadicType := method.Type().In(argc)

		subType := variadicType.Elem()
		for _, varg := range variadic {
			arg, err := eval.convert(subType, varg)
			if err != nil {
				return err
			}

			argv = append(argv, arg)
		}
	}

	result := method.Call(argv)
	switch len(result) {
	case 0:
		return nil
	case 1:
		last := result[0]
		switch v := last.Interface().(type) {
		case nil:
		case error:
			return v
		case booklit.Content:
			eval.Result = booklit.Append(eval.Result, v)
		default:
			return fmt.Errorf("unknown return type: %T", v)
		}
	case 2:
		first := result[0]
		switch v := first.Interface().(type) {
		case booklit.Content:
			eval.Result = booklit.Append(eval.Result, v)
		default:
			return fmt.Errorf("unknown first return type: %T", v)
		}

		last := result[1]
		switch v := last.Interface().(type) {
		case nil:
		case error:
			return v
		default:
			return fmt.Errorf("unknown second return type: %T", v)
		}
	default:
		return fmt.Errorf("expected 0-2 return values from %s, got %d", invoke.Method, len(result))
	}

	return nil
}

func (eval Evaluate) convert(to reflect.Type, node ast.Node) (reflect.Value, error) {
	iface := reflect.New(to).Interface()

	if _, ok := iface.(*ast.Node); ok {
		return reflect.ValueOf(node), nil
	}

	argEval := &Evaluate{
		Plugins: eval.Plugins,
	}

	err := node.Visit(argEval)
	if err != nil {
		return reflect.ValueOf(nil), err
	}

	content := argEval.Result

	switch iface.(type) {
	case *string:
		return reflect.ValueOf(content.String()), nil
	case *booklit.Content:
		return reflect.ValueOf(content), nil
	default:
		return reflect.ValueOf(nil), fmt.Errorf("unsupported argument type: %T", iface)
	}
}
