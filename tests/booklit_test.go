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
)
