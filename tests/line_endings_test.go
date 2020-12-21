package tests

import (
	"strings"

	. "github.com/onsi/ginkgo/extensions/table"
	_ "github.com/vito/booklit/tests/fixtures/stringer-plugin"
)

var _ = DescribeTable("Booklit", (Example).Run,
	Entry("simple 'Hello World'", Example{
		Input: crlf(`\title{Hello, world!}

How are you?
This is the same paragraph.

I'm another paragraph.
`),

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are you? This is the same paragraph.</p>

	<p>I'm another paragraph.</p>
</section>`,
		},
	}),
	Entry("comments", Example{
		Input: crlf(`\title{Hello, world!}

How are you?

{-
	This is the same paragraph.

	I'm another paragraph.
-}
`),

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are you?</p>
</section>`,
		},
	}),
	Entry("verbatim 'Hello World'", Example{
		Input: crlf(`\title{Hello, world!}

\code{{{
	How are you?
}}}
`),

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<pre>How are you?</pre>
</section>`,
		},
	}),
)

func crlf(str string) string {
	return strings.ReplaceAll(str, "\n", "\r\n")
}
