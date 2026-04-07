package tests

import (
	"testing"

	_ "github.com/vito/booklit/tests/fixtures/stringer-plugin"
)

func TestProse(t *testing.T) {
	for _, tt := range []struct {
		name    string
		example Example
	}{
		{
			name: "simple Hello World",
			example: Example{
				Input: `\title{Hello, world!}

How are you?
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are you?</p>
</section>`,
				},
			},
		},
		{
			name: "no trailing linebreak",
			example: Example{
				Input: `\title{Hello, world!}

How are you?`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are you?</p>
</section>`,
				},
			},
		},
		{
			name: "link",
			example: Example{
				Input: `\title{Hello, world!}

How are \link{you}{https://example.com}?
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are <a href="https://example.com">you</a>?</p>
</section>`,
				},
			},
		},
		{
			name: "images",
			example: Example{
				Input: `\title{Hello, world!}

Here's an \image{foo.png}{with alt text} and another \image{without.gif}.
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>Here's an <img src="foo.png" alt="with alt text" /> and another <img src="without.gif" alt="" />.</p>
</section>`,
				},
			},
		},
		{
			name: "italics",
			example: Example{
				Input: `\title{Hello, world!}

How are \italic{you}?
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are <em>you</em>?</p>
</section>`,
				},
			},
		},
		{
			name: "bold",
			example: Example{
				Input: `\title{Hello, world!}

How are \bold{you}?
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are <strong>you</strong>?</p>
</section>`,
				},
			},
		},
		{
			name: "larger",
			example: Example{
				Input: `\title{Hello, world!}

How are \larger{you}?
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are <span style="font-size: 120%">you</span>?</p>
</section>`,
				},
			},
		},
		{
			name: "smaller",
			example: Example{
				Input: `\title{Hello, world!}

How are \smaller{you}?
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are <span style="font-size: 80%">you</span>?</p>
</section>`,
				},
			},
		},
		{
			name: "strike",
			example: Example{
				Input: `\title{Hello, world!}

How are \strike{you}?
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are <span style="text-decoration: line-through">you</span>?</p>
</section>`,
				},
			},
		},
		{
			name: "superscript",
			example: Example{
				Input: `\title{Hello, world!}

How are \superscript{you}?
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are <sup>you</sup>?</p>
</section>`,
				},
			},
		},
		{
			name: "subscript",
			example: Example{
				Input: `\title{Hello, world!}

How are \subscript{you}?
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are <sub>you</sub>?</p>
</section>`,
				},
			},
		},
		{
			name: "multiple paragraphs",
			example: Example{
				Input: `\title{Hello, world!}

How are you?

I'm good, thanks!
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are you?</p>

	<p>I'm good, thanks!</p>
</section>
`,
				},
			},
		},
		{
			name: "invokes interspersed in words",
			example: Example{
				Ext: ".lit",

				Input: `\title{Hello, world!}

This{\italic{is}}a test.
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>This<em>is</em>a test.</p>
</section>
`,
				},
			},
		},
		{
			name: "word-wrapped lines",
			example: Example{
				Input: `\title{Hello, world!}

Lorem ipsum dolor sit amet, consectetur adipiscing elit. \italic{Curabitur
accumsan a ligula id feugiat. Quisque luctus semper ex sodales vulputate.} Sed
mi mi, rhoncus non justo et, aliquam dictum est. Donec egestas massa id
pharetra scelerisque. Nulla nunc quam, sagittis vel est sed, ultrices bibendum
magna. Nulla posuere ut erat eget tristique. Nullam vel nisl vitae dui
sollicitudin porta.

\section{
	\title{Indented}

	Integer malesuada purus dignissim turpis lacinia fringilla. Suspendisse
	potenti. Maecenas varius iaculis volutpat. \italic{Vestibulum sagittis lacus
	ut ex varius molestie.}
}
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. <em>Curabitur accumsan a ligula id feugiat. Quisque luctus semper ex sodales vulputate.</em> Sed mi mi, rhoncus non justo et, aliquam dictum est. Donec egestas massa id pharetra scelerisque. Nulla nunc quam, sagittis vel est sed, ultrices bibendum magna. Nulla posuere ut erat eget tristique. Nullam vel nisl vitae dui sollicitudin porta.</p>

	<h2>1 Indented</h2>

	<p>Integer malesuada purus dignissim turpis lacinia fringilla. Suspendisse potenti. Maecenas varius iaculis volutpat. <em>Vestibulum sagittis lacus ut ex varius molestie.</em></p>
</section>
`,
				},
			},
		},
		{
			name: "word-wrapped paragraphs",
			example: Example{
				Input: `\title{Hello, world!}

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur accumsan a
ligula id feugiat. Quisque luctus semper ex sodales vulputate. Sed mi mi,
rhoncus non justo et, aliquam dictum est. Donec egestas massa id pharetra
scelerisque. Nulla nunc quam, sagittis vel est sed, ultrices bibendum magna.
Nulla posuere ut erat eget tristique. Nullam vel nisl vitae dui sollicitudin
porta.

Integer malesuada purus dignissim turpis lacinia fringilla. Suspendisse
potenti. Maecenas varius iaculis volutpat. Vestibulum sagittis lacus ut ex
varius molestie.
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur accumsan a ligula id feugiat. Quisque luctus semper ex sodales vulputate. Sed mi mi, rhoncus non justo et, aliquam dictum est. Donec egestas massa id pharetra scelerisque. Nulla nunc quam, sagittis vel est sed, ultrices bibendum magna. Nulla posuere ut erat eget tristique. Nullam vel nisl vitae dui sollicitudin porta.</p>

	<p>Integer malesuada purus dignissim turpis lacinia fringilla. Suspendisse potenti. Maecenas varius iaculis volutpat. Vestibulum sagittis lacus ut ex varius molestie.</p>
</section>
`,
				},
			},
		},
		{
			name: "inline code and code blocks",
			example: Example{
				Ext: ".lit",

				Input: `\title{Hello, world!}

This is some \code{inline} code.

Here's a code block:

\code{
	I'm a code block.

		I'm indented more.

			I'm indented even more.

I'm indented less.

	\reference{hello-world}

	\\some-method\{Some argument.\}


	One more line, with meaning.
}

\code{{
	I'm a code block.

		I'm indented more.

			I'm indented even more.

I'm indented less.

	\reference{hello-world}

	\\some-method\{Some argument.\}


	One more line, with meaning.
}}

\code{{{
	I'm a code block.

		I'm indented more.

			I'm indented even more.

I'm indented less.

	\reference{hello-world}

	\\some-method\{Some argument.\}


	One more line, with meaning.
}}}

And here's some more content.
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>This is some <code>inline</code> code.</p>

	<p>Here's a code block:</p>

	<pre><p>I'm a code block.</p><p>I'm indented more.</p><p>I'm indented even more.</p><p>I'm indented less.</p><p><a href="hello-world.html">Hello, world!</a></p><p>\some-method{Some argument.}</p><p>One more line, with meaning.</p></pre>

	<pre>I'm a code block.

	I'm indented more.

		I'm indented even more.

I'm indented less.

<a href="hello-world.html">Hello, world!</a>

\some-method{Some argument.}


One more line, with meaning.</pre>

	<pre>I'm a code block.

	I'm indented more.

		I'm indented even more.

I'm indented less.

\reference{hello-world}

\\some-method\{Some argument.\}


One more line, with meaning.</pre>

	<p>And here's some more content.</p>
</section>
`,
				},
			},
		},
		{
			name: "code block indent tracking",
			example: Example{
				Ext: ".lit",

				Input: `\title{Hello, world!}

\code{{
I'm a code block.
}}

\section{
  \title{Sub-section}

	\code{{
	I'm a code block in a sub-section.
  }}

  \code{{
    {}   I'm a code block {- in a sub-section -}with a forced indent level.
  }}
}
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<pre>I'm a code block.</pre>

	<h2>1 Sub-section</h2>

	<pre>I'm a code block in a sub-section.</pre>

	<pre>   I'm a code block with a forced indent level.</pre>
</section>
`,
				},
			},
		},
		{
			name: "empty arguments",
			example: Example{
				Input: `\title{Hello, world!}

This is an \italic{} empty italic.

This is a \italic{ } space italic.

This is an \italic{  } even more spaced italic.
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>This is an <em></em> empty italic.</p>
	<p>This is a <em> </em> space italic.</p>
	<p>This is an <em>  </em> even more spaced italic.</p>
</section>
`,
				},
			},
		},
		{
			name: "empty multi-line arguments",
			example: Example{
				Input: `\title{Hello, world!}

\code{
}

\code{{
}}

\code{{{
}}}

\code{

}

\code{{

}}

\code{{{

}}}

\code{


}

\code{{


}}

\code{{{


}}}
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p><code></code></p>
	<pre></pre>
	<pre></pre>
	<p><code></code></p>
	<pre></pre>
	<pre></pre>
	<p><code></code></p>
	<pre>
</pre>
	<pre>
</pre>
</section>
`,
				},
			},
		},
		{
			name: "preformatted string arguments",
			example: Example{
				Input: `\title{Hello, world!}

\use-plugin{stringer}

Here's a code block:

\string{{{
	I'm a code block.

		I'm indented more.

			I'm indented even more.

I'm indented less.

	\reference{hello-world}

	\\some-method\{Some argument.\}


	One more line, with meaning.
}}}

And here I'm just using it to \string{{{escape {{wacky}} curlies}}}.
`,

				Outputs: Files{
					"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>Here's a code block:</p>

	<p>I'm a code block.

	I'm indented more.

		I'm indented even more.

I'm indented less.

\reference{hello-world}

\\some-method\{Some argument.\}


One more line, with meaning.
</p>

<p>And here I'm just using it to escape {{wacky}} curlies.</p>
</section>
`,
				},
			},
		},
	} {
		t.Run(tt.name, tt.example.Run)
	}
}
