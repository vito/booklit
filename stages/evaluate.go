package stages

import (
	"fmt"
	"sort"

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
	// The section which acts as the evaluation context.
	Section *booklit.Section

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

// VisitParagraph visits each line of the paragraph node, then segments
// the evaluated lines into runs of flow content (each wrapped in a
// booklit.Paragraph) interleaved with block content (emitted
// unwrapped). Mirrors CommonMark's behavior for block HTML embedded in
// a paragraph: preceding flow becomes its own `<p>`, the block element
// breaks out, trailing flow forms a new `<p>`.
//
// A Sequence line whose items are mixed flow + block (typical when an
// inline JSX element evaluates to block content mid-prose) splits
// inside the line at the block boundary. A line that's entirely flow
// stays one line of the surrounding Paragraph. A line that's entirely
// block flushes the flow buffer and emits standalone.
//
// Lines that evaluate to nil are skipped (an invoke that returned
// nothing). If the paragraph as a whole emits nothing, Result is
// unaffected.
func (eval *Evaluate) VisitParagraph(node ast.Paragraph) error {
	previous := eval.Result

	var lines []booklit.Content
	for _, line := range node {
		eval.Result = nil

		err := line.Visit(eval)
		if err != nil {
			return err
		}

		if eval.Result != nil {
			lines = append(lines, eval.Result)
		}
	}

	eval.Result = previous

	if len(lines) == 0 {
		return nil
	}

	flow := booklit.Paragraph{}
	flushFlow := func() {
		if len(flow) == 0 {
			return
		}
		eval.Result = booklit.Append(eval.Result, flow)
		flow = booklit.Paragraph{}
	}
	emit := func(item booklit.Content) {
		if item.IsFlow() {
			flow = append(flow, item)
			return
		}
		flushFlow()
		eval.Result = booklit.Append(eval.Result, item)
	}

	for _, line := range lines {
		// A non-flow Sequence (mixed flow + block items from an inline
		// JSX element evaluating to block content mid-prose) splits at
		// the block boundary; the contiguous flow before and after stays
		// inside the paragraph being built.
		if seq, ok := line.(booklit.Sequence); ok && !seq.IsFlow() {
			for _, item := range seq {
				emit(item)
			}
			continue
		}
		emit(line)
	}
	flushFlow()

	return nil
}

// VisitJSXElement dispatches a JSX element. Tags are split on case:
//
//   - **lowercase** tags are treated as raw HTML wrappers: children
//     are evaluated normally, attribute values are stringified (with
//     HTML-escaping), and the result is sandwiched between literal
//     `<name attrs>` / `</name>` opening and closing fragments.
//
//   - **PascalCase** tags route through three tiers in order: a Go
//     builtin in `builtins/`; a `pub PascalCase` function in Dang
//     scope; an mdx-template `.md` file. An unmatched PascalCase tag
//     is an error rather than an implicit Styled wrap.
func (eval *Evaluate) VisitJSXElement(node ast.JSXElement) error {
	eval.Section.InvokeLocation = node.Location

	if isLowerCaseTag(node.Name) {
		return eval.dispatchRawHTML(node)
	}

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

func isLowerCaseTag(name string) bool {
	if name == "" {
		return false
	}
	c := name[0]
	return c >= 'a' && c <= 'z'
}

// dispatchRawHTML emits a booklit.RawElement for a lowercase JSX
// element. Attributes are pre-rendered (string-valued props splice in
// verbatim, expression-valued props evaluate against the Dang env and
// stringify with HTML-escaping). A node with no children produces a
// RawElement with nil Content, which the renderer emits self-closing.
//
// Block/flow classification now comes from the tag name (via
// internal/htmltags), not from MultiLine — `<div>` is always block,
// `<span>` is always flow, regardless of how the source spans lines.
// The previous design used a Block flag fed from MultiLine to keep
// multi-line `<div>` blocks from being re-wrapped in `<p>`; tag-based
// classification gets the same result for the common cases and
// removes a footgun where a single-line `<div>x</div>` was getting
// paragraph-wrapped.
func (eval *Evaluate) dispatchRawHTML(node ast.JSXElement) error {
	attrs, err := eval.renderRawHTMLAttrs(node.Props)
	if err != nil {
		return fmt.Errorf("rendering attrs for <%s>: %w", node.Name, err)
	}

	if len(node.Children) == 0 {
		eval.Result = booklit.Append(eval.Result, booklit.RawElement{
			Tag:   node.Name,
			Attrs: attrs,
		})
		return nil
	}

	children, err := eval.evalArg(ast.Sequence(node.Children))
	if err != nil {
		return err
	}

	eval.Result = booklit.Append(eval.Result, booklit.RawElement{
		Tag:     node.Name,
		Attrs:   attrs,
		Content: children,
	})
	return nil
}

// renderRawHTMLAttrs concatenates ` name="value"` fragments for each
// prop. String values pass through verbatim (the source already wrote
// them as escaped HTML); expression values evaluate against the Dang
// env, stringify, and HTML-escape.
//
// Attrs are emitted in alphabetical order. Go's map iteration is
// randomized, and rendered HTML is observed by golden-file tests and
// the docs build, so a deterministic order matters. Alphabetical is
// the canonical pick — no need to track authored order on the AST.
func (eval *Evaluate) renderRawHTMLAttrs(props map[string]ast.Node) (string, error) {
	if len(props) == 0 {
		return "", nil
	}
	names := make([]string, 0, len(props))
	for name := range props {
		names = append(names, name)
	}
	sort.Strings(names)
	var out string
	for _, name := range names {
		val := props[name]
		switch v := val.(type) {
		case ast.String:
			out += " " + name + `="` + htmlEscapeAttr(string(v)) + `"`
		case ast.JSXExpression:
			if eval.Dang == nil {
				return "", fmt.Errorf("attr %q={%s}: no Dang evaluator", name, v.Raw)
			}
			dv, err := eval.Dang.Eval(v.Raw)
			if err != nil {
				return "", fmt.Errorf("attr %q={%s}: %w", name, v.Raw, err)
			}
			content, err := eval.Dang.ContentFromValue(dv, eval.Section)
			if err != nil {
				return "", fmt.Errorf("attr %q={%s}: %w", name, v.Raw, err)
			}
			s := ""
			if content != nil {
				s = content.String()
			}
			out += " " + name + `="` + htmlEscapeAttr(s) + `"`
		default:
			content, err := eval.evalArg(val)
			if err != nil {
				return "", fmt.Errorf("attr %q: %w", name, err)
			}
			s := ""
			if content != nil {
				s = content.String()
			}
			out += " " + name + `="` + htmlEscapeAttr(s) + `"`
		}
	}
	return out, nil
}

// htmlEscapeAttr escapes a string for safe inclusion in a double-quoted
// HTML attribute value.
func htmlEscapeAttr(s string) string {
	var b []byte
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '&':
			b = append(b, "&amp;"...)
		case '<':
			b = append(b, "&lt;"...)
		case '>':
			b = append(b, "&gt;"...)
		case '"':
			b = append(b, "&quot;"...)
		default:
			b = append(b, s[i])
		}
	}
	return string(b)
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
		Section:   eval.Section,
		Dang:      eval.Dang,
		Templates: eval.Templates,
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

func (eval Evaluate) evalArg(node ast.Node) (booklit.Content, error) {
	subEval := &Evaluate{
		Section:   eval.Section,
		Dang:      eval.Dang,
		Templates: eval.Templates,
	}

	err := node.Visit(subEval)
	if err != nil {
		return nil, err
	}

	return subEval.Result, nil
}
