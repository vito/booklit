package tests

import (
	"testing"
)

// TestJSX covers JSX dispatch end-to-end: parse, evaluate via builtins
// registry or template fallback, render. Each case mirrors what the
// equivalent `\foo{}` invocation produces so output equivalence is
// verifiable.
func TestJSX(t *testing.T) {
	for _, tt := range []struct {
		name    string
		example Example
	}{
		{
			name: "title via JSX",
			example: Example{
				Input: `<Title>Hello, world!</Title>

How are you?
`,
				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are you?</p>
</section>
`,
				},
			},
		},
		{
			name: "italic and bold inline",
			example: Example{
				Input: `<Title>Hello, world!</Title>

This is <Italic>emphasized</Italic> and this is <Bold>strong</Bold>.
`,
				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>This is <em>emphasized</em> and this is <strong>strong</strong>.</p>
</section>
`,
				},
			},
		},
		{
			name: "nested section with JSX title",
			example: Example{
				Input: `<Title>Hello, world!</Title>

<Section>
<Title>Sub-section</Title>

Sub body.
</Section>
`,
				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<h2>1 Sub-section</h2>

	<p>Sub body.</p>
</section>
`,
				},
			},
		},
		{
			name: "reference to a target",
			example: Example{
				Input: `<Title>Hello, world!</Title>

The <Target tag="middle">middle</Target> is here.

See <Reference tag="middle">it again</Reference>.
`,
				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>The <a id="middle"></a> is here.</p>

	<p>See <a href="hello-world.html#middle">it again</a>.</p>
</section>
`,
				},
			},
		},
		{
			name: "inline raw HTML passes through",
			example: Example{
				Input: `<Title>Hello, world!</Title>

This has <em>literal em</em> in it.
`,
				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>This has <em>literal em</em> in it.</p>
</section>
`,
				},
			},
		},
		{
			name: "block raw HTML passes through",
			example: Example{
				Input: `<Title>Hello, world!</Title>

<dl>
<dt><em>term</em></dt>
<dd>body</dd>
</dl>
`,
				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<dl>
		<dt><em>term</em></dt>
		<dd>body</dd>
	</dl>
</section>
`,
				},
			},
		},
		{
			name: "template fallback for unknown component",
			example: Example{
				Input: `<Title>Hello, world!</Title>

<Card title="Greetings">
Welcome to the test.
</Card>
`,
				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<div class="card">
		<h3>Greetings</h3>
		<div class="body"><p>Welcome to the test.</p></div>
	</div>
</section>
`,
				},
			},
		},
	} {
		t.Run(tt.name, tt.example.Run)
	}
}
