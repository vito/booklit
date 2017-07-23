package stages

import (
	"fmt"
	"reflect"

	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
)

type Evaluate struct {
	Section *booklit.Section

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
	for _, line := range node {
		eval.Result = nil

		err := line.Visit(eval)
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

	if len(para) == 1 && !para[0].IsFlow() {
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
	for _, line := range node {
		eval.Result = nil

		err := line.Visit(eval)
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
	methodName := invoke.Method()

	var method reflect.Value
	for _, p := range eval.Section.Plugins {
		value := reflect.ValueOf(p)
		method = value.MethodByName(methodName)
		if method.IsValid() {
			break
		}
	}

	if !method.IsValid() {
		return fmt.Errorf("undefined booklit function: %s", invoke.Function)
	}

	methodType := method.Type()

	rawArgs := invoke.Arguments

	argc := methodType.NumIn()
	if methodType.IsVariadic() {
		argc--

		if len(rawArgs) < argc {
			return fmt.Errorf("argument count mismatch for %s: given %d, need at least %d", invoke.Function, len(rawArgs), argc)
		}
	} else {
		if len(rawArgs) != argc {
			return fmt.Errorf("argument count mismatch for %s: given %d, need %d", invoke.Function, len(rawArgs), argc)
		}
	}

	argv := make([]reflect.Value, argc)
	for i := 0; i < argc; i++ {
		t := methodType.In(i)
		arg, err := eval.convert(t, rawArgs[i])
		if err != nil {
			return err
		}

		argv[i] = arg
	}

	if methodType.IsVariadic() {
		variadic := rawArgs[argc:]
		variadicType := methodType.In(argc)

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

	switch methodType.NumOut() {
	case 0:
		return nil
	case 1:
		val := result[0].Interface()
		valType := methodType.Out(0)

		switch reflect.New(valType).Interface().(type) {
		case *error:
			if val != nil {
				return val.(error)
			}
		case *booklit.Content:
			eval.Result = booklit.Append(eval.Result, val.(booklit.Content))
		default:
			return fmt.Errorf("unknown return type: %s", valType)
		}
	case 2:
		second := result[1].Interface()
		secondType := methodType.Out(1)
		switch reflect.New(secondType).Interface().(type) {
		case *error:
			if second != nil {
				return second.(error)
			}
		default:
			return fmt.Errorf("unknown second return type: %s", secondType)
		}

		first := result[0].Interface()
		firstType := methodType.Out(0)
		switch reflect.New(firstType).Interface().(type) {
		case *booklit.Content:
			eval.Result = booklit.Append(eval.Result, first.(booklit.Content))
		default:
			return fmt.Errorf("unknown first return type: %s", firstType)
		}
	default:
		return fmt.Errorf("expected 0-2 return values from %s, got %d", invoke.Function, len(result))
	}

	return nil
}

func (eval Evaluate) convert(to reflect.Type, node ast.Node) (reflect.Value, error) {
	switch reflect.New(to).Interface().(type) {
	case *string:
		content, err := eval.evalArg(node)
		if err != nil {
			return reflect.ValueOf(nil), err
		}

		return reflect.ValueOf(content.String()), nil
	case *booklit.Content:
		content, err := eval.evalArg(node)
		if err != nil {
			return reflect.ValueOf(nil), err
		}

		return reflect.ValueOf(content), nil
	case *ast.Node:
		return reflect.ValueOf(node), nil
	default:
		return reflect.ValueOf(nil), fmt.Errorf("unsupported argument type: %s.%s", to.PkgPath(), to.Name())
	}
}

func (eval Evaluate) evalArg(node ast.Node) (booklit.Content, error) {
	subEval := &Evaluate{
		Section: eval.Section,
	}

	err := node.Visit(subEval)
	if err != nil {
		return nil, err
	}

	return subEval.Result, nil
}
