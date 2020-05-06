package booklit

import (
	"fmt"
	"strings"

	"github.com/vito/booklit/ast"
)

func annotate(filePath string, loc ast.Location, msg string, args ...interface{}) string {
	if loc.Line == 0 {
		return fmt.Sprintf("%s: %s", filePath, fmt.Sprintf(msg, args...))
	} else {
		return fmt.Sprintf("%s:%d: %s", filePath, loc.Line, fmt.Sprintf(msg, args...))
	}
}

type UnknownReferenceError struct {
	TagName string

	FilePath string
	Location ast.Location
}

func (err UnknownReferenceError) Error() string {
	return annotate(err.FilePath, err.Location, "unknown tag '%s'", err.TagName)
}

type AmbiguousReferenceError struct {
	TagName          string
	DefinedLocations []string

	FilePath string
	Location ast.Location
}

func (err AmbiguousReferenceError) Error() string {
	return annotate(err.FilePath, err.Location,
		"ambiguous target for tag '%s'\n\ntag '%s' is defined in multiple locations:\n\n - %s",
		err.TagName,
		err.TagName,
		strings.Join(err.DefinedLocations, "\n - "),
	)
}

type UndefinedFunctionError struct {
	Function string
	FilePath string
	Location ast.Location
}

func (err UndefinedFunctionError) Error() string {
	return annotate(err.FilePath, err.Location,
		"undefined function \\%s",
		err.Function,
	)
}

type FailedFunctionError struct {
	Err error

	Function string
	FilePath string
	Location ast.Location
}

func (err FailedFunctionError) Error() string {
	return annotate(err.FilePath, err.Location,
		"function \\%s returned an error: %s",
		err.Function,
		err.Err,
	)
}
