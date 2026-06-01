package tests

import (
	"testing"
)

// TestMdxTemplateDispatch covers tier-4 JSX→mdx-template dispatch: when
// `<Foo/>` isn't a built-in and isn't a Dang function in scope, the
// evaluator looks up `Foo.md` and evaluates it with props bound in
// Dang scope and `children` carrying the JSX children's rendered
// content.
func TestMdxTemplateDispatch(t *testing.T) {
	for _, tt := range []struct {
		name    string
		example Example
	}{
		{
			name: "prop interpolation",
			example: Example{
				Inputs: Files{
					"components/Greet.md": `<RawHTML>hello, </RawHTML>{name}<RawHTML>!</RawHTML>`,
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
			name: "children interpolation via {children}",
			example: Example{
				Inputs: Files{
					"components/Wrap.md": `<Italic>{children}</Italic>`,
				},
				Input: `<Title>Hi</Title>

Wrapped: <Wrap>body text</Wrap>.
`,
				Outputs: Files{
					"hi.html": `<section>
	<h1>Hi</h1>

	<p>Wrapped: <em>body text</em>.</p>
</section>
`,
				},
			},
		},
		{
			name: "children interpolation via <Children/>",
			example: Example{
				Inputs: Files{
					"components/Wrap.md": `<Bold><Children/></Bold>`,
				},
				Input: `<Title>Hi</Title>

Wrapped: <Wrap>body text</Wrap>.
`,
				Outputs: Files{
					"hi.html": `<section>
	<h1>Hi</h1>

	<p>Wrapped: <strong>body text</strong>.</p>
</section>
`,
				},
			},
		},
		{
			name: "empty children renders nothing",
			example: Example{
				Inputs: Files{
					"components/Mark.md": `<Italic>[{children}]</Italic>`,
				},
				Input: `<Title>Hi</Title>

Empty: <Mark/>.
`,
				Outputs: Files{
					"hi.html": `<section>
	<h1>Hi</h1>

	<p>Empty: <em>[]</em>.</p>
</section>
`,
				},
			},
		},
		{
			name: "nested JSX inside template",
			example: Example{
				Inputs: Files{
					"components/Pair.md": `<Italic>{first}</Italic> and <Bold>{second}</Bold>`,
				},
				Input: `<Title>Hi</Title>

Picked: <Pair first="apples" second="oranges"/>.
`,
				Outputs: Files{
					"hi.html": `<section>
	<h1>Hi</h1>

	<p>Picked: <em>apples</em> and <strong>oranges</strong>.</p>
</section>
`,
				},
			},
		},
		{
			name: "multi-line template with raw HTML",
			example: Example{
				Inputs: Files{
					"components/Card.md": `<div class="card">
  <h3>{title}</h3>
  <div class="body">{children}</div>
</div>`,
				},
				Input: `<Title>Hi</Title>

<Card title="Welcome">
Body text here.
</Card>
`,
				Outputs: Files{
					"hi.html": `<section>
	<h1>Hi</h1>

	<div class="card">
		<h3>Welcome</h3>
		<div class="body"><p>Body text here.</p></div>
	</div>
</section>
`,
				},
			},
		},
		{
			// Per the React/MDX convention recorded in pivot.md design
			// notes: PascalCase JSX children are Markdown-parsed, lowercase
			// tags are raw HTML. Card wraps children in a lowercase <div>,
			// but Markdown is parsed at the *caller* level where <Card> is
			// invoked — `**bold**` is emphasized before reaching the
			// template.
			name: "Markdown inside Component children",
			example: Example{
				Inputs: Files{
					"components/Card.md": `<div class="card">{children}</div>`,
				},
				Input: `<Title>Hi</Title>

<Card>**bold** and *italic*.</Card>
`,
				Outputs: Files{
					"hi.html": `<section>
	<h1>Hi</h1>

	<div class="card"><strong>bold</strong> and <em>italic</em>.</div>
</section>
`,
				},
			},
		},
		{
			name: "templates have priority over legacy Styled fallback",
			example: Example{
				// Without Foo.md, this falls through to Styled fallback and
				// the renderer would error on a missing Foo.tmpl. With
				// Foo.md, the template handles it.
				Inputs: Files{
					"components/Foo.md": `<Italic>via template: {children}</Italic>`,
				},
				Input: `<Title>Hi</Title>

Picked: <Foo>x</Foo>.
`,
				Outputs: Files{
					"hi.html": `<section>
	<h1>Hi</h1>

	<p>Picked: <em>via template: x</em>.</p>
</section>
`,
				},
			},
		},
		{
			name: "Dang functions still beat templates",
			example: Example{
				// Both a Dang function AND a template for <Pick>. Per
				// phase-3b.md Q1, Dang wins because it's more specific.
				Inputs: Files{
					"helpers.dang": `pub Pick(): String! { "dang" }
`,
					"components/Pick.md": `<RawHTML>template</RawHTML>`,
				},
				Input: `<Title>Hi</Title>

Picked: <Italic><Pick/></Italic>.
`,
				Outputs: Files{
					"hi.html": `<section>
	<h1>Hi</h1>

	<p>Picked: <em>dang</em>.</p>
</section>
`,
				},
			},
		},
		{
			name: "expression prop passed to template",
			example: Example{
				Inputs: Files{
					"components/Echo.md": `<Italic>{val}</Italic>`,
				},
				Input: `<Title>Hi</Title>

Result: <Echo val={1 + 2}/>.
`,
				Outputs: Files{
					"hi.html": `<section>
	<h1>Hi</h1>

	<p>Result: <em>3</em>.</p>
</section>
`,
				},
			},
		},
	} {
		t.Run(tt.name, tt.example.Run)
	}
}
