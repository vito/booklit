//go:build cgo

package builtins

import (
	"testing"

	"github.com/vito/booklit"
)

func TestSyntaxLinksBooklitComponentTags(t *testing.T) {
	section := &booklit.Section{}

	content, err := syntax(section, "markdown", booklit.String(`<IncludeSection path="quotes.md"/>`), true)
	if err != nil {
		t.Fatal(err)
	}

	// Block syntax wraps in <pre><code class="language-X">...</code></pre>.
	pre, ok := content.(booklit.RawElement)
	if !ok || pre.Tag != "pre" {
		t.Fatalf("expected outer <pre> RawElement, got %T (tag=%q)", content, pre.Tag)
	}
	code, ok := pre.Content.(booklit.RawElement)
	if !ok || code.Tag != "code" {
		t.Fatalf("expected inner <code> RawElement, got %T (tag=%q)", pre.Content, code.Tag)
	}
	seq, ok := code.Content.(booklit.Sequence)
	if !ok {
		t.Fatalf("expected sequence body, got %T", code.Content)
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
