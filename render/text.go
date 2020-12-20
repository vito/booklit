package render

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/vito/booklit"
	"github.com/vito/booklit/render/text"
)

var initTextTmpl *template.Template

// TextFuncs is the set of functions available to all templates.
var TextFuncs = template.FuncMap{
	"render": func(booklit.Content) (string, error) {
		return "", errors.New("render stubbed out")
	},

	"url": func(ext string, tag booklit.Tag) (string, error) {
		return "", errors.New("url stubbed out")
	},

	"htmlURL": func(tag booklit.Tag) string {
		return sectionURL("html", tag.Section, tag.Anchor)
	},

	"stripAux": booklit.StripAux,

	"joinLines": func(prefix string, str string) string {
		return strings.Join(strings.Split(str, "\n"), "\n"+prefix)
	},
}

func init() {
	initTextTmpl = template.New("engine").Funcs(TextFuncs)

	for _, asset := range text.AssetNames() {
		info, err := text.AssetInfo(asset)
		if err != nil {
			panic(err)
		}

		content := strings.TrimRight(string(text.MustAsset(asset)), "\n")

		_, err = initTextTmpl.New(filepath.Base(info.Name())).Parse(content)
		if err != nil {
			panic(err)
		}
	}
}

// TextEngine renders sections as plaintext using Go's text/template system.
//
// Text templates may be provided to generate e.g. Markdown or other plaintext
// formats.
type TextEngine struct {
	fileExtension string

	tmpl         *template.Template
	tmplModTimes map[string]time.Time

	template *template.Template
	data     interface{}
}

// NewTextEngine constructs a new TextEngine with the basic set of text
// templates bundled with Booklit.
//
// A file extension must be provided, e.g. "md" for Markdown.
func NewTextEngine(fileExtension string) *TextEngine {
	engine := &TextEngine{
		fileExtension: fileExtension,

		tmplModTimes: map[string]time.Time{},
	}

	engine.resetTmpl()

	return engine
}

func (engine *TextEngine) resetTmpl() {
	engine.tmpl = template.Must(initTextTmpl.Clone())
	engine.tmpl.Funcs(template.FuncMap{
		"render": engine.subRender,
		"url": func(tag booklit.Tag) string {
			return sectionURL(engine.FileExtension(), tag.Section, tag.Anchor)
		},
	})
}

// LoadTemplates loads all *.tmpl files in the specified directory.
func (engine *TextEngine) LoadTemplates(templatesDir string) error {
	templates, err := filepath.Glob(filepath.Join(templatesDir, "*.tmpl"))
	if err != nil {
		return err
	}

	var shouldReload bool
	for _, path := range templates {
		info, err := os.Stat(path)
		if err != nil {
			return err
		}

		modTime := info.ModTime()

		lastModTime, found := engine.tmplModTimes[path]
		if !found || modTime.After(lastModTime) {
			shouldReload = true
		}

		engine.tmplModTimes[path] = modTime
	}

	if !shouldReload {
		return nil
	}

	engine.resetTmpl()

	for _, path := range templates {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		trimmed := strings.TrimRight(string(content), "\n")

		_, err = engine.tmpl.New(filepath.Base(path)).Parse(trimmed)
		if err != nil {
			return err
		}
	}

	return nil
}

// FileExtension returns the configured file extension.
func (engine *TextEngine) FileExtension() string {
	return engine.fileExtension
}

// URL returns the file name using te configured file extension, with an anchor
// if present.
func (engine *TextEngine) URL(tag booklit.Tag) string {
	return sectionURL(engine.FileExtension(), tag.Section, tag.Anchor)
}

// RenderSection renders the section to the writer using page.tmpl.
//
// If the section has Style set and a template named (Style)-page.tmpl exists
// it will be used instead.
func (engine *TextEngine) RenderSection(out io.Writer, con *booklit.Section) error {
	tmpl := "page"
	if con.Style != "" {
		tmpl = con.Style + "-page"
	}

	engine.data = con

	err := engine.setTmpl(tmpl)
	if err != nil {
		return err
	}

	return engine.render(out)
}

// VisitString renders con using string.tmpl.
func (engine *TextEngine) VisitString(con booklit.String) error {
	engine.data = con
	return engine.setTmpl("string")
}

// VisitReference renders con using reference.tmpl.
func (engine *TextEngine) VisitReference(con *booklit.Reference) error {
	engine.data = con
	return engine.setTmpl("reference")
}

// VisitSection renders con using section.tmpl.
//
// If the section has Style set and a template named (Style).tmpl exists it
// will be used instead.
func (engine *TextEngine) VisitSection(con *booklit.Section) error {
	tmpl := "section"
	if con.Style != "" {
		tmpl = con.Style
	}

	engine.data = con
	return engine.setTmpl(tmpl)
}

// VisitSequence renders con using sequence.tmpl.
func (engine *TextEngine) VisitSequence(con booklit.Sequence) error {
	engine.data = con
	return engine.setTmpl("sequence")
}

// VisitParagraph renders con using paragraph.tmpl.
func (engine *TextEngine) VisitParagraph(con booklit.Paragraph) error {
	engine.data = con
	return engine.setTmpl("paragraph")
}

// VisitPreformatted renders con using preformatted.tmpl.
func (engine *TextEngine) VisitPreformatted(con booklit.Preformatted) error {
	engine.data = con
	return engine.setTmpl("preformatted")
}

// VisitTableOfContents renders con using toc.tmpl.
func (engine *TextEngine) VisitTableOfContents(con booklit.TableOfContents) error {
	engine.data = con.Section
	return engine.setTmpl("toc")
}

// VisitStyled renders con using (Style).tmpl.
func (engine *TextEngine) VisitStyled(con booklit.Styled) error {
	engine.data = con
	return engine.setTmpl(string(con.Style))
}

// VisitTarget renders con using target.tmpl.
func (engine *TextEngine) VisitTarget(con booklit.Target) error {
	engine.data = con
	return engine.setTmpl("target")
}

// VisitImage renders con using image.tmpl.
func (engine *TextEngine) VisitImage(con booklit.Image) error {
	engine.data = con
	return engine.setTmpl("image")
}

// VisitList renders con using list.tmpl.
func (engine *TextEngine) VisitList(con booklit.List) error {
	engine.data = con
	return engine.setTmpl("list")
}

// VisitLink renders con using link.tmpl.
func (engine *TextEngine) VisitLink(con booklit.Link) error {
	engine.data = con
	return engine.setTmpl("link")
}

// VisitTable renders con using table.tmpl.
func (engine *TextEngine) VisitTable(con booklit.Table) error {
	engine.data = con
	return engine.setTmpl("table")
}

// VisitDefinitions renders con using definitions.tmpl.
func (engine *TextEngine) VisitDefinitions(con booklit.Definitions) error {
	engine.data = con
	return engine.setTmpl("definitions")
}

func (engine *TextEngine) setTmpl(name string) error {
	tmpl := engine.tmpl.Lookup(name + ".tmpl")

	if tmpl == nil {
		return fmt.Errorf("template '%s.tmpl' not found", name)
	}

	engine.template = tmpl

	return nil
}

func (engine *TextEngine) render(out io.Writer) error {
	if engine.template == nil {
		return fmt.Errorf("unknown template for '%s' (%T)", engine.data, engine.data)
	}

	return engine.template.Execute(out, engine.data)
}

func (engine *TextEngine) subRender(content booklit.Content) (string, error) {
	buf := new(bytes.Buffer)

	subEngine := NewTextEngine(engine.fileExtension)
	subEngine.tmpl = engine.tmpl

	err := content.Visit(subEngine)
	if err != nil {
		return "", err
	}

	err = subEngine.render(buf)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
