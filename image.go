package booklit

import "fmt"

// Image embeds an image as flow content.
type Image struct {
	// File path or URL.
	Path string

	// Description of the image.
	Description string
}

// IsFlow returns true.
func (con Image) IsFlow() bool {
	return true
}

// String summarizes the content for debugging purposes.
func (con Image) String() string {
	return fmt.Sprintf("{image: %s}", con.Path)
}

// Visit calls VisitImage.
func (con Image) Visit(visitor Visitor) error {
	return visitor.VisitImage(con)
}
