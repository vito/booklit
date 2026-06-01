//go:build cgo

package treehighlight

import (
	"strings"
	"testing"
)

func TestHTMLHighlightsKnownLanguages(t *testing.T) {
	cases := []struct {
		language string
		source   string
		want     string
	}{
		{"lit", `\title{Hello}`, `color:#ed6c30`},
		{"go", `package main
func main() {}`, `color:#ed6c30`},
		{"jsx", `<Foo bar="x"/>`, `color:#fcc21b`},
		{"html", `<div class="x">hi</div>`, `&lt;`},
		{"bash", `echo "$HOME"`, `color:#ed6c30`},
	}

	for _, tc := range cases {
		t.Run(tc.language, func(t *testing.T) {
			html, err := HTML(tc.language, tc.source, false)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(html, tc.want) {
				t.Fatalf("highlighted HTML did not contain %q:\n%s", tc.want, html)
			}
			if strings.Contains(html, tc.source) && tc.language == "html" {
				t.Fatalf("HTML source was not escaped: %s", html)
			}
		})
	}
}

func TestChunksLinkBooklitCommands(t *testing.T) {
	chunks, err := Chunks("lit", `\title{Hello}`, Options{LinkReferences: true})
	if err != nil {
		t.Fatal(err)
	}

	assertLinkChunk(t, chunks, "title", "title")
}

func TestChunksLinkBooklitComponentTags(t *testing.T) {
	chunks, err := Chunks("markdown", `<IncludeSection path="quotes.md"/>`, Options{LinkReferences: true})
	if err != nil {
		t.Fatal(err)
	}

	assertLinkChunk(t, chunks, "include-section", "IncludeSection")
}

func TestBooklitLinkTag(t *testing.T) {
	cases := map[string]string{
		"IncludeSection":  "include-section",
		"OutputFrame":     "output-frame",
		"TableOfContents": "table-of-contents",
		"HTMLRenderer":    "html-renderer",
		"title":           "title",
		"include-section": "include-section",
	}

	for in, want := range cases {
		if got := booklitLinkTag(in); got != want {
			t.Fatalf("booklitLinkTag(%q) = %q, want %q", in, got, want)
		}
	}
}

func assertLinkChunk(t *testing.T, chunks []Chunk, tag, text string) {
	t.Helper()
	for _, chunk := range chunks {
		if chunk.LinkTag == tag && chunk.LinkText == text {
			return
		}
	}
	t.Fatalf("expected link chunk %q/%q, got %#v", tag, text, chunks)
}

func TestUnknownLanguageFallsBackToEscapedHTML(t *testing.T) {
	html, err := HTML("made-up", `<x>`, false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, `&lt;x&gt;`) {
		t.Fatalf("expected escaped fallback, got %s", html)
	}
}
