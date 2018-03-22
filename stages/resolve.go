package stages

import (
	"fmt"

	"github.com/sirupsen/logrus"
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
	tags := resolve.Section.FindTag(con.TagName)

	var err error
	switch len(tags) {
	case 0:
		logrus.WithFields(logrus.Fields{
			"section": resolve.Section.Path,
		}).Warnf("broken reference: %s", con.TagName)

		err = booklit.ErrBrokenReference
	case 1:
		con.Tag = &tags[0]
	default:
		paths := []string{}
		for _, t := range tags {
			paths = append(paths, t.Section.Path)
		}

		logrus.WithFields(logrus.Fields{
			"section":    resolve.Section.Path,
			"defined-in": paths,
		}).Warnf("ambiguous reference: %s", con.TagName)

		err = booklit.ErrBrokenReference
	}

	if err == nil {
		return nil
	}

	if resolve.AllowBrokenReferences {
		con.Tag = &booklit.Tag{
			Name:    con.TagName,
			Anchor:  "broken",
			Title:   booklit.String(fmt.Sprintf("{broken reference: %s}", con.TagName)),
			Section: resolve.Section,
		}

		return nil
	}

	return err
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
