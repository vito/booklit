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

	Entry("forcing sections onto one page", Example{
		Input: `\title{Hello, world!}

How are you? See \reference{deep-inlined}.

\single-page
\split-sections

\section{
	\title{Section A}

	\split-sections

	Blah blah in section A.

	\section{
		\title{Nested Section}

		Good, thanks!
	}
}

\section{
	\title{Section B}

	Blah blah in section B.

	\section{
		\title{Nested Section 2}

		Foo bar.
	}

	\section{
		\title{Nested Section 3}

		\split-sections

		Fizz buzz.

		\section{
			\title{Super Duple Wrapped}{deep-inlined}

			Whoooooa.
		}
	}
}
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are you? See <a href="hello-world.html#deep-inlined">Super Duple Wrapped</a>.</p>

	<h1>1 Section A</h1>

	<p>Blah blah in section A.</p>

	<h1>1.1 Nested Section</h1>

	<p>Good, thanks!</p>

	<h1>2 Section B</h1>

	<p>Blah blah in section B.</p>

	<h2>2.1 Nested Section 2</h2>

	<p>Foo bar.</p>

	<h2>2.2 Nested Section 3</h2>

	<p>Fizz buzz.</p>

	<h1>2.2.1 Super Duple Wrapped</h1>

	<p>Whoooooa.</p>
</section>
`,
		},
	}),

	Entry("split sub-sub-sections", Example{
		Input: `\title{Hello, world!}

How are you?

\section{
	\title{How I'm Doing}

	\split-sections

	Good, thanks!

	\section{
		\title{Nested Section}

		Sup.
	}
}
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are you?</p>

	<h2>1 How I'm Doing</h2>

	<p>Good, thanks!</p>
</section>
`,
			"nested-section.html": `<section>
	<h1>1.1 Nested Section</h1>

	<p>Sup.</p>
</section>
`,
		},
	}),

	Entry("tables of contents", Example{
		Input: `\title{Hello, world!}

How are you?

\table-of-contents

This is some more content.

\section{
	\title{Top Section A}

	Foo bar.

	\section{
		\title{Nested Section}

		Fizz buzz.
	}

	\section{
		\title{Another Nested Section}

		Fiddlesticks.
	}

	\section{
		\title{Section with Omitted Children}

		I omit my children.

		\omit-children-from-table-of-contents

		\section{
			\title{Invisible Child}

			Boo!
		}
	}
}

\section{
	\title{Top Section B}

	Fiddlesticks is as far as I go.
}
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are you?</p>

	<ul>
		<li>
			<a href="hello-world.html#top-section-a">Top Section A</a>
			<ul>
				<li><a href="hello-world.html#nested-section">Nested Section</a></li>
				<li><a href="hello-world.html#another-nested-section">Another Nested Section</a></li>
				<li><a href="hello-world.html#section-with-omitted-children">Section with Omitted Children</a></li>
			</ul>
		</li>
		<li>
			<a href="hello-world.html#top-section-b">Top Section B</a>
		</li>
	</ul>

	<p>This is some more content.</p>

	<h2>1 Top Section A</h2>

	<p>Foo bar.</p>

	<h3>1.1 Nested Section</h3>

	<p>Fizz buzz.</p>

	<h3>1.2 Another Nested Section</h3>

	<p>Fiddlesticks.</p>

	<h3>1.3 Section with Omitted Children</h3>

	<p>I omit my children.</p>

	<h4>1.3.1 Invisible Child</h4>

	<p>Boo!</p>

	<h2>2 Top Section B</h2>

	<p>Fiddlesticks is as far as I go.</p>
</section>
`,
		},
	}),

	Entry("styled sections", Example{
		Input: `\title{Hello, world!}

\styled{top-template}

How are you?

\section{
	\title{How I'm doing}

	\styled{sub-template}

	Good, thanks! And you?

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
			"hello-world.html": `<section class="custom-top-page">
	<h1>Hello, world!</h1>

	<p>I'm a toplevel template! Here's my body:</p>

	<div class="custom-top-body">
		<p>How are you?</p>
	</div>

	<h2>1 How I'm doing</h2>

	<p>I'm a sub template! Here's my body:</p>

	<div class="custom-sub-body">
		<p>Good, thanks! And you?</p>
	</div>

	<h3>1.1 After Much Deliberation</h3>

	<p>I have decided that I'm doing well. How about you?</p>

	<h2>2 Their Reply</h2>

	<p>Good, thanks!</p>
</section>
`,
		},
	}),
)
