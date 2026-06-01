package tests

import (
	"testing"
)

func TestSections(t *testing.T) {
	for _, tt := range []struct {
		name    string
		example Example
	}{
		{
			name: "sub-sections",
			example: Example{
				Input: `# Hello, world!

How are you?

## How I'm doing

Good, thanks! And you?

## Their Reply

Good, thanks!
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are you?</p>

	<h2>1 How I'm doing</h2>

	<p>Good, thanks! And you?</p>

	<h2>2 Their Reply</h2>

	<p>Good, thanks!</p>
</section>
`,
				},
			},
		},
		{
			name: "sub-sections with explicit tags",
			example: Example{
				Input: `# Hello, world! {#hello}

How are you?

## How I'm doing {#how}

Good, thanks!
`,

				Outputs: Files{
					"hello.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are you?</p>

	<h2>1 How I'm doing</h2>

	<p>Good, thanks!</p>
</section>
`,
				},
			},
		},
		{
			name: "sub-sections from files",
			example: Example{
				Input: `# Hello, world!

How are you?

<IncludeSection path="how-im-doing.md"/>

## Their Reply

Good, thanks!
`,

				Inputs: Files{
					"how-im-doing.md": `# How I'm doing

Good, thanks! And you?
`,
				},

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are you?</p>

	<h2>1 How I'm doing</h2>

	<p>Good, thanks! And you?</p>

	<h2>2 Their Reply</h2>

	<p>Good, thanks!</p>
</section>
`,
				},
			},
		},
		{
			name: "including sections relative to the section's path",
			example: Example{
				Input: `# Hello, world!

How are you?

<IncludeSection path="./sub-path/how-im-doing.md"/>

## Their Reply

Good, thanks!
`,

				Inputs: Files{
					"sub-path/how-im-doing.md": `# How I'm doing

Good, thanks! And you?

<IncludeSection path="another-section.md"/>
`,

					"sub-path/another-section.md": `# My Response

Not bad, not bad.

## Including in an Inline Section

That's great.

<IncludeSection path="yet-another-section.md"/>
`,

					"sub-path/yet-another-section.md": `# Their Response to My Response

Sick.
`,
				},

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are you?</p>

	<h2>1 How I'm doing</h2>

	<p>Good, thanks! And you?</p>

	<h3>1.1 My Response</h3>

	<p>Not bad, not bad.</p>

	<h4>1.1.1 Including in an Inline Section</h4>

	<p>That's great.</p>

	<h5>1.1.1.1 Their Response to My Response</h5>

	<p>Sick.</p>

	<h2>2 Their Reply</h2>

	<p>Good, thanks!</p>
</section>
`,
				},
			},
		},
		{
			name: "nested sub-sections",
			example: Example{
				Input: `# Hello, world!

How are you?

## How I'm doing

### After Much Deliberation

I have decided that I'm doing well. How about you?

## Their Reply

Good, thanks!
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are you?</p>

	<h2>1 How I'm doing</h2>

	<h3>1.1 After Much Deliberation</h3>

	<p>I have decided that I'm doing well. How about you?</p>

	<h2>2 Their Reply</h2>

	<p>Good, thanks!</p>
</section>
`,
				},
			},
		},
		{
			name: "split sub-sections",
			example: Example{
				Input: `# Hello, world!

How are you?

<SplitSections/>

## How I'm Doing

Good, thanks!
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are you?</p>
</section>
`,
					"how-im-doing.html": `<section>
	<h1>1 How I'm Doing</h1>

	<p>Good, thanks!</p>
</section>
`,
				},
			},
		},
		{
			name: "forcing sections onto one page",
			example: Example{
				Input: `# Hello, world!

How are you? See <Reference tag="deep-inlined"/>.

<SinglePage/>
<SplitSections/>

## Section A

<SplitSections/>

Blah blah in section A.

### Nested Section

Good, thanks!

## Section B

Blah blah in section B.

### Nested Section 2

Foo bar.

### Nested Section 3

<SplitSections/>

Fizz buzz.

#### Super Duple Wrapped {#deep-inlined}

Whoooooa.
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are you? See <a href="hello-world.html#deep-inlined">Super Duple Wrapped</a>.</p>

	<h1>1 Section A</h1>

	<p>Blah blah in section A.</p>

	<h1>1.1 Nested Section</h1>

	<p>Good, thanks!</p>

	<h1>2 Section B</h1>

	<p>Blah blah in section B.</p>

	<h2>2.1 Nested Section 2</h2>

	<p>Foo bar.</p>

	<h2>2.2 Nested Section 3</h2>

	<p>Fizz buzz.</p>

	<h1>2.2.1 Super Duple Wrapped</h1>

	<p>Whoooooa.</p>
</section>
`,
				},
			},
		},
		{
			name: "split sub-sub-sections",
			example: Example{
				Input: `# Hello, world!

How are you?

## How I'm Doing

<SplitSections/>

Good, thanks!

### Nested Section

Sup.
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are you?</p>

	<h2>1 How I'm Doing</h2>

	<p>Good, thanks!</p>
</section>
`,
					"nested-section.html": `<section>
	<h1>1.1 Nested Section</h1>

	<p>Sup.</p>
</section>
`,
				},
			},
		},
		{
			name: "tables of contents",
			example: Example{
				Input: `# Hello, world!

How are you?

<TableOfContents/>

This is some more content.

## Top Section A

Foo bar.

### Nested Section

Fizz buzz.

### Another Nested Section

Fiddlesticks.

### Section with Omitted Children

I omit my children.

<OmitChildrenFromTableOfContents/>

#### Invisible Child

Boo!

## Top Section B

Fiddlesticks is as far as I go.
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are you?</p>

	<ul>
		<li>
			<a href="hello-world.html#top-section-a">Top Section A</a>
			<ul>
				<li><a href="hello-world.html#nested-section">Nested Section</a></li>
				<li><a href="hello-world.html#another-nested-section">Another Nested Section</a></li>
				<li><a href="hello-world.html#section-with-omitted-children">Section with Omitted Children</a></li>
			</ul>
		</li>
		<li>
			<a href="hello-world.html#top-section-b">Top Section B</a>
		</li>
	</ul>

	<p>This is some more content.</p>

	<h2>1 Top Section A</h2>

	<p>Foo bar.</p>

	<h3>1.1 Nested Section</h3>

	<p>Fizz buzz.</p>

	<h3>1.2 Another Nested Section</h3>

	<p>Fiddlesticks.</p>

	<h3>1.3 Section with Omitted Children</h3>

	<p>I omit my children.</p>

	<h4>1.3.1 Invisible Child</h4>

	<p>Boo!</p>

	<h2>2 Top Section B</h2>

	<p>Fiddlesticks is as far as I go.</p>
</section>
`,
				},
			},
		},
		{
			name: "styled sections",
			example: Example{
				Input: `# Hello, world!

<Styled name="top-template"/>

How are you?

## How I'm doing

<Styled name="sub-template"/>

Good, thanks! And you?

### After Much Deliberation

I have decided that I'm doing well. How about you?

## Their Reply

Good, thanks!
`,

				Outputs: Files{
					"hello-world.html": `<section class="custom-top-page">
	<h1>Hello, world!</h1>

	<p>I'm a toplevel template! Here's my body:</p>

	<div class="custom-top-body">
		<p>How are you?</p>
	</div>

	<h2>1 How I'm doing</h2>

	<p>I'm a sub template! Here's my body:</p>

	<div class="custom-sub-body">
		<p>Good, thanks! And you?</p>
	</div>

	<h3>1.1 After Much Deliberation</h3>

	<p>I have decided that I'm doing well. How about you?</p>

	<h2>2 Their Reply</h2>

	<p>Good, thanks!</p>
</section>
`,
				},
			},
		},
	} {
		t.Run(tt.name, tt.example.Run)
	}
}
