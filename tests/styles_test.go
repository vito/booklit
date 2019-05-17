package tests

import (
	. "github.com/onsi/ginkgo/extensions/table"
	_ "github.com/vito/booklit/tests/fixtures/partials-style-plugin"
)

var _ = DescribeTable("Booklit", (Example).Run,
	Entry("styling with partials", Example{
		Input: `\title{Hello, world!}

\use-plugin{partial-style}

\block-style{Title A}{
	Hello, \target{target-a}{some target} \reference{target-b}!
}

\block-style{Title B}{
	Hello again, \target{target-b}{some other target} \reference{target-a}!
}

This is an \inline-style{Title C}{inline style}!

\block-style{Title D}{This is a line forced into block style!}
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

  <div class="custom-style">
		<h3>Title A</h3>

		<p>Hello, <a id="target-a"></a> <a href="hello-world.html#target-b">some other target</a>!</p>
	</div>

  <div class="custom-style">
		<h3>Title B</h3>

		<p>Hello again, <a id="target-b"></a> <a href="hello-world.html#target-a">some target</a>!</p>
	</div>

	<p>This is an <span class="inline-style"><strong>Title C</strong>: inline style</span>!</p>

	<div class="custom-style">
		<h3>Title D</h3>

		This is a line forced into block style!
	</div>
</section>
`,
		},
	}),
)
