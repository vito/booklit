package tests

import (
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/onsi/gomega"
	_ "github.com/vito/booklit/tests/fixtures/stringer-plugin"
)

var _ = DescribeTable("Booklit", (Example).Run,
	Entry("unknown function", Example{
		Input: `\title{Hello, world!}

\banana{attack}
`,

		Err: gomega.ContainSubstring("unknown function.lit:3: undefined function \\banana"),
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

		Err: gomega.MatchRegexp(`ambiguous references.lit:3: ambiguous reference tag 'dupe-tag'\n\ndefined in multiple locations:\n - .*ambiguous references.lit:8\n - .*ambiguous references.lit:14`),
	}),

	Entry("missing references", Example{
		Input: `\title{Hello, world!}

See also \reference{missing-tag}{this tag}.
`,

		Err: gomega.MatchRegexp(`missing references.lit:3: unknown tag 'missing-tag'`),
	}),
)
