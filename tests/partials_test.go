package tests

import (
	"testing"

	_ "github.com/vito/booklit/tests/testdata/set-partials-plugin"
)

func TestPartials(t *testing.T) {
	for _, tt := range []struct {
		name    string
		example Example
	}{
		{
			name: "set in section and rendered via template",
			example: Example{
				Input: `\title{Set Partial Read Template}

I want to be some body.

\set-partial{FooBar}{
	I'm a partial!
}

Some more body.
`,

				Outputs: Files{
					"set-partial-read-template.html": `<div>
	Here's a partial:

	<p>I’m a partial!</p>
</div>

<div>
	Here's the partial again:

	<p>I’m a partial!</p>
</div>

<p>I want to be some body.</p>

<p>Some more body.</p>
`,
				},
			},
		},
		{
			name: "targets and references are collected and resolved",
			example: Example{
				Input: `\title{Set Partial Read Template}

I want to be some body.

\set-partial{FooBar}{
	\target{some-target}{Hello.}

	I'm a partial!

	\reference{set-partial-read-template}
}

Some more body. \reference{some-target}
`,

				Outputs: Files{
					"set-partial-read-template.html": `<div>
	Here's a partial:

	<p><a id="some-target"></a></p>

	<p>I’m a partial!</p>

	<p><a href="set-partial-read-template.html">Set Partial Read Template</a></p>
</div>

<div>
	Here's the partial again:

	<p><a id="some-target"></a></p>

	<p>I’m a partial!</p>

	<p><a href="set-partial-read-template.html">Set Partial Read Template</a></p>
</div>

<p>I want to be some body.</p>

<p>Some more body. <a href="set-partial-read-template.html#some-target">Hello.</a></p>
`,
				},
			},
		},
		{
			name: "set in plugin and rendered in template",
			example: Example{
				Input: `\title{Set Partial Read Template}

\use-plugin{set-partials}

I want to be some body.

\set-the-partial

Some more body.
`,

				Outputs: Files{
					"set-partial-read-template.html": `<div>
	Here's a partial:

	<p>I'm a partial!</p>
</div>

<div>
	Here's the partial again:

	<p>I'm a partial!</p>
</div>

<p>I want to be some body.</p>

<p>Some more body.</p>
`,
				},
			},
		},
	} {
		t.Run(tt.name, tt.example.Run)
	}
}
