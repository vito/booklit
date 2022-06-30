package booklit

import (
	"fmt"

	"github.com/vito/booklit/ast"
)

// Reference is flow content linking to a tag defined elsewhere.
type Reference struct {
	// The section to resolve the reference within.
	Section *Section

	// The tag to link to.
	TagName string

	// Optional content to display for the reference. If not present, the tag's
	// own display content will be used.
	Content Content

	// If set to true and resolving the tag does not succeed, display Content
	// instead of displaying a link.
	Optional bool

	// The original source location of the reference. Used when generating error
	// messages.
	Location ast.Location

	// The result of Resolving the reference.
	tag    *Tag
	broken bool
}

// IsFlow returns true.
func (con *Reference) IsFlow() bool {
	return true
}

// String summarizes the content for debugging purposes.
func (con *Reference) String() string {
	if con.Content != nil {
		return con.Content.String()
	}

	tag, err := con.Tag()
	if err != nil {
		return fmt.Sprintf("{broken reference: %s: %s}", con.TagName, err)
	}

	return tag.Title.String()
}

// Visit calls VisitReference.
func (con *Reference) Visit(visitor Visitor) error {
	return visitor.VisitReference(con)
}

// Display returns the content to display for the reference. If Content is set,
// it is returned, otherwise the resolved tag's Title is returned.
func (con *Reference) Display() (Content, error) {
	if con.Content != nil {
		return con.Content, nil
	}

	tag, err := con.Tag()
	if err != nil {
		return nil, err
	}

	return tag.Title, nil
}

func (con *Reference) Tag() (*Tag, error) {
	if con.tag != nil {
		return con.tag, nil
	}

	err := con.resolve()
	if err != nil {
		return nil, err
	}

	return con.tag, nil
}

// Resolve resolves the reference to a tag, returning an error if the tag is
// unknown or ambiguous.
func (con *Reference) resolve() error {
	tags := con.Section.FindTag(con.TagName)

	switch len(tags) {
	case 1: // unnambiguous reference
		con.tag = &tags[0]
		return nil

	case 0: // broken reference
		if con.Optional {
			// allow nonexistant tag; template must handle nil Tag
			return nil
		}

		return UnknownTagError{
			TagName:     con.TagName,
			SimilarTags: con.Section.SimilarTags(con.TagName),
			ErrorLocation: ErrorLocation{
				FilePath:     con.Section.FilePath(),
				NodeLocation: con.Location,
				Length:       len("\\reference"),
			},
		}

	default: // ambiguous reference
		locs := []ErrorLocation{}
		for _, t := range tags {
			locs = append(locs, ErrorLocation{
				FilePath:     t.Section.FilePath(),
				NodeLocation: t.Location,
			})
		}

		return AmbiguousReferenceError{
			TagName:          con.TagName,
			DefinedLocations: locs,
			ErrorLocation: ErrorLocation{
				FilePath:     con.Section.FilePath(),
				NodeLocation: con.Location,
				Length:       len("\\reference"),
			},
		}
	}
}
