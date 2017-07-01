package tests

import (
	. "github.com/onsi/ginkgo/extensions/table"
)

var _ = DescribeTable("Booklit", (Example).Run,
	Entry("simple 'Hello World'", Example{
		Input: `\title{Hello, world!}

How are you?
`,

		Outputs: Outputs{
			"hello-world.html": `<h1>Hello, world!</h1>

<p>How are you?</p>`,
		},
	}),

	Entry("multiple paragraphs", Example{
		Input: `\title{Hello, world!}

How are you?

I'm good, thanks!
`,

		Outputs: Outputs{
			"hello-world.html": `<h1>Hello, world!</h1>

<p>How are you?</p>

<p>I'm good, thanks!</p>
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

		Outputs: Outputs{
			"hello-world.html": `<h1>Hello, world!</h1>

<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur accumsan a ligula id feugiat. Quisque luctus semper ex sodales vulputate. Sed mi mi, rhoncus non justo et, aliquam dictum est. Donec egestas massa id pharetra scelerisque. Nulla nunc quam, sagittis vel est sed, ultrices bibendum magna. Nulla posuere ut erat eget tristique. Nullam vel nisl vitae dui sollicitudin porta.</p>

<p>Integer malesuada purus dignissim turpis lacinia fringilla. Suspendisse potenti. Maecenas varius iaculis volutpat. Vestibulum sagittis lacus ut ex varius molestie.</p>
`,
		},
	}),

	Entry("sub-sections", Example{
		Input: `\title{Hello, world!}

How are you?

\section{
	\title{How I'm doing}

	Good, thanks! And you?
}

\section{
	\title{Their Reply}

	Good, thanks!
}
`,

		Outputs: Outputs{
			"hello-world.html": `<h1>Hello, world!</h1>

<p>How are you?</p>

<h2>How I'm doing</h2>

<p>Good, thanks! And you?</p>

<h2>Their Reply</h2>

<p>Good, thanks!</p>
`,
		},
	}),

	Entry("split sub-sections", Example{
		Input: `\title{Hello, world!}

How are you?

\split-sections

\section{
	\title{How I'm doing}

	Good, thanks!
}
`,

		Outputs: Outputs{
			"hello-world.html": `<h1>Hello, world!</h1>

<p>How are you?</p>
`,
			"how-im-doing.html": `<h1>How I'm doing</h1>

<p>Good, thanks!</p>
`,
		},
	}),

	Entry("references to other sections on the same page", Example{
		Input: `\title{Hello, world!}

See also \reference{section-c}{the last section}.

\section{
	\title{Section A}

	See also \reference{section-b}.
}

\section{
	\title{Section B}

	See also \reference{section-a}.
}

\section{
	\title{Section C}

	See also \reference{hello-world}.
}
`,

		Outputs: Outputs{
			"hello-world.html": `<h1>Hello, world!</h1>

<p>See also <a href="hello-world.html#section-c">the last section</a>.</p>

<h2>Section A</h2>

<p>See also <a href="hello-world.html#section-b">Section B</a>.</p>

<h2>Section B</h2>

<p>See also <a href="hello-world.html#section-a">Section A</a>.</p>

<h2>Section C</h2>

<p>See also <a href="hello-world.html">Hello, world!</a>.</p>
`,
		},
	}),

	Entry("references to other sections on split pages", Example{
		Input: `\title{Hello, world!}

See also \reference{section-c}{the last section}.

\split-sections

\section{
	\title{Section A}

	See also \reference{section-b}.
}

\section{
	\title{Section B}

	See also \reference{section-a}.
}

\section{
	\title{Section C}

	See also \reference{hello-world}.
}
`,

		Outputs: Outputs{
			"hello-world.html": `<h1>Hello, world!</h1>

<p>See also <a href="section-c.html">the last section</a>.</p>
`,
			"section-a.html": `<h1>Section A</h1>

<p>See also <a href="section-b.html">Section B</a>.</p>
`,
			"section-b.html": `<h1>Section B</h1>

<p>See also <a href="section-a.html">Section A</a>.</p>
`,
			"section-c.html": `<h1>Section C</h1>

<p>See also <a href="hello-world.html">Hello, world!</a>.</p>
`,
		},
	}),
)
