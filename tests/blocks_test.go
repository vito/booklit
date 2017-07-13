package tests

import (
	. "github.com/onsi/ginkgo/extensions/table"
)

var _ = DescribeTable("Blocks", (Example).Run,
	Entry("lists", Example{
		Input: `\title{Hello, world!}

\list{a}{
	b
}{
	\code{{
	c
	}}
}
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

<ul><li>a</li><li><p>b</p></li><li><pre>c</pre></li></ul>
</section>`,
		},
	}),

	Entry("inset", Example{
		Input: `\title{Hello, world!}

\inset{
	Hello.
}
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

<div style="margin: 0 2em 1em" class="inset">
	<p>Hello.</p>
</div>
</section>`,
		},
	}),

	Entry("note", Example{
		Input: `\title{Hello, world!}

\note{
	Hello.
}
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

<blockquote class="note">
	<p>Hello.</p>
</blockquote>
</section>`,
		},
	}),

	Entry("ordered lists", Example{
		Input: `\title{Hello, world!}

\ordered-list{a}{
	b
}{
	\code{{
	c
	}}
}
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

<ol><li>a</li><li><p>b</p></li><li><pre>c</pre></li></ol>
</section>`,
		},
	}),
)
