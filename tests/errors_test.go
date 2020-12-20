package tests

import (
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/onsi/gomega"
	_ "github.com/vito/booklit/tests/fixtures/erroring-plugin"
)

var _ = DescribeTable("Booklit", (Example).Run,
	Entry("unknown function", Example{
		Input: `\title{Hello, world!}

\banana{attack}
`,

		Err: gomega.ContainSubstring("undefined function \\banana"),
	}),

	Entry("erroring single-return function", Example{
		Input: `\title{Hello, world!}

\use-plugin{errer}

\single-fail{some arg}

\multi-return-fail{some arg}
`,

		Err: gomega.ContainSubstring("function \\single-fail returned an error: oh no"),
	}),

	Entry("erroring multi-return function", Example{
		Input: `\title{Hello, world!}

\use-plugin{errer}

\multi-fail{some arg}
`,

		Err: gomega.ContainSubstring("function \\multi-fail returned an error: oh no"),
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

		Err: gomega.MatchRegexp(`ambiguous target for tag 'dupe-tag'`),
	}),

	Entry("missing references", Example{
		Input: `\title{Hello, world!}

See also \reference{missing-tag}{this tag}.
`,

		Err: gomega.MatchRegexp(`unknown tag 'missing-tag'`),
	}),

	Entry("setting title twice", Example{
		Input: `\title{Hello, world!}

\title{BAM}
`,

		Err: gomega.ContainSubstring("cannot set title twice"),
	}),
)
