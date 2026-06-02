package contentjson_test

import (
	"reflect"
	"testing"

	"github.com/vito/booklit"
	"github.com/vito/booklit/contentjson"
)

func TestRoundTrip(t *testing.T) {
	sec := &booklit.Section{}

	// One tree exercising every serializable content type, nested.
	content := booklit.Sequence{
		booklit.String("hello "),
		booklit.RawElement{Tag: "div", Attrs: ` class="x"`, Content: booklit.String("world")},
		booklit.RawElement{Tag: "img", Attrs: ` alt="alt text" src="/a.png"`},
		booklit.RawFragment{HTML: "<i>fragment</i>"},
		booklit.Aux{Content: booklit.String("aux")},
		&booklit.Reference{Section: sec, TagName: "some-tag", Content: booklit.String("display"), Optional: true},
		booklit.Target{TagName: "anchor", Title: booklit.String("title"), Content: booklit.String("body")},
		booklit.Paragraph{booklit.String("a sentence")},
	}

	data, err := contentjson.Marshal(content)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	got, err := contentjson.Unmarshal(data, sec)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if !reflect.DeepEqual(got, content) {
		t.Errorf("round trip mismatch:\n got: %#v\nwant: %#v", got, content)
	}
}

func TestReferenceRehydratesSection(t *testing.T) {
	// A reference serialized without a section must come back bound to the
	// section passed to Unmarshal — that is how cross-references resolve.
	ref := &booklit.Reference{TagName: "x", Optional: true}
	data, err := contentjson.Marshal(ref)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	sec := &booklit.Section{}
	got, err := contentjson.Unmarshal(data, sec)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	decoded, ok := got.(*booklit.Reference)
	if !ok {
		t.Fatalf("expected *booklit.Reference, got %T", got)
	}
	if decoded.Section != sec {
		t.Errorf("reference not rehydrated with section: got %p, want %p", decoded.Section, sec)
	}
	if decoded.TagName != "x" {
		t.Errorf("tag name = %q, want %q", decoded.TagName, "x")
	}
}

func TestMarshalRejectsInProcessContent(t *testing.T) {
	cases := []struct {
		name    string
		content booklit.Content
	}{
		{"section", &booklit.Section{}},
		{"table-of-contents", booklit.TableOfContents{}},
		{"lazy", &booklit.Lazy{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := contentjson.Marshal(tc.content); err == nil {
				t.Errorf("expected error marshaling %T, got nil", tc.content)
			}
		})
	}
}

func TestMarshalRejectsNestedInProcessContent(t *testing.T) {
	// In-process content nested inside an otherwise-serializable tree must
	// still be rejected, not silently dropped.
	content := booklit.Sequence{
		booklit.String("before"),
		booklit.TableOfContents{},
	}
	if _, err := contentjson.Marshal(content); err == nil {
		t.Error("expected error marshaling sequence containing TableOfContents, got nil")
	}
}

func TestNullAndEmpty(t *testing.T) {
	// Marshaling nil yields JSON null; Unmarshaling null yields Empty.
	data, err := contentjson.Marshal(nil)
	if err != nil {
		t.Fatalf("Marshal(nil): %v", err)
	}
	if string(data) != "null" {
		t.Errorf("Marshal(nil) = %q, want %q", data, "null")
	}

	got, err := contentjson.Unmarshal([]byte("null"), nil)
	if err != nil {
		t.Fatalf("Unmarshal(null): %v", err)
	}
	if got != booklit.Empty {
		t.Errorf("Unmarshal(null) = %#v, want booklit.Empty", got)
	}
}
