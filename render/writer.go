package render

import (
	"io"
	"os"
	"path/filepath"

	"github.com/vito/booklit"
)

type RenderingEngine interface {
	booklit.Visitor

	FileExtension() string
	RenderSection(io.Writer, *booklit.Section) error
}

type Writer struct {
	Engine RenderingEngine

	Destination string
}

func (writer Writer) WriteSection(section *booklit.Section) error {
	if section.Parent == nil || section.Parent.SplitSections {
		name := section.PrimaryTag.Name + "." + writer.Engine.FileExtension()
		path := filepath.Join(writer.Destination, name)

		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}

		err = writer.Engine.RenderSection(file, section)
		if err != nil {
			return err
		}
	}

	errs := make(chan error, len(section.Children))
	for _, child := range section.Children {
		go func() {
			errs <- writer.WriteSection(child)
		}()
	}

	var anyErr error
	for i := 0; i < len(section.Children); i++ {
		err := <-errs
		if err != nil {
			anyErr = err
		}
	}

	return anyErr
}
