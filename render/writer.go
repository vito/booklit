package render

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

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

	SaveSearchIndex bool
}

const SearchIndexFilename = "search_index.json"

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
		name := section.PrimaryTag.Name + "." + writer.Engine.FileExtension()
		path := filepath.Join(writer.Destination, name)

		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}

		err = writer.Engine.RenderSection(file, section)
		if err != nil {
			return err
		}
	}

	if writer.SaveSearchIndex {
		indexPath := filepath.Join(writer.Destination, SearchIndexFilename)

		indexFile, err := os.OpenFile(indexPath, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return err
		}

		index := SearchIndex{}
		err = json.NewDecoder(indexFile).Decode(&index)
		if err != nil && err != io.EOF {
			return err
		}

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

		err = indexFile.Truncate(0)
		if err != nil {
			return err
		}

		_, err = indexFile.Seek(0, 0)
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
	}

	for _, child := range section.Children {
		err := writer.WriteSection(child)
		if err != nil {
			return err
		}
	}

	return nil
}
