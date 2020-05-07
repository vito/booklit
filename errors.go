package booklit

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/segmentio/textio"
	"github.com/vito/booklit/ast"
)

type PrettyError interface {
	PrettyPrint(io.Writer)
}

type UnknownReferenceError struct {
	TagName string

	ErrorLocation
}

func (err UnknownReferenceError) Error() string {
	return fmt.Sprintf("unknown tag '%s'", err.TagName)
}

func (err UnknownReferenceError) PrettyPrint(out io.Writer) {
	fmt.Fprintf(out, err.Annotate("reference points to unknown tag '%s':\n\n", err.TagName))
	err.AnnotateLocation(out)
}

type AmbiguousReferenceError struct {
	TagName          string
	DefinedLocations []ErrorLocation

	ErrorLocation
}

func (err AmbiguousReferenceError) Error() string {
	return fmt.Sprintf(
		"ambiguous target for tag '%s'",
		err.TagName,
	)
}

func (err AmbiguousReferenceError) PrettyPrint(out io.Writer) {
	fmt.Fprintf(out, err.Annotate("%s:\n\n", err))
	err.AnnotateLocation(out)
	fmt.Fprintf(out, "the same tag was defined in the following locations:\n\n")

	for _, loc := range err.DefinedLocations {
		fmt.Fprintf(out, "- %s:\n", loc.FilePath)
		loc.AnnotateLocation(textio.NewPrefixWriter(out, "  "))
	}
}

type UndefinedFunctionError struct {
	Function string

	ErrorLocation
}

func (err UndefinedFunctionError) Error() string {
	return fmt.Sprintf(
		"undefined function \\%s",
		err.Function,
	)
}

func (err UndefinedFunctionError) PrettyPrint(out io.Writer) {
	fmt.Fprintf(out, err.Annotate("%s:\n\n", err))
	err.AnnotateLocation(out)
}

type FailedFunctionError struct {
	Function string
	Err      error

	ErrorLocation
}

func (err FailedFunctionError) Error() string {
	return fmt.Sprintf(
		"function \\%s returned an error: %s",
		err.Function,
		err.Err,
	)
}

func (err FailedFunctionError) PrettyPrint(out io.Writer) {
	fmt.Fprintf(out, err.Annotate("function \\%s returned an error:\n\n", err.Function))
	err.AnnotateLocation(out)
	fmt.Fprintf(out, "error: %s\n", err.Err)
}

type ErrorLocation struct {
	FilePath     string
	NodeLocation ast.Location
	Length       int
}

func (loc ErrorLocation) Annotate(msg string, args ...interface{}) string {
	if loc.NodeLocation.Line == 0 {
		return fmt.Sprintf("%s: %s", loc.FilePath, fmt.Sprintf(msg, args...))
	} else {
		return fmt.Sprintf("%s:%d: %s", loc.FilePath, loc.NodeLocation.Line, fmt.Sprintf(msg, args...))
	}
}

func (loc ErrorLocation) AnnotateLocation(out io.Writer) error {
	if loc.NodeLocation.Line == 0 {
		// location unavailable
		return nil
	}

	file, err := os.Open(loc.FilePath)
	if err != nil {
		return err
	}

	buf := bufio.NewReader(file)

	for i := 0; i < loc.NodeLocation.Line-1; i++ {
		_, _, err := buf.ReadLine()
		if err != nil {
			return err
		}
	}

	lineInQuestion, _, err := buf.ReadLine()
	if err != nil {
		return err
	}

	prefix := fmt.Sprintf("% 4d| ", loc.NodeLocation.Line)

	_, err = fmt.Fprintf(out, "%s%s\n", prefix, lineInQuestion)
	if err != nil {
		return err
	}

	pad := strings.Repeat(" ", len(prefix)+loc.NodeLocation.Col-1)
	_, err = fmt.Fprintf(out, "%s\x1b[31m%s\x1b[0m\n", pad, strings.Repeat("^", loc.Length))
	if err != nil {
		return err
	}

	return nil
}
