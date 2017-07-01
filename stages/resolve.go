package stages

import (
	"fmt"

	"github.com/vito/booklit"
)

type Resolve struct {
	Section *booklit.Section
}

func (resolve *Resolve) VisitString(booklit.String) error {
	return nil
}

func (resolve *Resolve) VisitSequence(con booklit.Sequence) error {
	for _, c := range con {
		err := c.Visit(resolve)
		if err != nil {
			return err
		}
	}

	return nil
}

func (resolve *Resolve) VisitParagraph(con booklit.Paragraph) error {
	for _, c := range con {
		err := c.Visit(resolve)
		if err != nil {
			return err
		}
	}

	return nil
}

func (resolve *Resolve) VisitReference(con *booklit.Reference) error {
	tag, found := resolve.Section.FindTag(con.TagName)
	if !found {
		return fmt.Errorf("could not find tag '%s' from within section '%s'", con.TagName, resolve.Section.String())
	}

	con.Tag = &tag

	return nil
}

func (resolve *Resolve) VisitSection(con *booklit.Section) error {
	err := con.Title.Visit(resolve)
	if err != nil {
		return err
	}

	err = con.Body.Visit(resolve)
	if err != nil {
		return err
	}

	// TODO: this probably does redundant resolving, since i think the section
	// was loaded via a processor in the first place
	for _, child := range con.Children {
		subResolver := &Resolve{
			Section: child,
		}

		err := child.Visit(subResolver)
		if err != nil {
			return err
		}
	}

	return nil
}
