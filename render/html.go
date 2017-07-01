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

var tmpl *template.Template

func init() {
	tmpl = template.New("engine").Funcs(template.FuncMap{
		"render": renderFunc,
		"url":    tagURLFunc,
		"sectionHeader": func(con *booklit.Section, content template.HTML) template.HTML {
			depth := 1
			for sec := con; sec.Parent != nil && !sec.Parent.SplitSections; sec = sec.Parent {
				depth++
			}

			if depth > 6 {
				depth = 6
			}

			return template.HTML(fmt.Sprintf("<h%d>%s</h%d>", depth, content, depth))
		},
	})

	for _, asset := range AssetNames() {
		info, err := AssetInfo(asset)
		if err != nil {
			panic(err)
		}

		content := strings.TrimRight(string(MustAsset(asset)), "\n")

		_, err = tmpl.New(filepath.Base(info.Name())).Parse(content)
		if err != nil {
			panic(err)
		}
	}
}

type HTMLRenderingEngine struct {
	template *template.Template
	data     interface{}
}

func NewHTMLRenderingEngine() *HTMLRenderingEngine {
	return &HTMLRenderingEngine{}
}

func (engine *HTMLRenderingEngine) FileExtension() string {
	return "html"
}

func (engine *HTMLRenderingEngine) VisitString(con booklit.String) error {
	engine.template = tmpl.Lookup("string.html")
	engine.data = con
	return nil
}

func (engine *HTMLRenderingEngine) VisitReference(con *booklit.Reference) error {
	engine.template = tmpl.Lookup("reference.html")
	engine.data = con
	return nil
}

func (engine *HTMLRenderingEngine) VisitSection(con *booklit.Section) error {
	engine.template = tmpl.Lookup("section.html")
	engine.data = con
	return nil
}

func (engine *HTMLRenderingEngine) VisitSequence(con booklit.Sequence) error {
	engine.template = tmpl.Lookup("sequence.html")
	engine.data = con
	return nil
}

func (engine *HTMLRenderingEngine) VisitParagraph(con booklit.Paragraph) error {
	engine.template = tmpl.Lookup("paragraph.html")
	engine.data = con
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
