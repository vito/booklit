package stages

import (
	"fmt"
	"reflect"

	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
)

// Evaluate is an ast.Visitor that builds up the booklit.Content for a given
// section.
type Evaluate struct {
	// The section which acts as the evaluation context. The section's plugins
	// are used for evaluating Invoke nodes.
	Section *booklit.Section

	// The evaluated content after calling (ast.Node).Visit.
	Result booklit.Content
}

// VisitString appends the string's text to the result using booklit.Append.
func (eval *Evaluate) VisitString(str ast.String) error {
	eval.Result = booklit.Append(eval.Result, booklit.String(str))
	return nil
}

// VisitSequence visits each node within the sequence.
func (eval *Evaluate) VisitSequence(seq ast.Sequence) error {
	for _, node := range seq {
		err := node.Visit(eval)
		if err != nil {
			return err
		}
	}

	return nil
}

// VisitParagraph visits each line in the paragraph and builds up a
// booklit.Paragraph containing each evaluated line.
//
// Any lines which evaluate to a nil result are skipped, i.e. a Invoke which
// performed side effects and returned nothing. If the paragraph is empty as a
// result, the Result is unaffected.
//
// If the Result contains a single non-flow content, it is unwrapped and
// appended to the Result.
//
// Otherwise, i.e. a normal paragraph of flow content, the paragraph is
// appended to the Result.
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

// VisitPreformatted behaves similarly to VisitParagraph, but with no
// special-case for block content, and it appends a booklit.Preformatted
// instead.
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

// VisitInvoke uses reflection to evaluate the corresponding method on the
// section's plugins, trying them in order.
//
// If no method is found, booklit.UndefinedFunctionError is returned.
//
// If the method's arity does not match the Invoke's arguments, an error is
// returned.
//
// The method must return either no value, a booklit.Content, an error, or
// (booklit.Content, error). Other return types will result in an error.
//
// If a PrettyError is returned, it is returned without wrapping.
//
// If an error is returned, it is wrapped in a booklit.FailedFunctionError.
//
// If a booklit.Content is returned, it will be appended to Result.
func (eval *Evaluate) VisitInvoke(invoke ast.Invoke) error {
	eval.Section.InvokeLocation = invoke.Location

	methodName := invoke.Method()

	var method reflect.Value
	for _, p := range eval.Section.Plugins {
		value := reflect.ValueOf(p)
		method = value.MethodByName(methodName)
		if method.IsValid() {
			break
		}
	}

	loc := booklit.ErrorLocation{
		FilePath:     eval.Section.FilePath(),
		NodeLocation: invoke.Location,
		Length:       len("\\" + invoke.Function),
	}

	if !method.IsValid() {
		return booklit.UndefinedFunctionError{
			Function:      invoke.Function,
			ErrorLocation: loc,
		}
	}

	methodType := method.Type()

	rawArgs := invoke.Arguments

	argc := methodType.NumIn()
	if methodType.IsVariadic() {
		argc--

		if len(rawArgs) < argc {
			return fmt.Errorf("argument count mismatch for %s: given %d, need at least %d", invoke.Function, len(rawArgs), argc)
		}
	} else if len(rawArgs) != argc {
		return fmt.Errorf("argument count mismatch for %s: given %d, need %d", invoke.Function, len(rawArgs), argc)
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
				if pErr, ok := val.(booklit.PrettyError); ok {
					return pErr
				}

				return booklit.FailedFunctionError{
					Function:      invoke.Function,
					Err:           val.(error),
					ErrorLocation: loc,
				}
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
				if pErr, ok := second.(booklit.PrettyError); ok {
					return pErr
				}

				return booklit.FailedFunctionError{
					Function:      invoke.Function,
					Err:           second.(error),
					ErrorLocation: loc,
				}
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
		name := to.Name()
		if to.PkgPath() != "" {
			name = to.PkgPath() + "." + name
		}

		return reflect.ValueOf(nil), fmt.Errorf("unsupported argument type: %s", name)
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
