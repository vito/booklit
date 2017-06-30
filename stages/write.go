package stages

import (
	"io"
	"os"
	"path/filepath"

	"github.com/vito/booklit"
)

type RenderingEngine interface {
	booklit.Visitor

	FileExtension() string
	Render(io.Writer) error
}

type Write struct {
	Engine RenderingEngine

	Destination string
}

func (stage Write) VisitString(str booklit.String) {}

func (stage Write) VisitSequence(seq booklit.Sequence) {}

func (stage Write) VisitSection(section *booklit.Section) {
	if section.Parent != nil {
		// TODO: or, if parent is not configured for split sections
		return
	}

	name := section.PrimaryTag() + "." + stage.Engine.FileExtension()
	path := filepath.Join(stage.Destination, name)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		// TODO: visitor should really be able to error
		panic(err)
	}

	section.Visit(stage.Engine)

	err = stage.Engine.Render(file)
	if err != nil {
		// TODO: visitor should really be able to error
		panic(err)
	}
}
