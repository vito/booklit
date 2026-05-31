package dangeval

import (
	"fmt"
	"strconv"

	"github.com/vito/booklit"
	"github.com/vito/booklit/contentjson"
	"github.com/vito/dang/pkg/dang"
)

// ToContent coerces a Dang value into booklit.Content. Scalars render
// as their string form; lists render as a Sequence of bridged elements;
// null renders as empty content. A ContentValue unwraps to its carried
// Content unchanged — this is how templates pass `children` to a
// `{children}` interpolation without going through string-level
// flattening. Anything richer (records, functions, modules) is an error
// for v1 — see jsx-dang.md Phase 3 open questions.
func ToContent(val dang.Value) (booklit.Content, error) {
	switch v := val.(type) {
	case nil:
		return booklit.Empty, nil
	case dang.NullValue:
		return booklit.Empty, nil
	case ContentValue:
		if v.Content == nil {
			return booklit.Empty, nil
		}
		return v.Content, nil
	case dang.StringValue:
		return booklit.String(v.Val), nil
	case dang.IntValue:
		return booklit.String(strconv.Itoa(v.Val)), nil
	case dang.FloatValue:
		return booklit.String(strconv.FormatFloat(v.Val, 'g', -1, 64)), nil
	case dang.BoolValue:
		return booklit.String(strconv.FormatBool(v.Val)), nil
	case dang.ListValue:
		seq := make(booklit.Sequence, 0, len(v.Elements))
		for _, el := range v.Elements {
			c, err := ToContent(el)
			if err != nil {
				return nil, err
			}
			if c != nil {
				seq = append(seq, c)
			}
		}
		return seq, nil
	default:
		return nil, fmt.Errorf("cannot render Dang value of type %T as content: %s", val, val.String())
	}
}

// ContentFromValue coerces a Dang value into booklit.Content, with two
// capabilities ToContent lacks: it decodes Booklit content returned by a
// Dagger module, and it rehydrates Reference/Target nodes against sec.
//
// A Dagger function returning the `JSON` scalar (a dang.ScalarValue whose
// scalar type is JSON) or a `JSONValue!` object (a dang.GraphQLValue) is
// understood as a serialized content tree (see package contentjson) and
// decoded back into native booklit.Content. The JSON scalar arrives already
// materialized; the JSONValue object is forced by selecting its `contents`
// field, which costs one engine round-trip but lets the module compose results
// lazily. Everything else falls back to ToContent's primitive handling.
func (e *Evaluator) ContentFromValue(val dang.Value, sec *booklit.Section) (booklit.Content, error) {
	switch v := val.(type) {
	case dang.ScalarValue:
		if isJSONScalar(v) {
			return contentjson.Unmarshal([]byte(v.Val), sec)
		}
		return booklit.String(v.Val), nil
	case dang.GraphQLValue:
		if v.TypeName == "JSONValue" {
			contents, err := e.jsonValueContents(v)
			if err != nil {
				return nil, err
			}
			return contentjson.Unmarshal([]byte(contents), sec)
		}
		return nil, fmt.Errorf("cannot render Dang value of GraphQL type %s as content", v.TypeName)
	case dang.ListValue:
		seq := make(booklit.Sequence, 0, len(v.Elements))
		for _, el := range v.Elements {
			c, err := e.ContentFromValue(el, sec)
			if err != nil {
				return nil, err
			}
			if c != nil {
				seq = append(seq, c)
			}
		}
		return seq, nil
	default:
		return ToContent(val)
	}
}

// isJSONScalar reports whether v is a Dagger `JSON` scalar, i.e. a serialized
// content tree rather than ordinary text.
func isJSONScalar(v dang.ScalarValue) bool {
	mod, ok := v.ScalarType.(*dang.Module)
	return ok && mod.Named == "JSON"
}

// jsonValueContents forces a `JSONValue!` object by selecting its `contents`
// field and returns the serialized JSON. This is the one engine round-trip the
// lazy-object path costs over the scalar path.
func (e *Evaluator) jsonValueContents(v dang.GraphQLValue) (string, error) {
	sel, err := v.SelectField(e.ctx, "contents")
	if err != nil {
		return "", fmt.Errorf("selecting JSONValue.contents: %w", err)
	}
	fn, ok := sel.(dang.Callable)
	if !ok {
		return "", fmt.Errorf("JSONValue.contents is not callable (got %T)", sel)
	}
	res, err := fn.Call(e.ctx, e.evalEnv, map[string]dang.Value{})
	if err != nil {
		return "", fmt.Errorf("forcing JSONValue.contents: %w", err)
	}
	switch r := res.(type) {
	case dang.ScalarValue:
		return r.Val, nil
	case dang.StringValue:
		return r.Val, nil
	default:
		return "", fmt.Errorf("JSONValue.contents returned %T, expected a scalar", res)
	}
}
