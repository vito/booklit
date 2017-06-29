package render

import (
	"bytes"
	"fmt"
	"html/template"
	"io"

	"github.com/vito/booklit"
)

var tmpl = template.Must(template.New("engine").Funcs(template.FuncMap{
	"render": renderFunc,
}).ParseGlob("render/html/*.html"))

type HTMLRenderingEngine struct {
	template *template.Template
	data     interface{}
}

func NewHTMLRenderingEngine() *HTMLRenderingEngine {
	return &HTMLRenderingEngine{}
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
