package tests

import (
	"testing"
)

func TestComments(t *testing.T) {
	for _, tt := range []struct {
		name    string
		example Example
	}{
		{
			name: "inline comments",
			example: Example{
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
			},
		},
		{
			name: "comments with delimitery things",
			example: Example{
				Input: `\title{Hello, world!}

Goodbye, {- {cru-el} -}world!
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

<p>Goodbye, world!</p>
</section>`,
				},
			},
		},
		{
			name: "nested comments",
			example: Example{
				Input: `\title{Hello, world!}

Goodbye, {- {cru-{- whoa -}el} -}world!
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

<p>Goodbye, world!</p>
</section>`,
				},
			},
		},
		{
			name: "block comments",
			example: Example{
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
			},
		},
	} {
		t.Run(tt.name, tt.example.Run)
	}
}
