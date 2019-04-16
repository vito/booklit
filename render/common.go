package render

import "github.com/vito/booklit"

type WalkContext struct {
	Current *booklit.Section
	Section *booklit.Section
}

func sectionPageOwner(section *booklit.Section) *booklit.Section {
	if section.Parent == nil {
		return section
	}

	if section.Parent.SplitSections {
		return section
	}

	return sectionPageOwner(section.Parent)
}

func sectionURL(ext string, section *booklit.Section, anchor string) string {
	owner := sectionPageOwner(section)

	if owner != section {
		if anchor == "" {
			anchor = section.PrimaryTag.Name
		}

		return sectionURL(ext, owner, anchor)
	}

	filename := section.PrimaryTag.Name + "." + ext

	if anchor != "" {
		filename += "#" + anchor
	}

	return filename
}
