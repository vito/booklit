package tests

import (
	. "github.com/onsi/ginkgo/extensions/table"
)

var _ = DescribeTable("Blocks", (Example).Run,
	Entry("inline comments", Example{
		Input: `\title{Hello, world!}

{- Goodbye, -}cruel world!

Goodbye, {- cruel -}world!

Goodbye, cruel{- world! -}
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

<p>cruel world!</p>
<p>Goodbye, world!</p>
<p>Goodbye, cruel</p>
</section>`,
		},
	}),

	Entry("comments with delimitery things", Example{
		Input: `\title{Hello, world!}

Goodbye, {- {cru-el} -}world!
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

<p>Goodbye, world!</p>
</section>`,
		},
	}),

	Entry("nested comments", Example{
		Input: `\title{Hello, world!}

Goodbye, {- {cru-{- whoa -}el} -}world!
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

<p>Goodbye, world!</p>
</section>`,
		},
	}),

	Entry("block comments", Example{
		Input: `\title{Hello, world!}

Goodbye, world!

{-
	I'm a big old block comment.
-}

I'm another paragraph.

\section{
	\title{Subsection}

	Sup?

	{-
		I'm a big old indented block comment.
	-}

	Not much.
}
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

<p>Goodbye, world!</p>

<p>I'm another paragraph.</p>

<h2>1 Subsection</h2>

<p>Sup?</p>

<p>Not much.</p>
</section>`,
		},
	}),
)
