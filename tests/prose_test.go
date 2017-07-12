package tests

import (
	. "github.com/onsi/ginkgo/extensions/table"
)

var _ = DescribeTable("Booklit", (Example).Run,
	Entry("simple 'Hello World'", Example{
		Input: `\title{Hello, world!}

How are you?
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are you?</p>
</section>`,
		},
	}),

	Entry("link", Example{
		Input: `\title{Hello, world!}

How are \link{you}{https://example.com}?
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are <a href="https://example.com">you</a>?</p>
</section>`,
		},
	}),

	Entry("italics", Example{
		Input: `\title{Hello, world!}

How are \italic{you}?
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are <em>you</em>?</p>
</section>`,
		},
	}),

	Entry("bold", Example{
		Input: `\title{Hello, world!}

How are \bold{you}?
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are <strong>you</strong>?</p>
</section>`,
		},
	}),

	Entry("multiple paragraphs", Example{
		Input: `\title{Hello, world!}

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
	}),

	Entry("word-wrapped paragraphs", Example{
		Input: `\title{Hello, world!}

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
	}),

	Entry("inline code and code blocks", Example{
		Input: `\title{Hello, world!}

This is some \code{inline} code.

Here's a code block:

\code{
	I'm a code block.

		I'm indented more.

			I'm indented even more.

I'm indented less.

	\reference{hello-world}

	\\some-method\{Some argument.\}


	One more line, with meaning.
}

\code{{
	I'm a code block.

		I'm indented more.

			I'm indented even more.

I'm indented less.

	\reference{hello-world}

	\\some-method\{Some argument.\}


	One more line, with meaning.
}}

\code{{{
	I'm a code block.

		I'm indented more.

			I'm indented even more.

I'm indented less.

	\reference{hello-world}

	\\some-method\{Some argument.\}


	One more line, with meaning.
}}}

And here's some more content.
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>This is some <code>inline</code> code.</p>

	<p>Here's a code block:</p>

	<pre><p>I'm a code block.</p><p>I'm indented more.</p><p>I'm indented even more.</p><p>I'm indented less.</p><p><a href="hello-world.html">Hello, world!</a></p><p>\some-method{Some argument.}</p><p>One more line, with meaning.</p></pre>

	<pre>I'm a code block.

	I'm indented more.

		I'm indented even more.

I'm indented less.

<a href="hello-world.html">Hello, world!</a>

\some-method{Some argument.}


One more line, with meaning.</pre>

	<pre>I'm a code block.

	I'm indented more.

		I'm indented even more.

I'm indented less.

\reference{hello-world}

\\some-method\{Some argument.\}


One more line, with meaning.</pre>

	<p>And here's some more content.</p>
</section>
`,
		},
	}),
)
