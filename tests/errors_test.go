package tests

import (
	"testing"
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
