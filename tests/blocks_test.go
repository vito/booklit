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
