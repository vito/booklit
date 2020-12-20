package booklit

// String is flow content containing arbitrary text.
//
// The text should not contain linebreaks.
type String string

// Empty is an empty String.
var Empty String

// IsFlow returns true.
func (str String) IsFlow() bool {
	return true
}

// String returns the string value.
func (str String) String() string {
	return string(str)
}

// Visit calls VisitString.
func (str String) Visit(visitor Visitor) error {
	return visitor.VisitString(str)
}
