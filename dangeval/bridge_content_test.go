package dangeval

import (
	"reflect"
	"testing"

	"github.com/vito/booklit"
	"github.com/vito/dang/pkg/dang"
)

// jsonScalar builds a Dang `JSON` scalar value carrying raw, as a Dagger
// module would return when handing back serialized Booklit content.
func jsonScalar(raw string) dang.ScalarValue {
	return dang.ScalarValue{Val: raw, ScalarType: &dang.Module{Named: "JSON"}}
}

func TestContentFromValueJSONScalarDecodes(t *testing.T) {
	e := &Evaluator{}
	got, err := e.ContentFromValue(
		jsonScalar(`{"k":"element","htag":"span","attrs":" class=\"x\"","content":{"k":"string","s":"hi"}}`),
		&booklit.Section{},
	)
	if err != nil {
		t.Fatalf("ContentFromValue: %v", err)
	}
	want := booklit.RawElement{Tag: "span", Attrs: ` class="x"`, Content: booklit.String("hi")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func TestContentFromValueJSONScalarRehydratesReference(t *testing.T) {
	e := &Evaluator{}
	sec := &booklit.Section{}
	got, err := e.ContentFromValue(jsonScalar(`{"k":"ref","tag":"x","optional":true}`), sec)
	if err != nil {
		t.Fatalf("ContentFromValue: %v", err)
	}
	ref, ok := got.(*booklit.Reference)
	if !ok {
		t.Fatalf("expected *booklit.Reference, got %T", got)
	}
	if ref.Section != sec {
		t.Errorf("reference not bound to section: got %p, want %p", ref.Section, sec)
	}
	if ref.TagName != "x" {
		t.Errorf("tag = %q, want %q", ref.TagName, "x")
	}
}

func TestContentFromValuePlainStringStaysText(t *testing.T) {
	e := &Evaluator{}
	got, err := e.ContentFromValue(dang.StringValue{Val: "hello"}, nil)
	if err != nil {
		t.Fatalf("ContentFromValue: %v", err)
	}
	if got != booklit.String("hello") {
		t.Errorf("got %#v, want booklit.String(%q)", got, "hello")
	}
}

func TestContentFromValueNonJSONScalarStaysText(t *testing.T) {
	// A custom scalar that is NOT JSON must render as its text, not be decoded
	// as content.
	e := &Evaluator{}
	val := dang.ScalarValue{Val: "raw", ScalarType: &dang.Module{Named: "CustomScalar"}}
	got, err := e.ContentFromValue(val, nil)
	if err != nil {
		t.Fatalf("ContentFromValue: %v", err)
	}
	if got != booklit.String("raw") {
		t.Errorf("got %#v, want booklit.String(%q)", got, "raw")
	}
}

func TestContentFromValueListDecodesElements(t *testing.T) {
	// A list mixing plain text and a content-JSON scalar bridges each element.
	e := &Evaluator{}
	got, err := e.ContentFromValue(dang.ListValue{Elements: []dang.Value{
		dang.StringValue{Val: "a"},
		jsonScalar(`{"k":"string","s":"b"}`),
	}}, nil)
	if err != nil {
		t.Fatalf("ContentFromValue: %v", err)
	}
	want := booklit.Sequence{booklit.String("a"), booklit.String("b")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v", got, want)
	}
}
