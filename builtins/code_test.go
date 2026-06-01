//go:build cgo

package builtins

import (
	"testing"

	"github.com/vito/booklit"
)

func TestSyntaxLinksBooklitComponentTags(t *testing.T) {
	section := &booklit.Section{}

	content, err := syntax(section, "markdown", booklit.Preformatted{booklit.String(`<IncludeSection path="quotes.md"/>`)})
	if err != nil {
		t.Fatal(err)
	}

	styled, ok := content.(booklit.Styled)
	if !ok {
		t.Fatalf("expected styled content, got %T", content)
	}
	seq, ok := styled.Content.(booklit.Sequence)
	if !ok {
		t.Fatalf("expected sequence content, got %T", styled.Content)
	}

	for _, item := range seq {
		ref, ok := item.(*booklit.Reference)
		if !ok {
			continue
		}
		if ref.TagName == "include-section" && ref.Content.String() == "IncludeSection" && ref.Optional {
			return
		}
	}
	t.Fatalf("expected IncludeSection reference in %#v", seq)
}
