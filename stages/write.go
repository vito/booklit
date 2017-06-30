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

func (stage Write) VisitString(booklit.String) error { return nil }

func (stage Write) VisitSequence(booklit.Sequence) error { return nil }

func (stage Write) VisitParagraph(booklit.Paragraph) error { return nil }

func (stage Write) VisitSection(section *booklit.Section) error {
	if section.Parent != nil {
		// TODO: or, if parent is not configured for split sections
		return nil
	}

	name := section.PrimaryTag() + "." + stage.Engine.FileExtension()
	path := filepath.Join(stage.Destination, name)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	err = section.Visit(stage.Engine)
	if err != nil {
		return err
	}

	return stage.Engine.Render(file)
}
