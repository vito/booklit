package tests

import (
	"strings"
	"testing"
)

func TestLineEndings(t *testing.T) {
	for _, tt := range []struct {
		name    string
		example Example
	}{
		{
			name: "simple Hello World",
			example: Example{
				Input: crlf(`# Hello, world!

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
			},
		},
		{
			name: "comments",
			example: Example{
				Input: crlf(`# Hello, world!

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
			},
		},
		{
			name: "verbatim Hello World",
			example: Example{
				Input: crlf("# Hello, world!\n\n" +
					"```\n" +
					"How are you?\n" +
					"```\n"),

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<pre>How are you?</pre>
</section>`,
				},
			},
		},
	} {
		t.Run(tt.name, tt.example.Run)
	}
}

func crlf(str string) string {
	return strings.ReplaceAll(str, "\n", "\r\n")
}
