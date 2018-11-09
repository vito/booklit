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
	strip.Result = Sequence(stripAuxSeq(con))
	return nil
}

func (strip *stripAuxVisitor) VisitParagraph(con Paragraph) error {
	strip.Result = Paragraph(stripAuxSeq(con))
	return nil
}

func (strip *stripAuxVisitor) VisitPreformatted(con Preformatted) error {
	strip.Result = Preformatted(stripAuxSeq(con))
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

	strippedPartials := Partials{}
	for k, v := range con.Partials {
		if v == nil {
			continue
		}

		strippedPartials[k] = StripAux(v)
	}

	con.Partials = strippedPartials

	strip.Result = con
	return nil
}

func (strip *stripAuxVisitor) VisitTarget(con Target) error {
	con.Title = StripAux(con.Title)
	con.Content = StripAux(con.Content)
	strip.Result = con
	return nil
}

func (strip *stripAuxVisitor) VisitImage(con Image) error {
	strip.Result = con
	return nil
}

func (strip *stripAuxVisitor) VisitList(con List) error {
	con.Items = stripAuxSeq(con.Items)
	strip.Result = con
	return nil
}

func (strip *stripAuxVisitor) VisitLink(con Link) error {
	con.Content = StripAux(con.Content)
	strip.Result = con
	return nil
}

func (strip *stripAuxVisitor) VisitTable(con Table) error {
	newTable := Table{}
	for _, row := range con.Rows {
		newTable.Rows = append(newTable.Rows, stripAuxSeq(row))
	}

	strip.Result = newTable
	return nil
}

func (strip *stripAuxVisitor) VisitDefinitions(con Definitions) error {
	stripped := Definitions{}
	for _, def := range con {
		stripped = append(stripped, Definition{
			Subject:    StripAux(def.Subject),
			Definition: StripAux(def.Definition),
		})
	}

	strip.Result = stripped
	return nil
}

func stripAuxSeq(seq []Content) []Content {
	stripped := []Content{}

	for _, c := range seq {
		if _, isAux := c.(Aux); !isAux {
			stripped = append(stripped, StripAux(c))
		}
	}

	return stripped
}
