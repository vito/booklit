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
)

var initTmpl *template.Template

type WalkContext struct {
	Current *booklit.Section
	Section *booklit.Section
}

func init() {
	initTmpl = template.New("engine").Funcs(template.FuncMap{
		"url": tagURLFunc,

		"stripAux": booklit.StripAux,

		"rawHTML": func(con booklit.Content) template.HTML {
			return template.HTML(con.String())
		},

		"render": func(booklit.Content) (template.HTML, error) {
			return "", errors.New("render stubbed out")
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
	})

	for _, asset := range AssetNames() {
		info, err := AssetInfo(asset)
		if err != nil {
			panic(err)
		}

		content := strings.TrimRight(string(MustAsset(asset)), "\n")

		_, err = initTmpl.New(filepath.Base(info.Name())).Parse(content)
		if err != nil {
			panic(err)
		}
	}
}

type HTMLRenderingEngine struct {
	tmpl         *template.Template
	tmplModTimes map[string]time.Time

	template *template.Template
	data     interface{}
}

func NewHTMLRenderingEngine() *HTMLRenderingEngine {
	engine := &HTMLRenderingEngine{
		tmplModTimes: map[string]time.Time{},
	}

	engine.resetTmpl()

	return engine
}

func (engine *HTMLRenderingEngine) resetTmpl() {
	engine.tmpl = template.Must(initTmpl.Clone())
	engine.tmpl.Funcs(template.FuncMap{
		"render": engine.subRender,
	})
}

func (engine *HTMLRenderingEngine) LoadTemplates(templatesDir string) error {
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

func (engine *HTMLRenderingEngine) FileExtension() string {
	return "html"
}

func (engine *HTMLRenderingEngine) URL(tag booklit.Tag) string {
	return tagURLFunc(tag)
}

func (engine *HTMLRenderingEngine) RenderSection(out io.Writer, con *booklit.Section) error {
	if con.Style != "" {
		engine.template = engine.tmpl.Lookup(con.Style + "-page.tmpl")
	} else {
		engine.template = engine.tmpl.Lookup("page.tmpl")
	}

	engine.data = con

	return engine.render(out)
}

func (engine *HTMLRenderingEngine) VisitString(con booklit.String) error {
	engine.template = engine.tmpl.Lookup("string.tmpl")
	engine.data = con
	return nil
}

func (engine *HTMLRenderingEngine) VisitReference(con *booklit.Reference) error {
	engine.template = engine.tmpl.Lookup("reference.tmpl")
	engine.data = con
	return nil
}

func (engine *HTMLRenderingEngine) VisitSection(con *booklit.Section) error {
	if con.Style != "" {
		engine.template = engine.tmpl.Lookup(con.Style + ".tmpl")
	} else {
		engine.template = engine.tmpl.Lookup("section.tmpl")
	}

	engine.data = con

	return nil
}

func (engine *HTMLRenderingEngine) VisitSequence(con booklit.Sequence) error {
	engine.template = engine.tmpl.Lookup("sequence.tmpl")
	engine.data = con
	return nil
}

func (engine *HTMLRenderingEngine) VisitParagraph(con booklit.Paragraph) error {
	engine.template = engine.tmpl.Lookup("paragraph.tmpl")
	engine.data = con
	return nil
}

func (engine *HTMLRenderingEngine) VisitPreformatted(con booklit.Preformatted) error {
	engine.template = engine.tmpl.Lookup("preformatted.tmpl")
	engine.data = con
	return nil
}

func (engine *HTMLRenderingEngine) VisitTableOfContents(con booklit.TableOfContents) error {
	engine.template = engine.tmpl.Lookup("toc.tmpl")
	engine.data = con.Section
	return nil
}

func (engine *HTMLRenderingEngine) VisitStyled(con booklit.Styled) error {
	engine.template = engine.tmpl.Lookup(string(con.Style) + ".tmpl")
	if engine.template == nil {
		return fmt.Errorf("no template defined for style: %s", con.Style)
	}

	engine.data = con

	return nil
}

func (engine *HTMLRenderingEngine) VisitTarget(con booklit.Target) error {
	engine.template = engine.tmpl.Lookup("target.tmpl")
	engine.data = con
	return nil
}

func (engine *HTMLRenderingEngine) VisitImage(con booklit.Image) error {
	engine.template = engine.tmpl.Lookup("image.tmpl")
	engine.data = con
	return nil
}

func (engine *HTMLRenderingEngine) VisitList(con booklit.List) error {
	engine.template = engine.tmpl.Lookup("list.tmpl")
	engine.data = con
	return nil
}

func (engine *HTMLRenderingEngine) VisitLink(con booklit.Link) error {
	engine.template = engine.tmpl.Lookup("link.tmpl")
	engine.data = con
	return nil
}

func (engine *HTMLRenderingEngine) VisitTable(con booklit.Table) error {
	engine.template = engine.tmpl.Lookup("table.tmpl")
	engine.data = con
	return nil
}

func (engine *HTMLRenderingEngine) VisitDefinitions(con booklit.Definitions) error {
	engine.template = engine.tmpl.Lookup("definitions.tmpl")
	engine.data = con
	return nil
}

func (engine *HTMLRenderingEngine) render(out io.Writer) error {
	if engine.template == nil {
		return fmt.Errorf("unknown template for '%s' (%T)", engine.data, engine.data)
	}

	return engine.template.Execute(out, engine.data)
}

func (engine *HTMLRenderingEngine) subRender(content booklit.Content) (template.HTML, error) {
	buf := new(bytes.Buffer)

	subEngine := NewHTMLRenderingEngine()
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

func tagURLFunc(tag booklit.Tag) string {
	return sectionURL(tag.Section, tag.Anchor)
}

func sectionURL(section *booklit.Section, anchor string) string {
	owner := sectionPageOwner(section)

	if owner != section {
		if anchor == "" {
			anchor = section.PrimaryTag.Name
		}

		return sectionURL(owner, anchor)
	}

	filename := section.PrimaryTag.Name + ".html"

	if anchor != "" {
		filename += "#" + anchor
	}

	return filename
}

func sectionPageOwner(section *booklit.Section) *booklit.Section {
	if section.Parent == nil {
		return section
	}

	if section.Parent.SplitSections {
		return section
	}

	return sectionPageOwner(section.Parent)
}
