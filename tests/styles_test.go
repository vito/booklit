package tests

import (
	"testing"
)

func TestStyles(t *testing.T) {
	for _, tt := range []struct {
		name    string
		example Example
	}{
		{
			name: "styled sections",
			example: Example{
				Input: `# Hello, world!

<Styled name="styled"/>

Sup?
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1 class="styled">Hello, world!</h1>

	<p>Sup?</p>
</section>
`,
				},
			},
		},
		{
			name: "styled pages",
			example: Example{
				Input: `# Hello, world!

<Styled name="full-styled"/>

Sup?
`,

				Outputs: Files{
					"hello-world.html": `<section class="full-styled-page">
	<h1 class="full-styled">Hello, world!</h1>

	<p>Sup?</p>
</section>
`,
				},
			},
		},
	} {
		t.Run(tt.name, tt.example.Run)
	}
}
