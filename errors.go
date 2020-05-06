package booklit

import (
	"fmt"
	"strings"

	"github.com/vito/booklit/ast"
)

type ErrorLocation struct {
	FilePath     string
	NodeLocation ast.Location
}

func (loc ErrorLocation) Annotate(msg string, args ...interface{}) string {
	if loc.NodeLocation.Line == 0 {
		return fmt.Sprintf("%s: %s", loc.FilePath, fmt.Sprintf(msg, args...))
	} else {
		return fmt.Sprintf("%s:%d: %s", loc.FilePath, loc.NodeLocation.Line, fmt.Sprintf(msg, args...))
	}
}

type UnknownReferenceError struct {
	TagName string

	ErrorLocation
}

func (err UnknownReferenceError) Error() string {
	return err.Annotate("unknown tag '%s'", err.TagName)
}

type AmbiguousReferenceError struct {
	TagName          string
	DefinedLocations []string

	ErrorLocation
}

func (err AmbiguousReferenceError) Error() string {
	return err.Annotate(
		"ambiguous target for tag '%s'\n\ntag '%s' is defined in multiple locations:\n\n - %s",
		err.TagName,
		err.TagName,
		strings.Join(err.DefinedLocations, "\n - "),
	)
}

type UndefinedFunctionError struct {
	Function string

	ErrorLocation
}

func (err UndefinedFunctionError) Error() string {
	return err.Annotate(
		"undefined function \\%s",
		err.Function,
	)
}

type FailedFunctionError struct {
	Function string
	Err      error

	ErrorLocation
}

func (err FailedFunctionError) Error() string {
	return err.Annotate(
		"function \\%s returned an error: %s",
		err.Function,
		err.Err,
	)
}
