package tests

import (
	"testing"
)

func TestBlocks(t *testing.T) {
	for _, tt := range []struct {
		name    string
		example Example
	}{
		{
			name: "lists",
			example: Example{
				Input: "# Hello, world!\n\n" +
					"- a\n" +
					"- b\n" +
					"- ```\n" +
					"  c\n" +
					"  ```\n",

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

<ul><li>a</li><li>b</li><li><pre>c</pre></li></ul>
</section>`,
				},
			},
		},
		{
			name: "definitions",
			example: Example{
				Input: `# Hello, world!

<Definitions>
<Definition term="a">1</Definition>
<Definition term="b">2</Definition>
</Definitions>
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
				Input: `# Hello, world!

> Hello.
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

<div class="inset" style="margin: 0 2em 1em">
	<p>Hello.</p>
</div>
</section>`,
				},
			},
		},
		{
			name: "aside",
			example: Example{
				Input: `# Hello, world!

<Aside>
Hello.
</Aside>
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
				Input: "# Hello, world!\n\n" +
					"1. a\n" +
					"2. b\n" +
					"3. ```\n" +
					"   c\n" +
					"   ```\n",

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

<ol><li>a</li><li>b</li><li><pre>c</pre></li></ol>
</section>`,
				},
			},
		},
		{
			name: "markdown tables",
			example: Example{
				Input: `# Hello, world!

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
