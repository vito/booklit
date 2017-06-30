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

\section{How I'm doing}{
	Good, thanks!
}
`,

		Outputs: Outputs{
			"hello-world.html": `<h1>Hello, world!</h1>

<p>How are you?</p>

<h2>How I'm doing</h2>

<p>Good, thanks!</p>
`,
		},
	}),

	Entry("split sub-sections", Example{
		Input: `\title{Hello, world!}

How are you?

\split-sections

\section{How I'm doing}{
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
)
