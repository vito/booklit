package tests

import (
	. "github.com/onsi/ginkgo/extensions/table"
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

	Entry("invokes interspersed in words", Example{
		Input: `\title{Hello, world!}

This{\italic{is}}a test.
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>This<em>is</em>a test.</p>
</section>
`,
		},
	}),

	Entry("word-wrapped lines", Example{
		Input: `\title{Hello, world!}

Lorem ipsum dolor sit amet, consectetur adipiscing elit. \italic{Curabitur
accumsan a ligula id feugiat. Quisque luctus semper ex sodales vulputate.} Sed
mi mi, rhoncus non justo et, aliquam dictum est. Donec egestas massa id
pharetra scelerisque. Nulla nunc quam, sagittis vel est sed, ultrices bibendum
magna. Nulla posuere ut erat eget tristique. Nullam vel nisl vitae dui
sollicitudin porta.

\section{
	\title{Indented}

	Integer malesuada purus dignissim turpis lacinia fringilla. Suspendisse
	potenti. Maecenas varius iaculis volutpat. \italic{Vestibulum sagittis lacus
	ut ex varius molestie.}
}
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

	Entry("code block indent tracking", Example{
		Input: `\title{Hello, world!}

\code{{
I'm a code block.
}}

\section{
  \title{Sub-section}

	\code{{
	I'm a code block in a sub-section.
  }}

  \code{{
    {}   I'm a code block {- in a sub-section -}with a forced indent level.
  }}
}
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<pre>I'm a code block.</pre>

	<h2>1 Sub-section</h2>

	<pre>I'm a code block in a sub-section.</pre>

	<pre>   I'm a code block with a forced indent level.</pre>
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

	Entry("preformatted string arguments", Example{
		Input: `\title{Hello, world!}

\use-plugin{stringer}

Here's a code block:

\string{{{
	I'm a code block.

		I'm indented more.

			I'm indented even more.

I'm indented less.

	\reference{hello-world}

	\\some-method\{Some argument.\}


	One more line, with meaning.
}}}
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>Here's a code block:</p>

	<p>I'm a code block.

	I'm indented more.

		I'm indented even more.

I'm indented less.

\reference{hello-world}

\\some-method\{Some argument.\}


One more line, with meaning.
</p>
</section>
`,
		},
	}),
)
