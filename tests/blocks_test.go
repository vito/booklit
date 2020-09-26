package tests

import (
	. "github.com/onsi/ginkgo/extensions/table"
	_ "github.com/vito/booklit/tests/fixtures/arbitrary-style-plugin"
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

	Entry("tables", Example{
		Input: `\title{Hello, world!}

\table{
	\table-row{a}{1}
}{
	\table-row{b}{2}
}
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

<table>
	<tr>
		<td>a</td>
		<td>1</td>
	</tr>
	<tr>
		<td>b</td>
		<td>2</td>
	</tr>
</table>
</section>`,
		},
	}),

	Entry("definitions", Example{
		Input: `\title{Hello, world!}

\definitions{
	\definition{a}{1}
}{
	\definition{b}{2}
}
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

<dl>
	<dt>a</dt>
		<dd>1</dd>

	<dt>b</dt>
		<dd>2</dd>
</dl>
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

	Entry("aside", Example{
		Input: `\title{Hello, world!}

\aside{
	Hello.
}
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

<blockquote class="aside">
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

	Entry("arbitrary styles", Example{
		Input: `\title{Hello, world!}

\use-plugin{arbitrary-style}

\arbitrary-style{
	Sup!
}
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

<blink><p>Sup!</p></blink>
</section>`,
		},
	}),
)
