package marklit

import "testing"

func TestSmartTypography(t *testing.T) {
	for _, tt := range []struct {
		name string
		prev rune
		in   string
		want string
	}{
		{"double quotes", 0, `say "hi" now`, "say “hi” now"},
		{"apostrophe", 0, "don't", "don’t"},
		{"single quotes", 0, "a 'quote' here", "a ‘quote’ here"},
		{"possessive plural", 0, "the dogs' bones", "the dogs’ bones"},
		{"leading quote opens", 0, `"x"`, "“x”"},
		{"escaped double stays straight", 0, `\"x\"`, `"x"`},
		{"escaped single stays straight", 0, `\'x\'`, `'x'`},
		{"prev letter closes double", 'd', `" rest`, "” rest"},
		{"prev letter apostrophe", 's', "' rest", "’ rest"},
		{"prev space opens double", ' ', `"x`, "“x"},
		{"en dash", 0, "a -- b", "a – b"},
		{"em dash", 0, "a --- b", "a — b"},
		{"ellipsis", 0, "wait...", "wait…"},
		{"escaped dash stays straight", 0, `a \-- b`, "a -- b"},
		{"escaped dot stays straight", 0, `a\.\.\. b`, "a... b"},
		{"single dash and dot untouched", 0, "a-b. end", "a-b. end"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := smartTypography(tt.prev, tt.in)
			if got != tt.want {
				t.Errorf("smartTypography(%q, %q) = %q, want %q", tt.prev, tt.in, got, tt.want)
			}
		})
	}
}
