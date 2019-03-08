package render

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/vito/booklit"
)

type RenderingEngine interface {
	booklit.Visitor

	FileExtension() string
	RenderSection(io.Writer, *booklit.Section) error
	URL(booklit.Tag) string
}

type Writer struct {
	Engine RenderingEngine

	Destination string
}

type SearchIndex map[string]SearchDocument

type SearchDocument struct {
	Location   string `json:"location"`
	Title      string `json:"title"`
	Text       string `json:"text"`
	Depth      int    `json:"depth"`
	SectionTag string `json:"section_tag"`
}

func (writer Writer) WriteSection(section *booklit.Section) error {
	if section.Parent == nil || section.Parent.SplitSections {
		err := writer.writeSingleSection(section)
		if err != nil {
			return err
		}
	}

	for _, child := range section.Children {
		err := writer.WriteSection(child)
		if err != nil {
			return err
		}
	}

	return nil
}

func (writer Writer) WriteSearchIndex(section *booklit.Section, path string) error {
	logrus.WithFields(logrus.Fields{
		"path": path,
	}).Infoln("writing search index")

	indexPath := filepath.Join(writer.Destination, path)

	index := SearchIndex{}
	writer.loadTags(index, section)

	indexFile, err := os.Create(indexPath)
	if err != nil {
		return err
	}

	err = json.NewEncoder(indexFile).Encode(index)
	if err != nil {
		return err
	}

	err = indexFile.Close()
	if err != nil {
		return err
	}

	return nil
}

func (writer Writer) loadTags(index SearchIndex, section *booklit.Section) {
	for _, tag := range section.Tags {
		var text string
		if tag.Content != nil {
			text = tag.Content.String()
		} else {
			text = tag.Section.Body.String()
		}

		index[tag.Name] = SearchDocument{
			Location:   writer.Engine.URL(tag),
			Title:      tag.Title.String(),
			Text:       text,
			Depth:      tag.Section.Depth(),
			SectionTag: tag.Section.PrimaryTag.Name,
		}
	}

	for _, child := range section.Children {
		writer.loadTags(index, child)
	}
}

func (writer Writer) writeSingleSection(section *booklit.Section) error {
	name := section.PrimaryTag.Name + "." + writer.Engine.FileExtension()
	path := filepath.Join(writer.Destination, name)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	logrus.WithFields(logrus.Fields{
		"section":  section.Path,
		"rendered": path,
	}).Info("rendering")

	err = writer.Engine.RenderSection(file, section)
	if err != nil {
		return err
	}

	return nil
}
