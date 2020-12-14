package booklit

// String is flow content containing arbitrary text.
//
// The text should not contain linebreaks.
type String string

// Empty is an empty String.
var Empty String

func (str String) IsFlow() bool {
	return true
}

func (str String) String() string {
	return string(str)
}

func (str String) Visit(visitor Visitor) error {
	return visitor.VisitString(str)
}
