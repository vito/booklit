package ast

// Error is a parse error along with its location in the source.
type Error struct {
	Inner    error
	Location Location
}

// UnpackError converts an internal parse error to an *Error.
func UnpackError(err error) (*Error, bool) {
	if list, ok := err.(errList); ok {
		for _, e := range list {
			if perr, ok := e.(*parserError); ok {
				return &Error{
					Inner: perr.Inner,
					Location: Location{
						Line:   perr.pos.line,
						Col:    perr.pos.col,
						Offset: perr.pos.offset,
					},
				}, true
			}
		}
	}

	return nil, false
}
