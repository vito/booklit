package tests

import (
	. "github.com/onsi/ginkgo/extensions/table"
	_ "github.com/vito/booklit/tests/fixtures/data-style-plugin"
)

var _ = DescribeTable("Booklit", (Example).Run,
	Entry("styling with custom data", Example{
		Input: `\title{Hello, world!}

\use-plugin{data-style}

\struct-style{Hello, \target{target-a}{some target} \reference{target-b}!}

\map-style{Hello again, \target{target-b}{some other target} \reference{target-a}!}

This is an \inline-style{inline style}!
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

  <div class="data-style">Hello, <a name="target-a"></a> <a href="hello-world.html#target-b">some other target</a>!</div>

  <div class="data-style">Hello again, <a name="target-b"></a> <a href="hello-world.html#target-a">some target</a>!</div>

  <p>This is an <span class="inline-style">inline style</span>!</p>
</section>
`,
		},
	}),
)
