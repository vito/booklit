package dangeval

import (
	"github.com/vito/booklit"
	"github.com/vito/dang/pkg/dang"
	"github.com/vito/dang/pkg/hm"
)

// ContentValue is a Dang value that carries Booklit content unchanged.
// Templates use it to bind `children` in Dang scope: `{children}` looks
// the value up, the bridge unwraps it, and the original content is
// emitted at the interpolation site without going through String
// stringification (which would lose structure like nested styles).
//
// The declared Dang type is `String!` so Dang's inferrer accepts
// `{children}` as a valid expression. Operations that would require
// real string semantics (e.g. `children + "x"`) will fail at runtime
// because the Go type isn't dang.StringValue — that's acceptable for v1
// (templates use `{children}` only for emission).
type ContentValue struct {
	Content booklit.Content
}

var _ dang.Value = ContentValue{}

// Type returns String! so the inferrer can resolve `{children}` against
// the template's type env.
func (c ContentValue) Type() hm.Type {
	return hm.NonNullType{Type: dang.StringType}
}

// String returns the content's string form for debugging/Dang's stringly
// operations. The bridge picks up the Content field directly, so this
// is only hit by fallback paths.
func (c ContentValue) String() string {
	if c.Content == nil {
		return ""
	}
	return c.Content.String()
}

// LookupValue returns the raw Dang value bound to name, or (nil, false,
// nil) if it isn't bound. Templates use this to fetch the `children`
// content from inside the `<Children/>` built-in.
func (e *Evaluator) LookupValue(name string) (dang.Value, bool, error) {
	val, ok, err := e.evalEnv.Lookup(e.ctx, name)
	if err != nil {
		return nil, false, err
	}
	return val, ok, nil
}
