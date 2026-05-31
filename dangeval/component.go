package dangeval

import (
	"context"
	"fmt"

	"github.com/vito/booklit"
	"github.com/vito/dang/pkg/dang"
	"github.com/vito/dang/pkg/hm"
)

// LookupCallable returns the named binding from the held eval env if it
// is callable (a Dang function, builtin, etc.). The second return is
// false when no binding exists or the binding isn't callable.
func (e *Evaluator) LookupCallable(name string) (dang.Callable, bool, error) {
	val, ok, err := e.evalEnv.Lookup(e.ctx, name)
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	callable, ok := val.(dang.Callable)
	if !ok {
		return nil, false, nil
	}
	return callable, true, nil
}

// CallComponent invokes a Dang function as a JSX component. props are
// supplied by name; body is invoked once per `body(...)` call inside the
// Dang function and is responsible for evaluating the JSX children with
// the named bindings pushed into Dang scope.
//
// The function's return value is returned to the caller. For body-ful
// components content flows out as a side effect of body invocations and
// the return is typically ignored; for body-less components (e.g. a
// helper that just synthesizes a string) the return is the value.
func (e *Evaluator) CallComponent(
	callable dang.Callable,
	props map[string]dang.Value,
	body func(args map[string]dang.Value) error,
) (dang.Value, error) {
	blockVal := &componentBlock{body: body}
	ctx := dang.ContextWithBlock(e.ctx, blockVal)
	return callable.Call(ctx, e.evalEnv, props)
}

// WithBindings runs fn with the given name→value bindings pushed into
// derived eval AND type scopes. Both are needed because {expr} snippets
// re-enter Dang's parser+inferrer; without the type-level binding,
// inference would fail with "name not found". The scopes are popped
// when fn returns.
//
// Single-threaded — relies on Booklit's evaluator running sequentially.
func (e *Evaluator) WithBindings(args map[string]dang.Value, fn func() error) error {
	prevEval := e.evalEnv
	prevType := e.typeEnv

	derivedEval := prevEval.Derive(true)
	derivedType := prevType.Clone()
	for name, val := range args {
		derivedEval.Bind(name, val, dang.PrivateVisibility)
		if t := val.Type(); t != nil {
			derivedType.Add(name, hm.NewScheme(nil, t))
		}
	}
	e.evalEnv = derivedEval
	e.typeEnv = derivedType
	defer func() {
		e.evalEnv = prevEval
		e.typeEnv = prevType
	}()
	return fn()
}

// PropToDang coerces a Booklit content value to a Dang value for use as
// a component prop. Strings round-trip as StringValue; richer content
// stringifies for v1.
func PropToDang(c booklit.Content) dang.Value {
	if c == nil {
		return dang.NullValue{}
	}
	return dang.StringValue{Val: c.String()}
}

// componentBlock implements dang.Callable for a Go-backed body. Each
// invocation from inside a Dang function forwards to the supplied
// closure with the named args.
type componentBlock struct {
	body func(args map[string]dang.Value) error
}

var _ dang.Value = (*componentBlock)(nil)
var _ dang.Callable = (*componentBlock)(nil)

func (b *componentBlock) Type() hm.Type {
	// Runtime-only value; type inference for the call site already
	// completed when the user .dang file was loaded.
	return nil
}

func (b *componentBlock) String() string {
	return "<jsx body>"
}

func (b *componentBlock) Call(_ context.Context, _ dang.EvalEnv, args map[string]dang.Value) (dang.Value, error) {
	if err := b.body(args); err != nil {
		return nil, fmt.Errorf("jsx body: %w", err)
	}
	return dang.NullValue{}, nil
}

func (b *componentBlock) ParameterNames() []string { return nil }
func (b *componentBlock) IsAutoCallable() bool     { return false }
