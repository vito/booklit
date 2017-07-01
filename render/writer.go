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
	Render(io.Writer) error
}

type Writer struct {
	Engine RenderingEngine

	Destination string
}

func (writer Writer) WriteSection(section *booklit.Section) error {
	if section.Parent != nil && !section.Parent.SplitSections {
		return nil
	}

	name := section.PrimaryTag.Name + "." + writer.Engine.FileExtension()
	path := filepath.Join(writer.Destination, name)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	var node booklit.Content = section
	if section.SplitSections {
		for _, child := range section.Children {
			err := writer.WriteSection(child)
			if err != nil {
				return err
			}
		}
	} else {
		for _, child := range section.Children {
			node = booklit.Append(node, child)
		}
	}

	err = node.Visit(writer.Engine)
	if err != nil {
		return err
	}

	return writer.Engine.Render(file)
}
