package booklit

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/segmentio/textio"
	"github.com/vito/booklit/ast"
	"github.com/vito/booklit/errhtml"
)

var errorTmpl *template.Template

func init() {
	errorTmpl = template.New("errors").Funcs(template.FuncMap{
		"error": func(err error) (template.HTML, error) {
			buf := new(bytes.Buffer)
			if prettyErr, ok := err.(PrettyError); ok {
				renderErr := prettyErr.PrettyHTML(buf)
				if renderErr != nil {
					return "", renderErr
				}

				return template.HTML(buf.String()), nil
			}

			return template.HTML(
				`<pre class="raw-error">` +
					template.HTMLEscapeString(err.Error()) +
					`</pre>`,
			), nil
		},

		"annotate": func(loc ErrorLocation) (template.HTML, error) {
			buf := new(bytes.Buffer)
			err := loc.AnnotatedHTML(buf)
			if err != nil {
				return "", err
			}

			return template.HTML(buf.String()), nil
		},
	})

	for _, asset := range errhtml.AssetNames() {
		info, err := errhtml.AssetInfo(asset)
		if err != nil {
			panic(err)
		}

		content := strings.TrimRight(string(errhtml.MustAsset(asset)), "\n")

		_, err = errorTmpl.New(filepath.Base(info.Name())).Parse(content)
		if err != nil {
			panic(err)
		}
	}
}

// ErrorResponse writes the error response page.
//
// If err implements PrettyHTML it can render its own HTML template with
// additional troubleshooting.
func ErrorResponse(w http.ResponseWriter, err error) {
	renderErr := errorTmpl.Lookup("page.tmpl").Execute(w, err)
	if renderErr != nil {
		fmt.Fprintf(w, "failed to render error page: %s", renderErr)
	}
}

// PrettyError is an interface for providing friendly error messages.
type PrettyError interface {
	error

	// PrettyPrint is called by the booklit CLI to print an error message to
	// stderr.
	PrettyPrint(io.Writer)

	// PrettyHTML is called by the error page template to render HTML within the
	// error page.
	PrettyHTML(io.Writer) error
}

// ParseError is returned when a Booklit document fails to parse.
type ParseError struct {
	Err error

	ErrorLocation
}

// Error returns a 'parse error' error message.
func (err ParseError) Error() string {
	return fmt.Sprintf("parse error: %s", err.Err)
}

// PrettyPrint prints the error message followed by a snippet of the source
// location where the error occurred.
func (err ParseError) PrettyPrint(out io.Writer) {
	fmt.Fprintln(out, err.Annotate("%s", err))
	fmt.Fprintln(out)
	err.AnnotateLocation(out)
}

// PrettyHTML renders an HTML template containing the error message followed by
// a snippet of the source location where the error occurred.
func (err ParseError) PrettyHTML(out io.Writer) error {
	return errorTmpl.Lookup("parse-error.tmpl").Execute(out, err)
}

// UnknownTagError is returned when a reference is made to an unknown tag.
type UnknownTagError struct {
	TagName string

	SimilarTags []Tag

	ErrorLocation
}

// Error returns an 'unknown tag' error message.
func (err UnknownTagError) Error() string {
	return fmt.Sprintf("unknown tag '%s'", err.TagName)
}

// PrettyPrint prints the error message, a snippet of the source code where the
// error occurred, and suggests similar tags.
func (err UnknownTagError) PrettyPrint(out io.Writer) {
	fmt.Fprintln(out, err.Annotate("%s", err))
	fmt.Fprintln(out)
	err.AnnotateLocation(out)

	if len(err.SimilarTags) == 0 {
		fmt.Fprintf(out, "I couldn't find any similar tags. :(\n")
	} else {
		fmt.Fprintf(out, "These tags seem similar:\n\n")

		for _, tag := range err.SimilarTags {
			fmt.Fprintf(out, "- %s\n", tag.Name)
		}

		fmt.Fprintf(out, "\nDid you mean one of these?\n")
	}
}

// PrettyHTML renders an HTML template containing the error message, a snippet
// of the source code where the error occurred, and suggests similar tags.
func (err UnknownTagError) PrettyHTML(out io.Writer) error {
	return errorTmpl.Lookup("unknown-tag.tmpl").Execute(out, err)
}

// AmbiguousReferenceError is returned when a referenced tag is defined in
// multiple places.
type AmbiguousReferenceError struct {
	TagName          string
	DefinedLocations []ErrorLocation

	ErrorLocation
}

// Error returns an 'ambiguous target for tag' error message.
func (err AmbiguousReferenceError) Error() string {
	return fmt.Sprintf(
		"ambiguous target for tag '%s'",
		err.TagName,
	)
}

// PrettyPrint prints the error message, a snippet of the source code where the
// error occurred, and snippets for the definition location of each tag that
// was found.
func (err AmbiguousReferenceError) PrettyPrint(out io.Writer) {
	fmt.Fprintln(out, err.Annotate("%s", err))
	fmt.Fprintln(out)
	err.AnnotateLocation(out)

	fmt.Fprintf(out, "The same tag was defined in the following locations:\n\n")

	for _, loc := range err.DefinedLocations {
		fmt.Fprintf(out, "- %s:\n", loc.FilePath)
		loc.AnnotateLocation(textio.NewPrefixWriter(out, "  "))
	}

	fmt.Fprintf(out, "Tags must be unique so I know where to link to!\n")
}

// PrettyHTML renders a HTML template containing the error message, a snippet
// of the source code where the error occurred, and snippets for the definition
// location of each tag that was found.
func (err AmbiguousReferenceError) PrettyHTML(out io.Writer) error {
	return errorTmpl.Lookup("ambiguous-reference.tmpl").Execute(out, err)
}

// UndefinedFunctionError is returned when a Booklit document tries to call a
// function that is not defined by any plugin.
type UndefinedFunctionError struct {
	Function string

	ErrorLocation
}

// Error returns an 'undefined function' error message.
func (err UndefinedFunctionError) Error() string {
	return fmt.Sprintf(
		"undefined function \\%s",
		err.Function,
	)
}

// PrettyPrint prints the error message and a snippet of the source code where
// the error occurred.
func (err UndefinedFunctionError) PrettyPrint(out io.Writer) {
	fmt.Fprintln(out, err.Annotate("%s", err))
	fmt.Fprintln(out)
	err.AnnotateLocation(out)
}

// PrettyHTML renders an HTML template containing the error message and a
// snippet of the source code where the error occurred.
func (err UndefinedFunctionError) PrettyHTML(out io.Writer) error {
	return errorTmpl.Lookup("undefined-function.tmpl").Execute(out, err)
}

// FailedFunctionError is returned when a plugin function called by a Booklit
// document returns an error.
type FailedFunctionError struct {
	Function string
	Err      error

	ErrorLocation
}

// Error returns a 'function \... returned an error' message specifying the
// function name and the error it returned.
func (err FailedFunctionError) Error() string {
	return fmt.Sprintf(
		"function \\%s returned an error: %s",
		err.Function,
		err.Err,
	)
}

// PrettyPrint prints the error message and a snippet of the source code where
// the error occurred.
//
// If the error returned by the function is a PrettyError, PrettyPrint is
// called and its output is indented.
//
// Otherwise, the error is printed normally.
func (err FailedFunctionError) PrettyPrint(out io.Writer) {
	fmt.Fprintln(out, err.Annotate("function \\%s returned an error", err.Function))
	fmt.Fprintln(out)
	err.AnnotateLocation(out)

	if prettyErr, ok := err.Err.(PrettyError); ok {
		prettyErr.PrettyPrint(textio.NewPrefixWriter(out, "  "))
	} else {
		fmt.Fprintf(out, "\x1b[33m%s\x1b[0m\n", err.Err)
	}
}

// PrettyHTML renders an HTML template containing the error message followed by
// a snippet of the source location where the error occurred.
//
// If the error returned by the function is a PrettyError, PrettyHTML will be
// called within the template to embed the error recursively.
func (err FailedFunctionError) PrettyHTML(out io.Writer) error {
	return errorTmpl.Lookup("function-error.tmpl").Execute(out, err)
}

// TitleTwiceError is returned when a section tries to set \title twice.
type TitleTwiceError struct {
	TitleLocation ErrorLocation

	ErrorLocation
}

// Error returns a 'cannot set title twice' message.
func (err TitleTwiceError) Error() string {
	return "cannot set title twice"
}

// PrettyPrint prints the error message and a snippet of the source code where
// the error occurred.
//
// If the error returned by the function is a PrettyError, PrettyPrint is
// called and its output is indented.
//
// Otherwise, the error is printed normally.
func (err TitleTwiceError) PrettyPrint(out io.Writer) {
	fmt.Fprintln(out, err.Annotate("%s", err))
	fmt.Fprintln(out)
	err.AnnotateLocation(out)

	fmt.Fprintf(out, "The section's title was first set here:\n\n")
	err.TitleLocation.AnnotateLocation(out)

	fmt.Fprintln(out, "Maybe the second \\title should be in a \\section{...}?")
}

// PrettyHTML renders an HTML template containing the error message followed by
// a snippet of the source location where the error occurred.
//
// If the error returned by the function is a PrettyError, PrettyHTML will be
// called within the template to embed the error recursively.
func (err TitleTwiceError) PrettyHTML(out io.Writer) error {
	return errorTmpl.Lookup("title-twice-error.tmpl").Execute(out, err)
}

// ErrorLocation is the source location in a Booklit document where an error
// occurred.
type ErrorLocation struct {
	FilePath     string
	NodeLocation ast.Location
	Length       int
}

// Annotate prepends the source location to the given message.
func (loc ErrorLocation) Annotate(msg string, args ...interface{}) string {
	if loc.NodeLocation.Line == 0 {
		return fmt.Sprintf("%s: %s", loc.FilePath, fmt.Sprintf(msg, args...))
	}

	return fmt.Sprintf("%s:%d: %s", loc.FilePath, loc.NodeLocation.Line, fmt.Sprintf(msg, args...))
}

// AnnotateLocation writes a plaintext snippet of the location in the Booklit
// document.
func (loc ErrorLocation) AnnotateLocation(out io.Writer) error {
	if loc.NodeLocation.Line == 0 {
		// location unavailable
		return nil
	}

	line, err := loc.lineInQuestion()
	if err != nil {
		return err
	}

	prefix := fmt.Sprintf("% 4d| ", loc.NodeLocation.Line)

	_, err = fmt.Fprintf(out, "%s%s\n", prefix, line)
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

type annotationData struct {
	FilePath                  string
	EOF                       bool
	Lineno                    string
	Prefix, Annotated, Suffix string
}

// AnnotatedHTML renders a HTML snippet of the error location in the Booklit
// document.
func (loc ErrorLocation) AnnotatedHTML(out io.Writer) error {
	if loc.NodeLocation.Line == 0 {
		// location unavailable
		return nil
	}

	line, err := loc.lineInQuestion()
	if err != nil {
		return err
	}

	data := annotationData{
		FilePath: loc.FilePath,
		Lineno:   fmt.Sprintf("% 4d", loc.NodeLocation.Line),
	}

	if line == "" {
		data.EOF = true
	}

	offset := loc.NodeLocation.Col - 1
	if len(line) >= offset+loc.Length {
		data.Prefix = line[0:offset]
		data.Annotated = line[offset : offset+loc.Length]
		data.Suffix = line[offset+loc.Length:]
	}

	return errorTmpl.Lookup("annotated-line.tmpl").Execute(out, data)
}

func (loc ErrorLocation) lineInQuestion() (string, error) {
	file, err := os.Open(loc.FilePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	buf := bufio.NewReader(file)

	for i := 0; i < loc.NodeLocation.Line-1; i++ {
		_, _, err := buf.ReadLine()
		if err != nil {
			return "", err
		}
	}

	lineInQuestion, _, err := buf.ReadLine()
	if err != nil {
		if err == io.EOF {
			return "", nil
		}

		return "", err
	}

	return string(lineInQuestion), nil
}
