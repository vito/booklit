package tests

import (
	"testing"

	_ "github.com/vito/booklit/tests/fixtures/arbitrary-style-plugin"
)

func TestBlocks(t *testing.T) {
	for _, tt := range []struct {
		name    string
		example Example
	}{
		{
			name: "lists",
			example: Example{
				Input: `\title{Hello, world!}

\list{a}{
	b
}{
	\code{{
	c
	}}
}
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

<ul><li>a</li><li><p>b</p></li><li><pre>c</pre></li></ul>
</section>`,
				},
			},
		},
		{
			name: "tables",
			example: Example{
				Input: `\title{Hello, world!}

\table{
	\table-row{a}{1}
}{
	\table-row{b}{2}
}
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

<table>
	<tr>
		<td>a</td>
		<td>1</td>
	</tr>
	<tr>
		<td>b</td>
		<td>2</td>
	</tr>
</table>
</section>`,
				},
			},
		},
		{
			name: "definitions",
			example: Example{
				Input: `\title{Hello, world!}

\definitions{
	\definition{a}{1}
}{
	\definition{b}{2}
}
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

<dl>
	<dt>a</dt>
		<dd>1</dd>

	<dt>b</dt>
		<dd>2</dd>
</dl>
</section>`,
				},
			},
		},
		{
			name: "inset",
			example: Example{
				Input: `\title{Hello, world!}

\inset{
	Hello.
}
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

<div style="margin: 0 2em 1em" class="inset">
	<p>Hello.</p>
</div>
</section>`,
				},
			},
		},
		{
			name: "aside",
			example: Example{
				Input: `\title{Hello, world!}

\aside{
	Hello.
}
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

<blockquote class="aside">
	<p>Hello.</p>
</blockquote>
</section>`,
				},
			},
		},
		{
			name: "ordered lists",
			example: Example{
				Input: `\title{Hello, world!}

\ordered-list{a}{
	b
}{
	\code{{
	c
	}}
}
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

<ol><li>a</li><li><p>b</p></li><li><pre>c</pre></li></ol>
</section>`,
				},
			},
		},
		{
			name: "arbitrary styles",
			example: Example{
				Input: `\title{Hello, world!}

\use-plugin{arbitrary-style}

\arbitrary-style{
	Sup!
}
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

<blink><p>Sup!</p></blink>
</section>`,
				},
			},
		},
		{
			name: "markdown tables",
			example: Example{
				Input: `\title{Hello, world!}

| a | 1 |
| --- | --- |
| b | 2 |
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

<table>
	<tr>
		<td>a</td>
		<td>1</td>
	</tr>
	<tr>
		<td>b</td>
		<td>2</td>
	</tr>
</table>
</section>`,
				},
			},
		},
	} {
		t.Run(tt.name, tt.example.Run)
	}
}
