package stages

import (
	"fmt"

	"github.com/vito/booklit"
)

type Resolve struct {
	AllowBrokenReferences bool

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

func (resolve *Resolve) VisitPreformatted(con booklit.Preformatted) error {
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
	if found {
		con.Tag = &tag
	} else if resolve.AllowBrokenReferences {
		con.Tag = &booklit.Tag{
			Name:    con.TagName,
			Anchor:  "broken",
			Display: booklit.String(fmt.Sprintf("{broken reference: %s}", con.TagName)),
			Section: resolve.Section,
		}
	} else {
		return fmt.Errorf("could not find tag '%s' from within section '%s'", con.TagName, resolve.Section.String())
	}

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

	for _, p := range con.Partials {
		err = p.Visit(resolve)
		if err != nil {
			return err
		}
	}

	// TODO: this probably does redundant resolving, since i think the section
	// was loaded via a processor in the first place
	for _, child := range con.Children {
		subResolver := &Resolve{
			AllowBrokenReferences: resolve.AllowBrokenReferences,
			Section:               child,
		}

		err := child.Visit(subResolver)
		if err != nil {
			return err
		}
	}

	return nil
}

func (resolve *Resolve) VisitTableOfContents(booklit.TableOfContents) error {
	return nil
}

func (resolve *Resolve) VisitStyled(con booklit.Styled) error {
	return con.Content.Visit(resolve)
}

func (resolve *Resolve) VisitTarget(con booklit.Target) error {
	return con.Display.Visit(resolve)
}

func (resolve *Resolve) VisitBlock(con booklit.Block) error {
	return con.Content.Visit(resolve)
}

func (resolve *Resolve) VisitElement(con booklit.Element) error {
	return con.Content.Visit(resolve)
}

func (resolve *Resolve) VisitImage(con booklit.Image) error {
	return nil
}

func (resolve *Resolve) VisitList(con booklit.List) error {
	for _, c := range con.Items {
		err := c.Visit(resolve)
		if err != nil {
			return err
		}
	}

	return nil
}

func (resolve *Resolve) VisitLink(con booklit.Link) error {
	return con.Content.Visit(resolve)
}
