package render

import (
	"io"

	"github.com/vito/booklit"
)

// Engine is the primary type for rendering sections.
type Engine interface {
	// RenderSection writes the given section to the writer.
	RenderSection(io.Writer, *booklit.Section) error

	// The canonical file extension for files written by the engine,
	// without the leading dot.
	FileExtension() string

	// The URL to reference rendered content for a given tag.
	URL(booklit.Tag) string
}
