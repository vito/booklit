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
				Input: `\title{Hello, world!}

\styled{styled}

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
				Input: `\title{Hello, world!}

\styled{full-styled}

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
