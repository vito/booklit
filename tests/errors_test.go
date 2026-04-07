package tests

import (
	"testing"

	_ "github.com/vito/booklit/tests/fixtures/erroring-plugin"
)

func TestErrors(t *testing.T) {
	for _, tt := range []struct {
		name    string
		example Example
	}{
		{
			name: "unknown function",
			example: Example{
				Input: `\title{Hello, world!}

\banana{attack}
`,

				LoadErr: `undefined function \banana`,
			},
		},
		{
			name: "erroring single-return function",
			example: Example{
				Input: `\title{Hello, world!}

\use-plugin{errer}

\single-fail{some arg}

\multi-return-fail{some arg}
`,

				LoadErr: `function \single-fail returned an error: oh no`,
			},
		},
		{
			name: "erroring multi-return function",
			example: Example{
				Input: `\title{Hello, world!}

\use-plugin{errer}

\multi-fail{some arg}
`,

				LoadErr: `function \multi-fail returned an error: oh no`,
			},
		},
		{
			name: "ambiguous references",
			example: Example{
				Input: `\title{Hello, world!}

See also \reference{dupe-tag}{this tag}.

\section{
	\title{First Tag}

	\target{dupe-tag}{I'm the first tag.}
}

\section{
	\title{Second Tag}

	\target{dupe-tag}{I'm the second tag.}
}
`,

				RenderErr: `ambiguous target for tag 'dupe-tag'`,
			},
		},
		{
			name: "missing references",
			example: Example{
				Input: `\title{Hello, world!}

See also \reference{missing-tag}{this tag}.
`,

				RenderErr: `unknown tag 'missing-tag'`,
			},
		},
		{
			name: "setting title twice",
			example: Example{
				Input: `\title{Hello, world!}

\title{BAM}
`,

				LoadErr: "cannot set title twice",
			},
		},
	} {
		t.Run(tt.name, tt.example.Run)
	}
}
