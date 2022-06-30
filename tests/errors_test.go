package tests

import (
	. "github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	_ "github.com/vito/booklit/tests/fixtures/erroring-plugin"
)

var _ = DescribeTable("Booklit", (Example).Run,
	Entry("unknown function", Example{
		Input: `\title{Hello, world!}

\banana{attack}
`,

		LoadErr: gomega.ContainSubstring("undefined function \\banana"),
	}),

	Entry("erroring single-return function", Example{
		Input: `\title{Hello, world!}

\use-plugin{errer}

\single-fail{some arg}

\multi-return-fail{some arg}
`,

		LoadErr: gomega.ContainSubstring("function \\single-fail returned an error: oh no"),
	}),

	Entry("erroring multi-return function", Example{
		Input: `\title{Hello, world!}

\use-plugin{errer}

\multi-fail{some arg}
`,

		LoadErr: gomega.ContainSubstring("function \\multi-fail returned an error: oh no"),
	}),

	Entry("ambiguous references", Example{
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

		RenderErr: gomega.MatchRegexp(`ambiguous target for tag 'dupe-tag'`),
	}),

	Entry("missing references", Example{
		Input: `\title{Hello, world!}

See also \reference{missing-tag}{this tag}.
`,

		RenderErr: gomega.MatchRegexp(`unknown tag 'missing-tag'`),
	}),

	Entry("setting title twice", Example{
		Input: `\title{Hello, world!}

\title{BAM}
`,

		LoadErr: gomega.ContainSubstring("cannot set title twice"),
	}),
)
