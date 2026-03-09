package tests

import (
	. "github.com/onsi/ginkgo/v2"
)

var _ = DescribeTable("Comments", (Example).Run,
	Entry("inline comment", Example{
		Input: `\title{Hello, world!}

How {- secret -}are you?
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are you?</p>
</section>
`,
		},
	}),

	Entry("block comment", Example{
		Input: `\title{Hello, world!}

First paragraph.

{- this whole paragraph
is commented out -}

Second paragraph.
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>First paragraph.</p>

	<p>Second paragraph.</p>
</section>
`,
		},
	}),

	Entry("nested comment", Example{
		Input: `\title{Hello, world!}

Hello {- outer {- inner -} still gone -}world.
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>Hello world.</p>
</section>
`,
		},
	}),
)
