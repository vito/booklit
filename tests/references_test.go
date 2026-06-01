package tests

import (
	"testing"
)

func TestReferences(t *testing.T) {
	for _, tt := range []struct {
		name    string
		example Example
	}{
		{
			name: "references to other sections on the same page",
			example: Example{
				Input: `# Hello, world!

See also [the last section](#section-c).

## Section A

See also [#section-b].

## Section B

See also [#section-a].

## Section C

See also [#hello-world].
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>See also <a href="hello-world.html#section-c">the last section</a>.</p>

	<h2>1 Section A</h2>

	<p>See also <a href="hello-world.html#section-b">Section B</a>.</p>

	<h2>2 Section B</h2>

	<p>See also <a href="hello-world.html#section-a">Section A</a>.</p>

	<h2>3 Section C</h2>

	<p>See also <a href="hello-world.html">Hello, world!</a>.</p>
</section>
`,
				},
			},
		},
		{
			name: "references to other sections on split pages",
			example: Example{
				Input: `# Hello, world!

See also [the last section](#section-c).

<SplitSections/>

## Section A

See also [#section-b].

## Section B

See also [#section-a].

## Section C

See also [#hello-world].
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>See also <a href="section-c.html">the last section</a>.</p>
</section>
`,
					"section-a.html": `<section>
	<h1>1 Section A</h1>

	<p>See also <a href="section-b.html">Section B</a>.</p>
</section>
`,
					"section-b.html": `<section>
	<h1>2 Section B</h1>

	<p>See also <a href="section-a.html">Section A</a>.</p>
</section>
`,
					"section-c.html": `<section>
	<h1>3 Section C</h1>

	<p>See also <a href="hello-world.html">Hello, world!</a>.</p>
</section>
`,
				},
			},
		},
		{
			name: "explicit target elements",
			example: Example{
				Input: `# Hello, world!

Title fallback: [#target-a].

With display: [with display](#target-a).

Tag-name fallback: [#target-without-display].

With display: [with display](#target-without-display).

## Some Section

Foo bar.

<Target tag="target-a" title="Target A"/>

<Target tag="target-without-display"/>
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>Title fallback: <a href="hello-world.html#target-a">Target A</a>.</p>
	<p>With display: <a href="hello-world.html#target-a">with display</a>.</p>
	<p>Tag-name fallback: <a href="hello-world.html#target-without-display">target-without-display</a>.</p>
	<p>With display: <a href="hello-world.html#target-without-display">with display</a>.</p>

	<h2>1 Some Section</h2>

	<p>Foo bar.</p>

	<a id="target-a"></a>
	<a id="target-without-display"></a>
</section>
`,
				},
			},
		},
		{
			name: "aux",
			example: Example{
				Input: `# Hello, world!<Aux>: Foo Bar</Aux>

See also [the last section](#section-c).

<TableOfContents/>

## Section A

See also [#section-b].

## Section B<Aux>aby</Aux>

See also [#some-anchor].

## Section C

<Target tag="some-anchor">I'm an<Aux> awesome</Aux> anchor.</Target>

See also [#hello-world].
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!: Foo Bar</h1>

	<p>See also <a href="hello-world.html#section-c">the last section</a>.</p>

	<ul>
		<li><a href="hello-world.html#section-a">Section A</a></li>
		<li><a href="hello-world.html#section-b">Section B</a></li>
		<li><a href="hello-world.html#section-c">Section C</a></li>
	</ul>

	<h2>1 Section A</h2>

	<p>See also <a href="hello-world.html#section-b">Section B</a>.</p>

	<h2>2 Section Baby</h2>

	<p>See also <a href="hello-world.html#some-anchor">I'm an anchor.</a>.</p>

	<h2>3 Section C</h2>

	<a id="some-anchor"></a>

	<p>See also <a href="hello-world.html">Hello, world!</a>.</p>
</section>
`,
				},
			},
		},
	} {
		t.Run(tt.name, tt.example.Run)
	}
}
