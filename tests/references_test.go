package tests

import (
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/vito/booklit"
)

var _ = DescribeTable("Booklit", (Example).Run,
	Entry("references to other sections on the same page", Example{
		Input: `\title{Hello, world!}

See also \reference{section-c}{the last section}.

\section{
	\title{Section A}

	See also \reference{section-b}.
}

\section{
	\title{Section B}

	See also \reference{section-a}.
}

\section{
	\title{Section C}

	See also \reference{hello-world}.
}
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>See also <a href="hello-world.html#section-c">the last section</a>.</p>

	<h2>1 Section A</h2>

	<p>See also <a href="hello-world.html#section-b">Section B</a>.</p>

	<h2>2 Section B</h2>

	<p>See also <a href="hello-world.html#section-a">Section A</a>.</p>

	<h2>3 Section C</h2>

	<p>See also <a href="hello-world.html">Hello, world!</a>.</p>
</section>
`,
		},
	}),

	Entry("references to other sections by title", Example{
		Input: `\title{Hello, world!}

See also \reference{Section C}{the last section}.

\section{
	\title{Section A}

	See also \reference{Section B}.
}

\section{
	\title{Section B}

	See also \reference{Section A}.
}

\section{
	\title{Section C}

	See also \reference{Hello, world!}.
}
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>See also <a href="hello-world.html#section-c">the last section</a>.</p>

	<h2>1 Section A</h2>

	<p>See also <a href="hello-world.html#section-b">Section B</a>.</p>

	<h2>2 Section B</h2>

	<p>See also <a href="hello-world.html#section-a">Section A</a>.</p>

	<h2>3 Section C</h2>

	<p>See also <a href="hello-world.html">Hello, world!</a>.</p>
</section>
`,
		},
	}),

	Entry("references to other sections on split pages", Example{
		Input: `\title{Hello, world!}

See also \reference{section-c}{the last section}.

\split-sections

\section{
	\title{Section A}

	See also \reference{section-b}.
}

\section{
	\title{Section B}

	See also \reference{section-a}.
}

\section{
	\title{Section C}

	See also \reference{hello-world}.
}
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>See also <a href="section-c.html">the last section</a>.</p>
</section>
`,
			"section-a.html": `<section>
	<h1>1 Section A</h1>

	<p>See also <a href="section-b.html">Section B</a>.</p>
</section>
`,
			"section-b.html": `<section>
	<h1>2 Section B</h1>

	<p>See also <a href="section-a.html">Section A</a>.</p>
</section>
`,
			"section-c.html": `<section>
	<h1>3 Section C</h1>

	<p>See also <a href="hello-world.html">Hello, world!</a>.</p>
</section>
`,
		},
	}),

	Entry("explicit target elements", Example{
		Input: `\title{Hello, world!}

\reference{target-a}

\reference{target-a}{with display}

\reference{target-without-display}

\reference{target-without-display}{with display}

\section{
	\title{Some Section}

	Foo bar.

	\target{target-a}{Target A} Here's target A.

	\target{target-without-display} Here's another target.
}
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p><a href="hello-world.html#target-a">Target A</a></p>
	<p><a href="hello-world.html#target-a">with display</a></p>
	<p><a href="hello-world.html#target-without-display">target-without-display</a></p>
	<p><a href="hello-world.html#target-without-display">with display</a></p>

	<h2>1 Some Section</h2>

	<p>Foo bar.</p>
	<p><a id="target-a"></a> Here's target A.</p>
	<p><a id="target-without-display"></a> Here's another target.</p>
</section>
`,
		},
	}),

	Entry("aux", Example{
		Input: `\title{Hello, world!\aux{: Foo Bar}}

See also \reference{section-c}{the last section}.

\table-of-contents

\section{
	\title{Section A}

	See also \reference{section-b}.
}

\section{
	\title{Section B\aux{aby}}

	See also \reference{some-anchor}.
}

\section{
	\title{Section C}

	\target{some-anchor}{I'm an\aux{ awesome} anchor.}See also \reference{hello-world}.
}
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!: Foo Bar</h1>

	<p>See also <a href="hello-world.html#section-c">the last section</a>.</p>

	<ul>
		<li><a href="hello-world.html#section-a">Section A</a></li>
		<li><a href="hello-world.html#section-b">Section B</a></li>
		<li><a href="hello-world.html#section-c">Section C</a></li>
	</ul>

	<h2>1 Section A</h2>

	<p>See also <a href="hello-world.html#section-b">Section B</a>.</p>

	<h2>2 Section Baby</h2>

	<p>See also <a href="hello-world.html#some-anchor">I'm an anchor.</a>.</p>

	<h2>3 Section C</h2>

	<p><a id="some-anchor"></a>See also <a href="hello-world.html">Hello, world!</a>.</p>
</section>
`,
		},
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

		Err: booklit.ErrBrokenReference,
	}),

	Entry("missing references", Example{
		Input: `\title{Hello, world!}

See also \reference{missing-tag}{this tag}.
`,

		Err: booklit.ErrBrokenReference,
	}),
)
