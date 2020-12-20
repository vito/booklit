package stages

import "github.com/vito/booklit"

// Collect is a booklit.Visitor that collects all tag targets and sets them on
// the section.
type Collect struct {
	Section *booklit.Section
}

// VisitString does nothing.
func (collect *Collect) VisitString(booklit.String) error {
	return nil
}

// VisitSequence visits each content in the sequence.
func (collect *Collect) VisitSequence(con booklit.Sequence) error {
	for _, c := range con {
		err := c.Visit(collect)
		if err != nil {
			return err
		}
	}

	return nil
}

// VisitParagraph visits each line.
func (collect *Collect) VisitParagraph(con booklit.Paragraph) error {
	for _, c := range con {
		err := c.Visit(collect)
		if err != nil {
			return err
		}
	}

	return nil
}

// VisitPreformatted visits each line.
func (collect *Collect) VisitPreformatted(con booklit.Preformatted) error {
	for _, c := range con {
		err := c.Visit(collect)
		if err != nil {
			return err
		}
	}

	return nil
}

// VisitReference does nothing.
func (collect *Collect) VisitReference(con *booklit.Reference) error {
	return nil
}

// VisitSection visits the section's title, body, and partials, followed by
// each child section which is visited with their own *Collect.
func (collect *Collect) VisitSection(con *booklit.Section) error {
	err := con.Title.Visit(collect)
	if err != nil {
		return err
	}

	err = con.Body.Visit(collect)
	if err != nil {
		return err
	}

	for _, p := range con.Partials {
		err = p.Visit(collect)
		if err != nil {
			return err
		}
	}

	for _, child := range con.Children {
		subCollector := &Collect{
			Section: child,
		}

		err := child.Visit(subCollector)
		if err != nil {
			return err
		}
	}

	return nil
}

// VisitTableOfContents does nothing.
func (collect *Collect) VisitTableOfContents(booklit.TableOfContents) error {
	return nil
}

// VisitStyled visits the content and partials.
func (collect *Collect) VisitStyled(con booklit.Styled) error {
	err := con.Content.Visit(collect)
	if err != nil {
		return err
	}

	for _, v := range con.Partials {
		if v == nil {
			continue
		}

		err := v.Visit(collect)
		if err != nil {
			return err
		}
	}

	return nil
}

// VisitImage does nothing.
func (collect *Collect) VisitImage(con booklit.Image) error {
	return nil
}

// VisitTarget sets an anchored tag on the Section.
func (collect *Collect) VisitTarget(con booklit.Target) error {
	collect.Section.SetTagAnchored(con.TagName, con.Title, con.Location, con.Content, con.TagName)
	return nil
}

// VisitList visits each item in the list.
func (collect *Collect) VisitList(con booklit.List) error {
	for _, c := range con.Items {
		err := c.Visit(collect)
		if err != nil {
			return err
		}
	}

	return nil
}

// VisitLink visits the link's content.
func (collect *Collect) VisitLink(con booklit.Link) error {
	return con.Content.Visit(collect)
}

// VisitTable visits each table cell.
func (collect *Collect) VisitTable(con booklit.Table) error {
	for _, row := range con.Rows {
		for _, c := range row {
			err := c.Visit(collect)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// VisitDefinitions visits each subject and definition.
func (collect *Collect) VisitDefinitions(con booklit.Definitions) error {
	for _, def := range con {
		err := def.Subject.Visit(collect)
		if err != nil {
			return err
		}

		err = def.Definition.Visit(collect)
		if err != nil {
			return err
		}
	}

	return nil
}
