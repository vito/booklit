package ast

func UnpackError(err error) (error, Location, bool) {
	if list, ok := err.(errList); ok {
		for _, e := range list {
			if perr, ok := e.(*parserError); ok {
				return perr.Inner, Location{
					Line:   perr.pos.line,
					Col:    perr.pos.col,
					Offset: perr.pos.offset,
				}, true
			}
		}
	}

	return nil, Location{}, false
}
