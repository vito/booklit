package booklit

type Aux struct {
	Content
}

func StripAux(content Content) Content {
	visitor := &stripAuxVisitor{}

	_ = content.Visit(visitor)

	return visitor.Result
}

type stripAuxVisitor struct {
	Result Content
}

func (strip *stripAuxVisitor) VisitString(con String) error {
	strip.Result = con
	return nil
}

func (strip *stripAuxVisitor) VisitSequence(con Sequence) error {
	stripped := Sequence{}

	for _, c := range con {
		if _, isAux := c.(Aux); isAux {
			continue
		}

		stripped = append(stripped, StripAux(c))
	}

	strip.Result = stripped
	return nil
}

func (strip *stripAuxVisitor) VisitParagraph(con Paragraph) error {
	stripped := Paragraph{}

	for _, c := range con {
		if _, isAux := c.(Aux); !isAux {
			stripped = append(stripped, StripAux(c))
		}
	}

	strip.Result = stripped
	return nil
}

func (strip *stripAuxVisitor) VisitPreformatted(con Preformatted) error {
	stripped := Preformatted{}

	for _, c := range con {
		if _, isAux := c.(Aux); !isAux {
			stripped = append(stripped, StripAux(c))
		}
	}

	strip.Result = stripped
	return nil
}

func (strip *stripAuxVisitor) VisitReference(con *Reference) error {
	ref := *con
	ref.Content = StripAux(ref.Content)
	strip.Result = &ref
	return nil
}

func (strip *stripAuxVisitor) VisitSection(con *Section) error {
	strip.Result = con
	return nil
}

func (strip *stripAuxVisitor) VisitTableOfContents(con TableOfContents) error {
	strip.Result = con
	return nil
}

func (strip *stripAuxVisitor) VisitStyled(con Styled) error {
	con.Content = StripAux(con.Content)
	strip.Result = con
	return nil
}

func (strip *stripAuxVisitor) VisitTarget(con Target) error {
	con.Display = StripAux(con.Display)
	strip.Result = con
	return nil
}

func (strip *stripAuxVisitor) VisitBlock(con Block) error {
	con.Content = StripAux(con.Content)
	strip.Result = con
	return nil
}

func (strip *stripAuxVisitor) VisitElement(con Element) error {
	con.Content = StripAux(con.Content)
	strip.Result = con
	return nil
}

func (strip *stripAuxVisitor) VisitImage(con Image) error {
	strip.Result = con
	return nil
}

func (strip *stripAuxVisitor) VisitList(con List) error {
	stripped := []Content{}

	for _, c := range con.Items {
		if _, isAux := c.(Aux); !isAux {
			stripped = append(stripped, StripAux(c))
		}
	}

	con.Items = stripped
	strip.Result = con
	return nil
}

func (strip *stripAuxVisitor) VisitLink(con Link) error {
	con.Content = StripAux(con.Content)
	strip.Result = con
	return nil
}
