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
	return con.Content.Visit(collect)
}

func (collect *Collect) VisitBlock(con booklit.Block) error {
	return con.Content.Visit(collect)
}

func (collect *Collect) VisitTarget(con booklit.Target) error {
	collect.Section.Tags = append(collect.Section.Tags, booklit.Tag{
		Name:    con.TagName,
		Section: collect.Section,
		Display: con.Display,
		Anchor:  con.TagName,
	})

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
