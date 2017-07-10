package tests

import (
	. "github.com/onsi/ginkgo/extensions/table"
)

var _ = DescribeTable("Booklit", (Example).Run,
	Entry("sub-sections", Example{
		Input: `\title{Hello, world!}

How are you?

\section{
	\title{How I'm doing}

	Good, thanks! And you?
}

\section{
	\title{Their Reply}

	Good, thanks!
}
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are you?</p>

	<h2>1 How I'm doing</h2>

	<p>Good, thanks! And you?</p>

	<h2>2 Their Reply</h2>

	<p>Good, thanks!</p>
</section>
`,
		},
	}),

	Entry("sub-sections from files", Example{
		Input: `\title{Hello, world!}

How are you?

\include-section{how-im-doing.lit}

\section{
	\title{Their Reply}

	Good, thanks!
}
`,

		Inputs: Files{
			"how-im-doing.lit": `\title{How I'm doing}

Good, thanks! And you?
`,
		},

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are you?</p>

	<h2>1 How I'm doing</h2>

	<p>Good, thanks! And you?</p>

	<h2>2 Their Reply</h2>

	<p>Good, thanks!</p>
</section>
`,
		},
	}),

	Entry("nested sub-sections", Example{
		Input: `\title{Hello, world!}

How are you?

\section{
	\title{How I'm doing}

	\section{
		\title{After Much Deliberation}

		I have decided that I'm doing well. How about you?
	}
}

\section{
	\title{Their Reply}

	Good, thanks!
}
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are you?</p>

	<h2>1 How I'm doing</h2>

	<h3>1.1 After Much Deliberation</h3>

	<p>I have decided that I'm doing well. How about you?</p>

	<h2>2 Their Reply</h2>

	<p>Good, thanks!</p>
</section>
`,
		},
	}),

	Entry("split sub-sections", Example{
		Input: `\title{Hello, world!}

How are you?

\split-sections

\section{
	\title{How I'm Doing}

	Good, thanks!
}
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are you?</p>
</section>
`,
			"how-im-doing.html": `<section>
	<h1>1 How I'm Doing</h1>

	<p>Good, thanks!</p>
</section>
`,
		},
	}),
)
