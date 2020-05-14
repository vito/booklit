package tests

import (
	. "github.com/onsi/ginkgo/extensions/table"
	_ "github.com/vito/booklit/tests/fixtures/stringer-plugin"
)

// TODO: backfill tests showing that the additional notation doesn't disrupt edge cases

var _ = DescribeTable("Booklit", (Example).Run,
	Entry("italics", Example{
		Input: `\title{Hello, world!}

How _are_ *you*?

_I'm_ good.
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How <em>are</em> <em>you</em>?</p>

	<p><em>I'm</em> good.</p>
</section>`,
		},
	}),

	Entry("bold", Example{
		Input: `\title{Hello, world!}

How __are__ **you**?
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How <strong>are</strong> <strong>you</strong>?</p>
</section>`,
		},
	}),

	Entry("inline code", Example{
		Input: "\\title{Hello, world!}\n\nHow `are` you?\n\nI'm good, thanks!\n",

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How <code>are</code> you?</p>

	<p>I'm good, thanks!</p>
</section>
`,
		},
	}),

	Entry("mixing and nesting like a madman", Example{
		Input: "\\title{Hello, world!}\n\nHow _**`are`**_ you?\n\nHow **_`are`_** you?\n\nHow *__`are`__* you?\n\nHow __*`are`*__ you?\n\nHow ___`are`___ you?\n\nHow ***`are`*** you?\n\nHow `**_are_**` you?\n",

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How <em><strong><code>are</code></strong></em> you?</p>
	<p>How <strong><em><code>are</code></em></strong> you?</p>
	<p>How <em><strong><code>are</code></strong></em> you?</p>
	<p>How <strong><em><code>are</code></em></strong> you?</p>
	<p>How <strong><em><code>are</code></em></strong> you?</p>
	<p>How <strong><em><code>are</code></em></strong> you?</p>
	<p>How <code>**_are_**</code> you?</p>
</section>`,
		},
	}),

	Entry("links", Example{
		Input: `\title{Hello, world!}

How are [you](https://example.com)?
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>How are <a href="https://example.com">you</a>?</p>
</section>`,
		},
	}),

	Entry("images", Example{
		Input: `\title{Hello, world!}

Here's an ![with alt text](foo.png) and another ![](without.gif).
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>Here's an <img src="foo.png" alt="with alt text" /> and another <img src="without.gif" alt="" />.</p>
</section>`,
		},
	}),

	Entry("headers", Example{
		Input: `\title{Hello, world!}

# How are you?
## How are you?
### How are you?
#### How are you?
##### How are you?
###### How are you?
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<h1>How are you?</h1>
	<h2>How are you?</h2>
	<h3>How are you?</h3>
	<h4>How are you?</h4>
	<h5>How are you?</h5>
	<h6>How are you?</h6>
</section>`,
		},
	}),

	Entry("blockquote", Example{
		Input: `\title{Hello, world!}

> This is a one-liner.

> This is line 1.
>
> This is line 2.

  > This one has spaces!
 > Isn't that neat?
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<blockquote>
		<p>This is a one-liner.</p>
	</blockquote>

	<blockquote>
		<p>This is line 1.</p>
		<p>This is line 2.</p>
	</blockquote>
</section>`,
		},
	}),

	Entry("tight unordered lists", Example{
		Input: `\title{Hello, world!}

* This.
* Is.
* A tight list.

- This.
- Is.
- Another tight list.
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<ul>
		<li>This.</li>
		<li>Is.</li>
		<li>A tight list.</li>
    </ul>

	<ul>
		<li>This.</li>
		<li>Is.</li>
		<li>Another tight list.</li>
    </ul>
</section>`,
		},
	}),
	Entry("loose unordered lists", Example{
		Input: `\title{Hello, world!}

* This is a list item.
  
  With multiple paragraphs.
* Therefore...
* Is is a loose list.

- This is a list item.
  
  With multiple paragraphs.
- Therefore...
- Is is another loose list.
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<ul>
		<li>
			<p>This is a list item.</p>
			<p>With multiple paragraphs.</p>
		</li>
		<li><p>Therefore...</p></li>
		<li><p>It is a loose list.</p></li>
    </ul>

	<ul>
		<li>
			<p>This is a list item.</p>
			<p>With multiple paragraphs.</p>
		</li>
		<li><p>Therefore...</p></li>
		<li><p>It is another loose list.</p></li>
    </ul>
</section>`,
		},
	}),

	Entry("tight ordered lists", Example{
		Input: `\title{Hello, world!}

1. This.
2. Is.
1. A tight list.

1) This.
2) Is.
1) Another tight list.
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<ol>
		<li>This.</li>
		<li>Is.</li>
		<li>A tight list.</li>
    </ol>

	<ul>
		<li>This.</li>
		<li>Is.</li>
		<li>Another tight list.</li>
    </ul>
</section>`,
		},
	}),
	Entry("loose ordered lists", Example{
		Input: `\title{Hello, world!}

1. This is a list item.
  
  With multiple paragraphs.
2. Therefore...
1. Is is a loose list.

1) This is a list item.
  
  With multiple paragraphs.
2) Therefore...
1) Is is another loose list.
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<ol>
		<li>
			<p>This is a list item.</p>
			<p>With multiple paragraphs.</p>
		</li>
		<li><p>Therefore...</p></li>
		<li><p>It is a loose list.</p></li>
    </ol>

	<ol>
		<li>
			<p>This is a list item.</p>
			<p>With multiple paragraphs.</p>
		</li>
		<li><p>Therefore...</p></li>
		<li><p>It is another loose list.</p></li>
    </ol>
</section>`,
		},
	}),

	Entry("horizontal rule", Example{
		Input: `\title{Hello, world!}

One!

---

Two!

***

Three!
`,

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>One!</p>

	<hr />

	<p>Two!</p>

	<hr />

	<p>Three!</p>
</section>`,
		},
	}),

	Entry("fenced codeblock", Example{
		Input: "\\title{Hello, world!}\n\nCheck out my code!\n\n```\nThis is a test.\n\n  This is wild and crazy.\n```\n\nSee?\n",

		Outputs: Files{
			"hello-world.html": `<section>
	<h1>Hello, world!</h1>

	<p>Check out my code!</p>

	<pre>This is a test.
  This is wild and crazy.</pre>

	<p>See?</p>
</section>
`,
		},
	}),
)
