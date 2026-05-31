package stages

import (
	"fmt"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
	"github.com/vito/booklit/builtins"
	"github.com/vito/booklit/dangeval"
	"github.com/vito/booklit/templates"
	"github.com/vito/dang/pkg/dang"
)

// Evaluate is an ast.Visitor that builds up the booklit.Content for a given
// section.
type Evaluate struct {
	// The section which acts as the evaluation context. The section's plugins
	// are used for evaluating Invoke nodes.
	Section *booklit.Section

	// Duration after which to log a warning for a slow \invoke.
	SlowInvokeThreshold time.Duration

	// Dang interpreter used to evaluate JSX {expr} interpolations. Nil
	// means no Dang env was bootstrapped; expressions will error.
	Dang *dangeval.Evaluator

	// Templates is the tier-4 mdx-template registry. When a JSX element
	// isn't a built-in and isn't a Dang function, the evaluator looks
	// for an `<dir>/<Name>.md` template and dispatches to it. Nil means
	// no template directory was configured; tier-4 misses silently and
	// the evaluator falls back to the legacy Styled wrap.
	Templates *templates.Registry

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

// VisitJSXElement dispatches a JSX element across three tiers. An
// unknown name is an error rather than an implicit Styled wrap — the
// hard-cutover posture from phase-3b.md Q9 (a).
//
//  1. **Built-in**: a Go function registered in the builtins package.
//  2. **Dang scope**: a `pub PascalCase` function in scope, dispatched
//     with props bridged as named args and children compiled as a
//     `&body` block whose invocation pushes the named args into Dang
//     scope and re-evaluates the children.
//  3. **mdx template**: an `.md` file in the configured templates
//     directory; evaluated with props bound in Dang scope and
//     `children` carrying the JSX children's rendered content.
//
// Dagger dispatch is the eventual fourth tier (see jsx-dang.md).
func (eval *Evaluate) VisitJSXElement(node ast.JSXElement) error {
	eval.Section.InvokeLocation = node.Location

	ctx := &builtins.Context{
		Section:  eval.Section,
		Evaluate: eval.evalArg,
		Dang:     eval.Dang,
	}

	if fn, ok := builtins.Lookup(node.Name); ok {
		content, err := fn(ctx, node.Props, node.Children)
		if err != nil {
			return err
		}
		if content != nil {
			eval.Result = booklit.Append(eval.Result, content)
		}
		return nil
	}

	if eval.Dang != nil {
		callable, found, err := eval.Dang.LookupCallable(node.Name)
		if err != nil {
			return fmt.Errorf("looking up <%s> in Dang scope: %w", node.Name, err)
		}
		if found {
			return eval.dispatchDang(node, callable)
		}
	}

	if eval.Templates != nil && eval.Dang != nil {
		tmpl, found, err := eval.Templates.Load(node.Name)
		if err != nil {
			return fmt.Errorf("loading template for <%s>: %w", node.Name, err)
		}
		if found {
			return eval.dispatchTemplate(node, tmpl)
		}
	}

	return fmt.Errorf("unknown JSX component <%s>: no built-in, no Dang function, no <%s>.md template", node.Name, node.Name)
}

// dispatchDang invokes a Dang function representing a JSX component.
// Props are bridged to Dang values; the children are wrapped in a body
// closure that, each time the Dang function calls it, pushes the named
// args into Dang scope and re-evaluates the children, appending their
// content to a per-invocation accumulator. The accumulated content is
// appended to Result once the Dang function returns.
func (eval *Evaluate) dispatchDang(node ast.JSXElement, callable dang.Callable) error {
	props, err := eval.propsToDang(node.Props)
	if err != nil {
		return err
	}

	var accumulator booklit.Content
	body := func(args map[string]dang.Value) error {
		return eval.Dang.WithBindings(args, func() error {
			content, err := eval.evalArg(ast.Sequence(node.Children))
			if err != nil {
				return err
			}
			if content != nil {
				accumulator = booklit.Append(accumulator, content)
			}
			return nil
		})
	}

	ret, err := eval.Dang.CallComponent(callable, props, body)
	if err != nil {
		return fmt.Errorf("dispatching <%s>: %w", node.Name, err)
	}

	// Body-ful components emit content via the accumulator and ignore
	// their return value. Body-less components return a value directly.
	if accumulator != nil {
		eval.Result = booklit.Append(eval.Result, accumulator)
		return nil
	}
	content, err := eval.Dang.ContentFromValue(ret, eval.Section)
	if err != nil {
		return fmt.Errorf("bridging <%s> return value: %w", node.Name, err)
	}
	if content != nil {
		eval.Result = booklit.Append(eval.Result, content)
	}
	return nil
}

// dispatchTemplate evaluates an mdx template (`<dir>/<Name>.md`) with
// props bound in Dang scope by name and `children` bound to the JSX
// children's evaluated content. The template AST is visited by a sub-
// evaluator so its built-ins, Dang interpolations, and nested
// dispatches run with the bindings in scope.
//
// Children are rendered eagerly (once) into a ContentValue bound as
// `children` — `{children}` and `<Children/>` both read from this
// binding. Templates that don't reference children skip the render
// cost? No: they still pay it. Acceptable for v1; the side-effect
// closure model from phase-3b.md Q5 is a future optimization.
func (eval *Evaluate) dispatchTemplate(node ast.JSXElement, tmpl ast.Node) error {
	props, err := eval.propsToDang(node.Props)
	if err != nil {
		return err
	}

	var childContent booklit.Content
	if len(node.Children) > 0 {
		childContent, err = eval.evalArg(ast.Sequence(node.Children))
		if err != nil {
			return err
		}
	}

	bindings := make(map[string]dang.Value, len(props)+1)
	for k, v := range props {
		bindings[k] = v
	}
	bindings["children"] = dangeval.ContentValue{Content: childContent}

	subEval := &Evaluate{
		Section:             eval.Section,
		SlowInvokeThreshold: eval.SlowInvokeThreshold,
		Dang:                eval.Dang,
		Templates:           eval.Templates,
	}
	err = eval.Dang.WithBindings(bindings, func() error {
		return tmpl.Visit(subEval)
	})
	if err != nil {
		return fmt.Errorf("evaluating template <%s>: %w", node.Name, err)
	}
	if subEval.Result != nil {
		eval.Result = booklit.Append(eval.Result, subEval.Result)
	}
	return nil
}

// propsToDang bridges JSX prop values to Dang values for a component
// call. Literal string attrs map to StringValue verbatim; {expr} attrs
// evaluate against the held Dang env and pass the raw value through;
// anything else evaluates to content and stringifies.
func (eval *Evaluate) propsToDang(props map[string]ast.Node) (map[string]dang.Value, error) {
	out := make(map[string]dang.Value, len(props))
	for name, propNode := range props {
		switch n := propNode.(type) {
		case ast.String:
			out[name] = dang.StringValue{Val: string(n)}
		case ast.JSXExpression:
			val, err := eval.Dang.Eval(n.Raw)
			if err != nil {
				return nil, fmt.Errorf("prop %q: %w", name, err)
			}
			out[name] = val
		default:
			content, err := eval.evalArg(propNode)
			if err != nil {
				return nil, err
			}
			out[name] = dangeval.PropToDang(content)
		}
	}
	return out, nil
}

// VisitJSXExpression parses node.Raw as a Dang snippet, evaluates it
// against the held Dang env, and appends the bridged value to Result.
// Without a Dang evaluator (i.e. embedding contexts that haven't
// bootstrapped one) the expression is an error.
func (eval *Evaluate) VisitJSXExpression(node ast.JSXExpression) error {
	if eval.Dang == nil {
		return fmt.Errorf("no Dang evaluator configured: cannot evaluate {%s}", node.Raw)
	}
	val, err := eval.Dang.Eval(node.Raw)
	if err != nil {
		return fmt.Errorf("evaluating {%s}: %w", node.Raw, err)
	}
	content, err := eval.Dang.ContentFromValue(val, eval.Section)
	if err != nil {
		return fmt.Errorf("bridging {%s}: %w", node.Raw, err)
	}
	if content != nil {
		eval.Result = booklit.Append(eval.Result, content)
	}
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

var complainL = new(sync.Mutex)

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

	if eval.SlowInvokeThreshold > 0 {
		start := time.Now()
		complain := time.AfterFunc(eval.SlowInvokeThreshold, func() {
			complainL.Lock()
			defer complainL.Unlock()
			logrus.WithField("elapsed", time.Since(start)).
				Warn(loc.Annotate("slow invoke: \\%s (still running)", invoke.Function))
			loc.AnnotateLocation(os.Stderr)
		})

		defer func() {
			complainL.Lock()
			defer complainL.Unlock()
			if complain.Stop() {
				logrus.WithField("duration", time.Since(start)).
					Debug(loc.Annotate("fast invoke: \\%s", invoke.Function))
			} else {
				logrus.WithField("duration", time.Since(start)).
					Info(loc.Annotate("slow invoke: \\%s (finished)", invoke.Function))
			}
		}()
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
		Section:             eval.Section,
		SlowInvokeThreshold: eval.SlowInvokeThreshold,
		Dang:                eval.Dang,
		Templates:           eval.Templates,
	}

	err := node.Visit(subEval)
	if err != nil {
		return nil, err
	}

	return subEval.Result, nil
}
