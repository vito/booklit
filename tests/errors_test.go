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

		Err: gomega.ContainSubstring("unknown function.lit:3: undefined function \\banana"),
	}),

	Entry("erroring single-return function", Example{
		Input: `\title{Hello, world!}

\use-plugin{errer}

\single-fail{some arg}

\multi-return-fail{some arg}
`,

		Err: gomega.ContainSubstring("erroring single-return function.lit:5: failed to evaluate \\single-fail: oh no"),
	}),

	Entry("erroring multi-return function", Example{
		Input: `\title{Hello, world!}

\use-plugin{errer}

\multi-fail{some arg}
`,

		Err: gomega.ContainSubstring("erroring multi-return function.lit:5: failed to evaluate \\multi-fail: oh no"),
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

		Err: gomega.MatchRegexp(`ambiguous references.lit:3: ambiguous target for tag 'dupe-tag'\n\ntag 'dupe-tag' is defined in multiple locations:\n\n - .*ambiguous references.lit:8\n - .*ambiguous references.lit:14`),
	}),

	Entry("missing references", Example{
		Input: `\title{Hello, world!}

See also \reference{missing-tag}{this tag}.
`,

		Err: gomega.MatchRegexp(`missing references.lit:3: unknown tag 'missing-tag'`),
	}),
)
