package tests

import "testing"

// TestControlFlow covers the <For>, <If>, and <Unless> built-ins.
// These wrap the body-block + WithBindings pattern from tier-3 Dang
// dispatch so authors don't need to define per-project Each / If
// Dang helpers for the common iteration and conditional cases.
func TestControlFlow(t *testing.T) {
	for _, tt := range []struct {
		name    string
		example Example
	}{
		{
			name: "For iterates string list with default `item` binding",
			example: Example{
				Input: `<Title>List</Title>

Joined: *<For each={["a", "b", "c"]}>{item}</For>*.
`,
				Outputs: Files{
					"list.html": `<section>
	<h1>List</h1>

	<p>Joined: <em>abc</em>.</p>
</section>
`,
				},
			},
		},
		{
			name: "For iterates record list with custom binding name",
			example: Example{
				Inputs: Files{
					"helpers.dang": `pub items = [
  {{ name: "Int!", docs: "integer" }},
  {{ name: "String!", docs: "text" }},
]
`,
				},
				Input: `<Title>Types</Title>

Listing: <For each={items} as="t"><em>{t.name}</em>=<strong>{t.docs}</strong> </For>.
`,
				Outputs: Files{
					"types.html": `<section>
	<h1>Types</h1>

	<p>Listing: <em>Int!</em>=<strong>integer</strong> <em>String!</em>=<strong>text</strong> .</p>
</section>
`,
				},
			},
		},
		{
			name: "If renders children when cond is true",
			example: Example{
				Input: `<Title>Cond</Title>

Result: *<If cond={true}>yes</If><If cond={false}>no</If>*.
`,
				Outputs: Files{
					"cond.html": `<section>
	<h1>Cond</h1>

	<p>Result: <em>yes</em>.</p>
</section>
`,
				},
			},
		},
		{
			name: "Unless renders children when cond is false",
			example: Example{
				Input: `<Title>Cond</Title>

Result: *<Unless cond={false}>yes</Unless><Unless cond={true}>no</Unless>*.
`,
				Outputs: Files{
					"cond.html": `<section>
	<h1>Cond</h1>

	<p>Result: <em>yes</em>.</p>
</section>
`,
				},
			},
		},
		{
			name: "For with empty list emits nothing",
			example: Example{
				Input: `<Title>Empty</Title>

Result: *[<For each={[]}>x</For>]*.
`,
				Outputs: Files{
					"empty.html": `<section>
	<h1>Empty</h1>

	<p>Result: <em>[]</em>.</p>
</section>
`,
				},
			},
		},
	} {
		t.Run(tt.name, tt.example.Run)
	}
}
