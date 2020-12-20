package render

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/vito/booklit"
)

// Writer writes rendered content using an Engine to the given
// destination.
type Writer struct {
	Engine Engine

	Destination string
}

// SearchIndex is a mapping from tag names to a summary useful for
// inline search.
type SearchIndex map[string]SearchDocument

// SearchDocument contains data useful for implementing inline
// search.
type SearchDocument struct {
	// The tag's URL.
	Location string `json:"location"`

	// The title of the tag.
	Title string `json:"title"`

	// The text content for the tag, or the section's text if the tag
	// does not have its own content.
	Text string `json:"text"`

	// The depth of the tag's section.
	Depth int `json:"depth"`

	// The containing section's primary tag.
	SectionTag string `json:"section_tag"`
}

// WriteSection renders the given section to disk if it has no
// parent or if its parent is configured with SplitSections.
//
// After rendering, WriteSection recurses to the section's Children.
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

// WriteSearchIndex generates and writes a SearchIndex in JSON
// format to the given path within the destination.
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
