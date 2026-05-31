package builtins

import (
	"fmt"

	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
	"github.com/vito/dang/pkg/dang"
)

func init() {
	Register("For", forFunc)
	Register("If", ifFunc)
	Register("Unless", unlessFunc)
}

// forFunc — `<For each={items} as="item">...</For>`. Evaluates `each`
// as a Dang expression to get a list, then re-evaluates the children
// once per element with the named binding (default `item`) pushed into
// Dang scope. The children's content is concatenated.
//
// Authors who'd otherwise write a per-project `pub Each(items, &body)`
// Dang helper get the same behavior out of the box.
func forFunc(ctx *Context, props map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	if ctx.Dang == nil {
		return nil, fmt.Errorf("<For> requires a Dang evaluator: not configured")
	}

	eachNode, ok := props["each"]
	if !ok {
		return nil, fmt.Errorf("<For> requires prop \"each\"")
	}
	eachExpr, ok := eachNode.(ast.JSXExpression)
	if !ok {
		return nil, fmt.Errorf("<For each> must be a {expr}, got %T", eachNode)
	}

	listVal, err := ctx.Dang.Eval(eachExpr.Raw)
	if err != nil {
		return nil, fmt.Errorf("<For each={%s}>: %w", eachExpr.Raw, err)
	}
	list, ok := listVal.(dang.ListValue)
	if !ok {
		return nil, fmt.Errorf("<For each={%s}> must evaluate to a list, got %T", eachExpr.Raw, listVal)
	}

	bindingName := "item"
	if asNode, ok := props["as"]; ok {
		asContent, err := ctx.Evaluate(asNode)
		if err != nil {
			return nil, fmt.Errorf("<For as>: %w", err)
		}
		bindingName = asContent.String()
	}

	if len(children) == 0 {
		return nil, nil
	}

	var result booklit.Content
	body := ast.Sequence(children)
	for _, item := range list.Elements {
		err := ctx.Dang.WithBindings(map[string]dang.Value{bindingName: item}, func() error {
			content, err := ctx.Evaluate(body)
			if err != nil {
				return err
			}
			if content != nil {
				result = booklit.Append(result, content)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// ifFunc — `<If cond={booleanExpr}>...</If>`. Evaluates `cond` as a
// Dang expression; renders children when the result is truthy
// (non-false, non-null). Returns nothing when falsy.
func ifFunc(ctx *Context, props map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	cond, err := evalCondProp(ctx, props, "If")
	if err != nil {
		return nil, err
	}
	if !cond {
		return nil, nil
	}
	return EvaluateChildren(ctx, children)
}

// unlessFunc — `<Unless cond={booleanExpr}>...</Unless>`. Negation of
// `<If>`; renders children when the condition is falsy.
func unlessFunc(ctx *Context, props map[string]ast.Node, children []ast.Node) (booklit.Content, error) {
	cond, err := evalCondProp(ctx, props, "Unless")
	if err != nil {
		return nil, err
	}
	if cond {
		return nil, nil
	}
	return EvaluateChildren(ctx, children)
}

// evalCondProp evaluates a `cond={expr}` prop to a Go bool. Null and
// false are false; any other value is true.
func evalCondProp(ctx *Context, props map[string]ast.Node, component string) (bool, error) {
	if ctx.Dang == nil {
		return false, fmt.Errorf("<%s> requires a Dang evaluator: not configured", component)
	}
	condNode, ok := props["cond"]
	if !ok {
		return false, fmt.Errorf("<%s> requires prop \"cond\"", component)
	}
	condExpr, ok := condNode.(ast.JSXExpression)
	if !ok {
		return false, fmt.Errorf("<%s cond> must be a {expr}, got %T", component, condNode)
	}
	val, err := ctx.Dang.Eval(condExpr.Raw)
	if err != nil {
		return false, fmt.Errorf("<%s cond={%s}>: %w", component, condExpr.Raw, err)
	}
	switch v := val.(type) {
	case dang.BoolValue:
		return v.Val, nil
	case dang.NullValue:
		return false, nil
	case nil:
		return false, nil
	default:
		return true, nil
	}
}
