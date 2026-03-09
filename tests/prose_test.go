package tests

import (
	. "github.com/onsi/ginkgo/v2"
	_ "github.com/vito/booklit/tests/fixtures/stringer-plugin"
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

	Entry("no trailing linebreak", Example{
		Input: `\title{Hello, world!}

How are you?`,

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

	Entry("images", Example{
		Input: `\title{Hello, world!}

Here's an \image{foo.png}{with alt text} and another \image{without.gif}.
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>Here's an <img src="foo.png" alt="with alt text" /> and another <img src="without.gif" alt="" />.</p>
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

	Entry("larger", Example{
		Input: `\title{Hello, world!}

How are \larger{you}?
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are <span style="font-size: 120%">you</span>?</p>
</section>`,
		},
	}),

	Entry("smaller", Example{
		Input: `\title{Hello, world!}

How are \smaller{you}?
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are <span style="font-size: 80%">you</span>?</p>
</section>`,
		},
	}),

	Entry("strike", Example{
		Input: `\title{Hello, world!}

How are \strike{you}?
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are <span style="text-decoration: line-through">you</span>?</p>
</section>`,
		},
	}),

	Entry("superscript", Example{
		Input: `\title{Hello, world!}

How are \superscript{you}?
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are <sup>you</sup>?</p>
</section>`,
		},
	}),

	Entry("subscript", Example{
		Input: `\title{Hello, world!}

How are \subscript{you}?
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are <sub>you</sub>?</p>
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

	Entry("inline code", Example{
		Input: `\title{Hello, world!}

This is some \code{inline} code.
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>This is some <code>inline</code> code.</p>
</section>
`,
		},
	}),

	Entry("empty arguments", Example{
		Input: `\title{Hello, world!}

This is an \italic{} empty italic.

This is a \italic{ } space italic.

This is an \italic{  } even more spaced italic.
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
	}),
)
