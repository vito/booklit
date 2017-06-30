package render

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"path/filepath"

	"github.com/vito/booklit"
)

var tmpl *template.Template

func init() {
	tmpl = template.New("engine").Funcs(template.FuncMap{
		"render": renderFunc,
	})

	for _, asset := range AssetNames() {
		info, err := AssetInfo(asset)
		if err != nil {
			panic(err)
		}

		tmpl.New(filepath.Base(info.Name())).Parse(string(MustAsset(asset)))
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

func (engine *HTMLRenderingEngine) VisitString(str booklit.String) {
	engine.template = tmpl.Lookup("string.html")
	engine.data = str
}

func (engine *HTMLRenderingEngine) VisitSection(str *booklit.Section) {
	engine.template = tmpl.Lookup("section.html")
	engine.data = str
}

func (engine *HTMLRenderingEngine) VisitSequence(str booklit.Sequence) {
	engine.template = tmpl.Lookup("sequence.html")
	engine.data = str
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
	content.Visit(engine)
	err := engine.Render(buf)
	if err != nil {
		return "", err
	}

	return template.HTML(buf.String()), nil
}
