package tests

import (
	"testing"
)

func TestErrors(t *testing.T) {
	for _, tt := range []struct {
		name    string
		example Example
	}{
		{
			name: "unknown JSX component",
			example: Example{
				Input: `# Hello, world!

<Banana>attack</Banana>
`,

				LoadErr: `unknown JSX component <Banana>`,
			},
		},
		{
			name: "ambiguous references",
			example: Example{
				Input: `# Hello, world!

See also [this tag](#dupe-tag).

## First Tag

<Target tag="dupe-tag">I'm the first tag.</Target>

## Second Tag

<Target tag="dupe-tag">I'm the second tag.</Target>
`,

				RenderErr: `ambiguous target for tag 'dupe-tag'`,
			},
		},
		{
			name: "missing references",
			example: Example{
				Input: `# Hello, world!

See also [this tag](#missing-tag).
`,

				RenderErr: `unknown tag 'missing-tag'`,
			},
		},
		{
			name: "setting title twice",
			example: Example{
				Input: `# Hello, world!

# BAM
`,

				LoadErr: "cannot set title twice",
			},
		},
	} {
		t.Run(tt.name, tt.example.Run)
	}
}
