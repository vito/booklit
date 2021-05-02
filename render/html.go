package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vito/booklit"
	"github.com/vito/booklit/render/html"
)

var initHTMLTmpl *template.Template

// HTMLFuncs is the set of functions available to all templates.
//
// See https://booklit.page/html-renderer.html#template-functions for more
// information.
var HTMLFuncs = template.FuncMap{
	"render": func(booklit.Content) (template.HTML, error) {
		return "", errors.New("render stubbed out")
	},

	"url": func(tag booklit.Tag) string {
		return sectionURL("html", tag.Section, tag.Anchor)
	},

	"stripAux": booklit.StripAux,

	"rawHTML": func(con booklit.Content) template.HTML {
		return template.HTML(con.String())
	},

	"rawURL": func(con booklit.Content) template.URL {
		return template.URL(con.String())
	},

	"walkContext": func(current *booklit.Section, section *booklit.Section) WalkContext {
		return WalkContext{
			Current: current,
			Section: section,
		}
	},

	"headerDepth": func(con *booklit.Section) int {
		depth := con.PageDepth() + 1
		if depth > 6 {
			depth = 6
		}

		return depth
	},
}

func init() {
	initHTMLTmpl = template.New("engine").Funcs(HTMLFuncs)

	for _, asset := range html.AssetNames() {
		info, err := html.AssetInfo(asset)
		if err != nil {
			panic(err)
		}

		content := strings.TrimRight(string(html.MustAsset(asset)), "\n")

		_, err = initHTMLTmpl.New(filepath.Base(info.Name())).Parse(content)
		if err != nil {
			panic(err)
		}
	}
}

// HTMLEngine renders sections as HTML using Go's html/template
// system.
type HTMLEngine struct {
	tmpl         *template.Template
	tmplModTimes map[string]time.Time

	template *template.Template
	data     interface{}
}

// NewHTMLEngine constructs a new HTMLEngine with the basic set of HTML
// templates bundled with Booklit.
func NewHTMLEngine() *HTMLEngine {
	engine := &HTMLEngine{
		tmplModTimes: map[string]time.Time{},
	}

	engine.resetTmpl()

	return engine
}

func (engine *HTMLEngine) resetTmpl() {
	engine.tmpl = template.Must(initHTMLTmpl.Clone())
	engine.tmpl.Funcs(template.FuncMap{
		"render": engine.subRender,
	})
}

// LoadTemplates loads all *.tmpl files in the specified directory.
func (engine *HTMLEngine) LoadTemplates(templatesDir string) error {
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

// FileExtension returns "html".
func (engine *HTMLEngine) FileExtension() string {
	return "html"
}

// URL returns the HTML file name with an anchor if present.
func (engine *HTMLEngine) URL(tag booklit.Tag) string {
	return sectionURL(engine.FileExtension(), tag.Section, tag.Anchor)
}

// RenderSection renders the section to the writer using page.tmpl.
//
// If the section has Style set and a template named (Style)-page.tmpl exists
// it will be used instead.
func (engine *HTMLEngine) RenderSection(out io.Writer, con *booklit.Section) error {
	engine.data = con

	try := []string{}

	if con.Style != "" {
		try = append(try, con.Style+"-page")
	}

	try = append(try, "page")

	var err error
	for _, tmpl := range try {
		err = engine.setTmpl(tmpl)
		if err == nil {
			break
		}
	}
	if err != nil {
		return err
	}

	return engine.render(out)
}

// VisitString renders con using string.tmpl.
func (engine *HTMLEngine) VisitString(con booklit.String) error {
	engine.data = con
	return engine.setTmpl("string")
}

// VisitReference renders con using reference.tmpl.
func (engine *HTMLEngine) VisitReference(con *booklit.Reference) error {
	engine.data = con
	return engine.setTmpl("reference")
}

// VisitSection renders con using section.tmpl.
//
// If the section has Style set and a template named (Style).tmpl exists it
// will be used instead.
func (engine *HTMLEngine) VisitSection(con *booklit.Section) error {
	tmpl := "section"
	if con.Style != "" {
		tmpl = con.Style
	}

	engine.data = con
	return engine.setTmpl(tmpl)
}

// VisitSequence renders con using sequence.tmpl.
func (engine *HTMLEngine) VisitSequence(con booklit.Sequence) error {
	engine.data = con
	return engine.setTmpl("sequence")
}

// VisitParagraph renders con using paragraph.tmpl.
func (engine *HTMLEngine) VisitParagraph(con booklit.Paragraph) error {
	engine.data = con
	return engine.setTmpl("paragraph")
}

// VisitPreformatted renders con using preformatted.tmpl.
func (engine *HTMLEngine) VisitPreformatted(con booklit.Preformatted) error {
	engine.data = con
	return engine.setTmpl("preformatted")
}

// VisitTableOfContents renders con using toc.tmpl.
func (engine *HTMLEngine) VisitTableOfContents(con booklit.TableOfContents) error {
	engine.data = con.Section
	return engine.setTmpl("toc")
}

// VisitStyled renders con using (Style).tmpl.
func (engine *HTMLEngine) VisitStyled(con booklit.Styled) error {
	engine.data = con
	return engine.setTmpl(string(con.Style))
}

// VisitTarget renders con using target.tmpl.
func (engine *HTMLEngine) VisitTarget(con booklit.Target) error {
	engine.data = con
	return engine.setTmpl("target")
}

// VisitImage renders con using image.tmpl.
func (engine *HTMLEngine) VisitImage(con booklit.Image) error {
	engine.data = con
	return engine.setTmpl("image")
}

// VisitList renders con using list.tmpl.
func (engine *HTMLEngine) VisitList(con booklit.List) error {
	engine.data = con
	return engine.setTmpl("list")
}

// VisitLink renders con using link.tmpl.
func (engine *HTMLEngine) VisitLink(con booklit.Link) error {
	engine.data = con
	return engine.setTmpl("link")
}

// VisitTable renders con using table.tmpl.
func (engine *HTMLEngine) VisitTable(con booklit.Table) error {
	engine.data = con
	return engine.setTmpl("table")
}

// VisitDefinitions renders con using definitions.tmpl.
func (engine *HTMLEngine) VisitDefinitions(con booklit.Definitions) error {
	engine.data = con
	return engine.setTmpl("definitions")
}

func (engine *HTMLEngine) setTmpl(name string) error {
	tmpl := engine.tmpl.Lookup(name + ".tmpl")

	if tmpl == nil {
		return fmt.Errorf("template '%s.tmpl' not found", name)
	}

	engine.template = tmpl

	return nil
}

func (engine *HTMLEngine) render(out io.Writer) error {
	if engine.template == nil {
		return fmt.Errorf("unknown template for '%s' (%T)", engine.data, engine.data)
	}

	return engine.template.Execute(out, engine.data)
}

func (engine *HTMLEngine) subRender(content booklit.Content) (template.HTML, error) {
	buf := new(bytes.Buffer)

	subEngine := NewHTMLEngine()
	subEngine.tmpl = engine.tmpl

	err := content.Visit(subEngine)
	if err != nil {
		return "", err
	}

	err = subEngine.render(buf)
	if err != nil {
		return "", err
	}

	return template.HTML(buf.String()), nil
}
