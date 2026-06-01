package tests

import (
	"testing"
)

func TestProse(t *testing.T) {
	for _, tt := range []struct {
		name    string
		example Example
	}{
		{
			name: "simple Hello World",
			example: Example{
				Input: `# Hello, world!

How are you?
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are you?</p>
</section>`,
				},
			},
		},
		{
			name: "no trailing linebreak",
			example: Example{
				Input: `# Hello, world!

How are you?`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are you?</p>
</section>`,
				},
			},
		},
		{
			name: "link",
			example: Example{
				Input: `# Hello, world!

How are [you](https://example.com)?
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are <a href="https://example.com">you</a>?</p>
</section>`,
				},
			},
		},
		{
			name: "images",
			example: Example{
				Input: `# Hello, world!

Here's an ![with alt text](foo.png) and another ![](without.gif).
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>Here's an <img src="foo.png" alt="with alt text" /> and another <img src="without.gif" alt="" />.</p>
</section>`,
				},
			},
		},
		{
			name: "italics",
			example: Example{
				Input: `# Hello, world!

How are *you*?
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are <em>you</em>?</p>
</section>`,
				},
			},
		},
		{
			name: "bold",
			example: Example{
				Input: `# Hello, world!

How are **you**?
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are <strong>you</strong>?</p>
</section>`,
				},
			},
		},
		{
			name: "larger",
			example: Example{
				Input: `# Hello, world!

How are <Larger>you</Larger>?
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are <span style="font-size: 120%">you</span>?</p>
</section>`,
				},
			},
		},
		{
			name: "smaller",
			example: Example{
				Input: `# Hello, world!

How are <Smaller>you</Smaller>?
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are <span style="font-size: 80%">you</span>?</p>
</section>`,
				},
			},
		},
		{
			name: "strike",
			example: Example{
				Input: `# Hello, world!

How are <Strike>you</Strike>?
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are <span style="text-decoration: line-through">you</span>?</p>
</section>`,
				},
			},
		},
		{
			name: "superscript",
			example: Example{
				Input: `# Hello, world!

How are <sup>you</sup>?
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are <sup>you</sup>?</p>
</section>`,
				},
			},
		},
		{
			name: "subscript",
			example: Example{
				Input: `# Hello, world!

How are <sub>you</sub>?
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are <sub>you</sub>?</p>
</section>`,
				},
			},
		},
		{
			name: "multiple paragraphs",
			example: Example{
				Input: `# Hello, world!

How are you?

I'm good, thanks!
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are you?</p>

	<p>I'm good, thanks!</p>
</section>
`,
				},
			},
		},
		{
			name: "word-wrapped lines",
			example: Example{
				Input: `# Hello, world!

Lorem ipsum dolor sit amet, consectetur adipiscing elit. *Curabitur
accumsan a ligula id feugiat. Quisque luctus semper ex sodales vulputate.* Sed
mi mi, rhoncus non justo et, aliquam dictum est. Donec egestas massa id
pharetra scelerisque. Nulla nunc quam, sagittis vel est sed, ultrices bibendum
magna. Nulla posuere ut erat eget tristique. Nullam vel nisl vitae dui
sollicitudin porta.

## Indented

Integer malesuada purus dignissim turpis lacinia fringilla. Suspendisse
potenti. Maecenas varius iaculis volutpat. *Vestibulum sagittis lacus
ut ex varius molestie.*
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. <em>Curabitur accumsan a ligula id feugiat. Quisque luctus semper ex sodales vulputate.</em> Sed mi mi, rhoncus non justo et, aliquam dictum est. Donec egestas massa id pharetra scelerisque. Nulla nunc quam, sagittis vel est sed, ultrices bibendum magna. Nulla posuere ut erat eget tristique. Nullam vel nisl vitae dui sollicitudin porta.</p>

	<h2>1 Indented</h2>

	<p>Integer malesuada purus dignissim turpis lacinia fringilla. Suspendisse potenti. Maecenas varius iaculis volutpat. <em>Vestibulum sagittis lacus ut ex varius molestie.</em></p>
</section>
`,
				},
			},
		},
		{
			name: "word-wrapped paragraphs",
			example: Example{
				Input: `# Hello, world!

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur accumsan a
ligula id feugiat. Quisque luctus semper ex sodales vulputate. Sed mi mi,
rhoncus non justo et, aliquam dictum est. Donec egestas massa id pharetra
scelerisque. Nulla nunc quam, sagittis vel est sed, ultrices bibendum magna.
Nulla posuere ut erat eget tristique. Nullam vel nisl vitae dui sollicitudin
porta.

Integer malesuada purus dignissim turpis lacinia fringilla. Suspendisse
potenti. Maecenas varius iaculis volutpat. Vestibulum sagittis lacus ut ex
varius molestie.
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur accumsan a ligula id feugiat. Quisque luctus semper ex sodales vulputate. Sed mi mi, rhoncus non justo et, aliquam dictum est. Donec egestas massa id pharetra scelerisque. Nulla nunc quam, sagittis vel est sed, ultrices bibendum magna. Nulla posuere ut erat eget tristique. Nullam vel nisl vitae dui sollicitudin porta.</p>

	<p>Integer malesuada purus dignissim turpis lacinia fringilla. Suspendisse potenti. Maecenas varius iaculis volutpat. Vestibulum sagittis lacus ut ex varius molestie.</p>
</section>
`,
				},
			},
		},
		{
			name: "empty styled arguments",
			example: Example{
				Input: `# Hello, world!

This is an <em></em> empty italic.

This is a <em> </em> space italic.

This is an <em>  </em> even more spaced italic.
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>This is an <em></em> empty italic.</p>
	<p>This is a <em> </em> space italic.</p>
	<p>This is an <em>  </em> even more spaced italic.</p>
</section>
`,
				},
			},
		},
		{
			name: "fenced code blocks",
			example: Example{
				Input: "# Hello, world!\n\n" +
					"```\nHow are you?\n```\n",

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<pre>How are you?</pre>
</section>
`,
				},
			},
		},
		{
			name: "inline JSX interspersed in words",
			example: Example{
				Input: `# Hello, world!

This<em>is</em>a test.
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>This<em>is</em>a test.</p>
</section>
`,
				},
			},
		},
	} {
		t.Run(tt.name, tt.example.Run)
	}
}
