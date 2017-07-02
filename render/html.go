package render

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"strings"

	"github.com/vito/booklit"
)

var initTmpl *template.Template

type WalkContext struct {
	Current *booklit.Section
	Section *booklit.Section
}

func init() {
	initTmpl = template.New("engine").Funcs(template.FuncMap{
		"render": renderFunc,
		"url":    tagURLFunc,
		"walkContext": func(current *booklit.Section, section *booklit.Section) WalkContext {
			return WalkContext{
				Current: current,
				Section: section,
			}
		},

		"headerDepth": func(con *booklit.Section) int {
			depth := 1
			for sec := con; sec.Parent != nil && !sec.Parent.SplitSections; sec = sec.Parent {
				depth++
			}

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
	tmpl *template.Template

	template *template.Template
	data     interface{}
}

func NewHTMLRenderingEngine() *HTMLRenderingEngine {
	return &HTMLRenderingEngine{
		tmpl: template.Must(initTmpl.Clone()),
	}
}

func (engine *HTMLRenderingEngine) LoadTemplates(templatesDir string) error {
	tmpl, err := engine.tmpl.ParseGlob(filepath.Join(templatesDir, "*.tmpl"))
	if err != nil {
		return err
	}

	engine.tmpl = tmpl

	return nil
}

func (engine *HTMLRenderingEngine) FileExtension() string {
	return "html"
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
	var pageTemplate *template.Template

	if con.Parent == nil || con.Parent.SplitSections {
		pageTemplate = engine.tmpl.Lookup("page.tmpl")
	}

	if pageTemplate == nil {
		pageTemplate = engine.tmpl.Lookup("section.tmpl")
	}

	engine.template = pageTemplate
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

func (engine *HTMLRenderingEngine) VisitTableOfContents(con booklit.TableOfContents) error {
	engine.template = engine.tmpl.Lookup("toc.tmpl")
	engine.data = con.Section
	return nil
}

func (engine *HTMLRenderingEngine) Render(out io.Writer) error {
	if engine.template == nil {
		return fmt.Errorf("unknown template for %T", engine.data)
	}

	return engine.template.Execute(out, engine.data)
}

func renderFunc(content booklit.Content) (template.HTML, error) {
	buf := new(bytes.Buffer)

	engine := NewHTMLRenderingEngine()

	err := content.Visit(engine)
	if err != nil {
		return "", err
	}

	err = engine.Render(buf)
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
