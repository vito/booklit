package stages

import "github.com/vito/booklit"

type Collect struct {
	Section *booklit.Section
}

func (collect *Collect) VisitString(booklit.String) error {
	return nil
}

func (collect *Collect) VisitSequence(con booklit.Sequence) error {
	for _, c := range con {
		err := c.Visit(collect)
		if err != nil {
			return err
		}
	}

	return nil
}

func (collect *Collect) VisitParagraph(con booklit.Paragraph) error {
	for _, c := range con {
		err := c.Visit(collect)
		if err != nil {
			return err
		}
	}

	return nil
}

func (collect *Collect) VisitPreformatted(con booklit.Preformatted) error {
	for _, c := range con {
		err := c.Visit(collect)
		if err != nil {
			return err
		}
	}

	return nil
}

func (collect *Collect) VisitReference(con *booklit.Reference) error {
	return nil
}

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

	// TODO: this probably does redundant resolving, since i think the section
	// was loaded via a processor in the first place
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

func (collect *Collect) VisitTableOfContents(booklit.TableOfContents) error {
	return nil
}

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

func (collect *Collect) VisitImage(con booklit.Image) error {
	return nil
}

func (collect *Collect) VisitTarget(con booklit.Target) error {
	collect.Section.SetTagAnchored(con.TagName, con.Title, con.Content, con.TagName)
	return nil
}

func (collect *Collect) VisitList(con booklit.List) error {
	for _, c := range con.Items {
		err := c.Visit(collect)
		if err != nil {
			return err
		}
	}

	return nil
}

func (collect *Collect) VisitLink(con booklit.Link) error {
	return con.Content.Visit(collect)
}

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
