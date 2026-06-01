package tests

import (
	"testing"
)

// TestDangExpressions covers JSX {expr} interpolations. Each case wraps
// the expression inside a JSX element (the parser only triggers on `<`,
// so `{expr}` outside JSX is just literal text).
func TestDangExpressions(t *testing.T) {
	for _, tt := range []struct {
		name    string
		example Example
	}{
		{
			name: "int literal in child",
			example: Example{
				Input: `<Title>Hello</Title>

The answer is *{42}*.
`,
				Outputs: Files{
					"hello.html": `<section>
	<h1>Hello</h1>

	<p>The answer is <em>42</em>.</p>
</section>
`,
				},
			},
		},
		{
			name: "string literal in child",
			example: Example{
				Input: `<Title>Hello</Title>

Greeting: *{"hi"}*.
`,
				Outputs: Files{
					"hello.html": `<section>
	<h1>Hello</h1>

	<p>Greeting: <em>hi</em>.</p>
</section>
`,
				},
			},
		},
		{
			name: "arithmetic expression in child",
			example: Example{
				Input: `<Title>Hello</Title>

Result: *{1 + 2}*.
`,
				Outputs: Files{
					"hello.html": `<section>
	<h1>Hello</h1>

	<p>Result: <em>3</em>.</p>
</section>
`,
				},
			},
		},
		{
			name: "boolean in child",
			example: Example{
				Input: `<Title>Hello</Title>

Flag: *{true}*.
`,
				Outputs: Files{
					"hello.html": `<section>
	<h1>Hello</h1>

	<p>Flag: <em>true</em>.</p>
</section>
`,
				},
			},
		},
		{
			name: "expression in prop renders via template fallback",
			example: Example{
				Input: `<Title>Hello</Title>

<Card title={"Greetings"}>
Welcome.
</Card>
`,
				Outputs: Files{
					"hello.html": `<section>
	<h1>Hello</h1>

	<div class="card">
		<h3>Greetings</h3>
		<div class="body"><p>Welcome.</p></div>
	</div>
</section>
`,
				},
			},
		},
		{
			name: "list expression flattens into content",
			example: Example{
				Input: `<Title>Hello</Title>

Joined: *{["a", "b", "c"]}*.
`,
				Outputs: Files{
					"hello.html": `<section>
	<h1>Hello</h1>

	<p>Joined: <em>abc</em>.</p>
</section>
`,
				},
			},
		},
		{
			name: "parse error surfaces",
			example: Example{
				Input: `<Title>Hello</Title>

Bad: *{1 +}*.
`,
				LoadErr: "evaluating {1 +}",
			},
		},
	} {
		t.Run(tt.name, tt.example.Run)
	}
}
