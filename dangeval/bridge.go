package dangeval

import (
	"fmt"
	"strconv"

	"github.com/vito/booklit"
	"github.com/vito/dang/pkg/dang"
)

// ToContent coerces a Dang value into booklit.Content. Scalars render
// as their string form; lists render as a Sequence of bridged elements;
// null renders as empty content. Anything richer (records, functions,
// modules) is an error for v1 — see jsx-dang.md Phase 3 open questions.
func ToContent(val dang.Value) (booklit.Content, error) {
	switch v := val.(type) {
	case nil:
		return booklit.Empty, nil
	case dang.NullValue:
		return booklit.Empty, nil
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
