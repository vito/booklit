package tests

import (
	"testing"
)

// TestDangComponentDispatch covers tier-3 JSX→Dang dispatch: when a JSX
// element's name isn't a built-in, the evaluator looks it up as a Dang
// function in the project's auto-loaded .dang files.
func TestDangComponentDispatch(t *testing.T) {
	for _, tt := range []struct {
		name    string
		example Example
	}{
		{
			name: "body-less component returns a string",
			example: Example{
				Inputs: Files{
					"helpers.dang": `pub Greet(name: String!): String! {
  "hello, " + name + "!"
}
`,
				},
				Input: `<Title>Hi</Title>

Says: <Italic><Greet name="world"/></Italic>.
`,
				Outputs: Files{
					"hi.html": `<section>
	<h1>Hi</h1>

	<p>Says: <em>hello, world!</em>.</p>
</section>
`,
				},
			},
		},
		{
			name: "body component binds named args in children",
			example: Example{
				Inputs: Files{
					"helpers.dang": `pub Greet(name: String!, &body(greeting: String!): Boolean!): Boolean! {
  body(greeting: "hello, " + name + "!")
  true
}
`,
				},
				Input: `<Title>Hi</Title>

Says: <Italic><Greet name="world">{greeting}</Greet></Italic>.
`,
				Outputs: Files{
					"hi.html": `<section>
	<h1>Hi</h1>

	<p>Says: <em>hello, world!</em>.</p>
</section>
`,
				},
			},
		},
		{
			name: "iteration via repeated body calls",
			example: Example{
				Inputs: Files{
					"helpers.dang": `pub Each(items: [String!]!, &body(item: String!): Boolean!): Boolean! {
  items.each { item => body(item: item) }
  true
}
`,
				},
				Input: `<Title>List</Title>

Joined: <Italic><Each items={["a", "b", "c"]}>{item}</Each></Italic>.
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
			name: "missing function falls through to template fallback",
			example: Example{
				// No helpers.dang, no <Unknown> built-in → tier-3
				// misses, tier-4 (template) takes over. The fallback
				// emits a Styled{Style: "Unknown"} which the renderer
				// will try to look up; without a template, the render
				// fails. We assert on the render failure.
				Input: `<Title>Test</Title>

<Unknown>body</Unknown>
`,
				RenderErr: "Unknown",
			},
		},
		{
			name: "non-callable Dang binding falls through",
			example: Example{
				Inputs: Files{
					// pub Pal is a plain string, not callable.
					"helpers.dang": `pub Pal = "world"
`,
				},
				Input: `<Title>Test</Title>

<Pal>body</Pal>
`,
				// Falls through to template fallback; no Pal template
				// exists, so render fails.
				RenderErr: "Pal",
			},
		},
	} {
		t.Run(tt.name, tt.example.Run)
	}
}



// TestDangComponentRecordIteration covers the headline iteration
// pattern: a Dang function iterates a list-of-records and the JSX
// body uses record-field access on `item` to project each field into
// a child component or expression.
func TestDangComponentRecordIteration(t *testing.T) {
	example := Example{
		Inputs: Files{
			"helpers.dang": `pub Each(items: [a!]!, &body(item: a!): Boolean!): Boolean! {
  items.each { item => body(item: item) }
  true
}

pub primitiveTypes = [
  {{ name: "Int!", docs: "integer" }},
  {{ name: "String!", docs: "text" }},
]
`,
		},
		Input: `<Title>Types</Title>

Listing: <Each items={primitiveTypes}><Italic>{item.name}</Italic>=<Bold>{item.docs}</Bold> </Each>.
`,
		Outputs: Files{
			"types.html": `<section>
	<h1>Types</h1>

	<p>Listing: <em>Int!</em>=<strong>integer</strong> <em>String!</em>=<strong>text</strong> .</p>
</section>
`,
		},
	}
	t.Run("records iterate", example.Run)
}
