package tests

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
)

var _ = DescribeTable("Booklit", (Example).Run,
	Entry("CRLF line endings", Example{
		Input: crlf(`@title{Hello, world!}

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
)

func crlf(str string) string {
	return strings.ReplaceAll(str, "\n", "\r\n")
}
