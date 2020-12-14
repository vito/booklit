package booklit

import "fmt"

// Image embeds an image as flow content.
type Image struct {
	// File path or URL.
	Path string

	// Description of the image.
	Description string
}

func (con Image) IsFlow() bool {
	return true
}

func (con Image) String() string {
	return fmt.Sprintf("{image: %s}", con.Path)
}

func (con Image) Visit(visitor Visitor) error {
	return visitor.VisitImage(con)
}
