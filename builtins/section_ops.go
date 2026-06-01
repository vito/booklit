package builtins

import (
	"path/filepath"

	"github.com/vito/booklit"
	"github.com/vito/booklit/ast"
)

func init() {
	Register("Styled", styledFunc)
	Register("IncludeSection", includeSectionFunc)
	Register("SinglePage", singlePageFunc)
	Register("SplitSections", splitSectionsFunc)
	Register("OmitChildrenFromTableOfContents", omitChildrenFromTOCFunc)
	Register("TableOfContents", tableOfContentsFunc)
}

// styledFunc — `<Styled name="..."/>`. Sets the current section's Style;
// mirrors \styled{name}.
func styledFunc(ctx *Context, props map[string]ast.Node, _ []ast.Node) (booklit.Content, error) {
	name, err := requireStringProp(ctx, props, "name", "Styled")
	if err != nil {
		return nil, err
	}
	ctx.Section.Style = name
	return nil, nil
}

// includeSectionFunc — `<IncludeSection path="..."/>`. Mirrors
// \include-section{path}.
func includeSectionFunc(ctx *Context, props map[string]ast.Node, _ []ast.Node) (booklit.Content, error) {
	path, err := requireStringProp(ctx, props, "path", "IncludeSection")
	if err != nil {
		return nil, err
	}
	sectionPath := filepath.Join(filepath.Dir(ctx.Section.FilePath()), path)
	section, err := ctx.Section.Processor.EvaluateFile(ctx.Section, sectionPath)
	if err != nil {
		return nil, err
	}
	ctx.Section.Children = append(ctx.Section.Children, section)
	return nil, nil
}

// singlePageFunc — `<SinglePage/>`. Mirrors \single-page.
func singlePageFunc(ctx *Context, _ map[string]ast.Node, _ []ast.Node) (booklit.Content, error) {
	ctx.Section.PreventSplitSections = true
	return nil, nil
}

// splitSectionsFunc — `<SplitSections/>`. Mirrors \split-sections.
func splitSectionsFunc(ctx *Context, _ map[string]ast.Node, _ []ast.Node) (booklit.Content, error) {
	ctx.Section.ResetDepth = true
	if !ctx.Section.SplitSectionsPrevented() {
		ctx.Section.SplitSections = true
	}
	return nil, nil
}

// omitChildrenFromTOCFunc — `<OmitChildrenFromTableOfContents/>`. Mirrors
// \omit-children-from-table-of-contents.
func omitChildrenFromTOCFunc(ctx *Context, _ map[string]ast.Node, _ []ast.Node) (booklit.Content, error) {
	ctx.Section.OmitChildrenFromTableOfContents = true
	return nil, nil
}

// tableOfContentsFunc — `<TableOfContents/>`. Mirrors \table-of-contents.
func tableOfContentsFunc(ctx *Context, _ map[string]ast.Node, _ []ast.Node) (booklit.Content, error) {
	return booklit.TableOfContents{Section: ctx.Section}, nil
}
