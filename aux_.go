package booklit

// Aux is auxiliary content, typically in section titles, which can be removed
// in certain contexts such as references or table-of-content lists.
//
// It is primarily used with the `stripAux` template function.
//
// See https://booklit.page/baselit.html#aux for more information.
type Aux struct {
	Content
}

// StripAux recurses through Content and removes any Aux content from any
// nested sequences.
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

func (strip *stripAuxVisitor) VisitTarget(con Target) error {
	con.Title = StripAux(con.Title)
	con.Content = StripAux(con.Content)
	strip.Result = con
	return nil
}

func (strip *stripAuxVisitor) VisitLazy(con *Lazy) error {
	strip.Result = con
	return nil
}

func (strip *stripAuxVisitor) VisitRawElement(con RawElement) error {
	if con.Content != nil {
		con.Content = StripAux(con.Content)
	}
	strip.Result = con
	return nil
}

func (strip *stripAuxVisitor) VisitRawFragment(con RawFragment) error {
	strip.Result = con
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
