package stages

import (
	"github.com/sirupsen/logrus"
	"github.com/vito/booklit"
)

// Resolve is a booklit.Visitor that locates the target tag for each
// booklit.Reference.
type Resolve struct {
	AllowBrokenReferences bool
}

// VisitString does nothing.
func (resolve *Resolve) VisitString(booklit.String) error {
	return nil
}

// VisitSequence visits each content in the sequence.
func (resolve *Resolve) VisitSequence(con booklit.Sequence) error {
	for _, c := range con {
		err := c.Visit(resolve)
		if err != nil {
			return err
		}
	}

	return nil
}

// VisitParagraph visits each line.
func (resolve *Resolve) VisitParagraph(con booklit.Paragraph) error {
	for _, c := range con {
		err := c.Visit(resolve)
		if err != nil {
			return err
		}
	}

	return nil
}

// VisitPreformatted visits each line.
func (resolve *Resolve) VisitPreformatted(con booklit.Preformatted) error {
	for _, c := range con {
		err := c.Visit(resolve)
		if err != nil {
			return err
		}
	}

	return nil
}

// VisitReference finds the referenced tag through Section.FindTag.
//
// If 1 tag is found, it is assigned on the Reference.
//
// If 0 tags are found and AllowBrokenReferences is false,
// booklit.UnknownTagError is returned.
//
// If more than one tag is returned and AllowBrokenReferences is false,
// booklit.AmbiguousReferenceError is returned.
//
// If AllowBrokenReferences is true, a made-up tag is assigned on the Reference
// instead of returning either of the above errors.
func (resolve *Resolve) VisitReference(con *booklit.Reference) error {
	_, err := con.Tag()
	if err != nil && resolve.AllowBrokenReferences {
		logrus.WithFields(logrus.Fields{"section": con.Section.Path}).
			Warnf("broken reference: %s", err)
		return nil
	}

	return err
}

// VisitSection visits the section's title, body, and partials, followed by
// each child section which is visited with their own *Resolve.
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

	for _, child := range con.Children {
		err := child.Visit(resolve)
		if err != nil {
			return err
		}
	}

	return nil
}

// VisitTableOfContents does nothing.
func (resolve *Resolve) VisitTableOfContents(booklit.TableOfContents) error {
	return nil
}

// VisitStyled visits the content and partials.
func (resolve *Resolve) VisitStyled(con booklit.Styled) error {
	err := con.Content.Visit(resolve)
	if err != nil {
		return err
	}

	for _, v := range con.Partials {
		if v == nil {
			continue
		}

		err := v.Visit(resolve)
		if err != nil {
			return err
		}
	}

	return nil
}

// VisitTarget visits the target's title and content.
func (resolve *Resolve) VisitTarget(con booklit.Target) error {
	err := con.Title.Visit(resolve)
	if err != nil {
		return err
	}

	if con.Content != nil {
		err := con.Content.Visit(resolve)
		if err != nil {
			return err
		}
	}

	return nil
}

// VisitImage does nothing.
func (resolve *Resolve) VisitImage(con booklit.Image) error {
	return nil
}

// VisitList visits each item in the list.
func (resolve *Resolve) VisitList(con booklit.List) error {
	for _, c := range con.Items {
		err := c.Visit(resolve)
		if err != nil {
			return err
		}
	}

	return nil
}

// VisitLink visits the link's content.
func (resolve *Resolve) VisitLink(con booklit.Link) error {
	return con.Content.Visit(resolve)
}

// VisitTable visits each table cell.
func (resolve *Resolve) VisitTable(con booklit.Table) error {
	for _, row := range con.Rows {
		for _, c := range row {
			err := c.Visit(resolve)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// VisitDefinitions visits each subject and definition.
func (resolve *Resolve) VisitDefinitions(con booklit.Definitions) error {
	for _, def := range con {
		err := def.Subject.Visit(resolve)
		if err != nil {
			return err
		}

		err = def.Definition.Visit(resolve)
		if err != nil {
			return err
		}
	}

	return nil
}

// VisitLazy does nothing.
func (resolve *Resolve) VisitLazy(con *booklit.Lazy) error {
	return nil
}
