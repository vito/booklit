package tests

import (
	. "github.com/onsi/ginkgo/extensions/table"
)

var _ = DescribeTable("Search Index", (Example).Run,
	Entry("sections", Example{
		Input: `\title{Hello, world!}

How are you?

Here's another paragraph.

\section{
	\title{How I'm doing}

	Good, thanks! And you?
}

\section{
	\title{Their Reply}

	Good, thanks!
}
`,

		SearchIndex: `{
			"hello-world": {
				"location": "hello-world.html",
				"title": "Hello, world!",
				"text": "How are you?\n\nHere's another paragraph.\n\n",
				"depth": 0,
				"section_tag": "hello-world"
			},
			"how-im-doing": {
				"location": "hello-world.html#how-im-doing",
				"title": "How I'm doing",
				"text": "Good, thanks! And you?\n\n",
				"depth": 1,
				"section_tag": "how-im-doing"
			},
			"their-reply": {
				"location": "hello-world.html#their-reply",
				"title": "Their Reply",
				"text": "Good, thanks!\n\n",
				"depth": 1,
				"section_tag": "their-reply"
			}
		}`,
	}),

	Entry("targets", Example{
		Input: `\title{Hello, world!}

How are you?

Here's another paragraph.

\section{
	\title{Sub-Section}

	Sub-section content.

	\target{some-target}{Some Target}{
		This is more information about some-target.
	}
}
`,

		SearchIndex: `{
			"hello-world": {
				"location": "hello-world.html",
				"title": "Hello, world!",
				"text": "How are you?\n\nHere's another paragraph.\n\n",
				"depth": 0,
				"section_tag": "hello-world"
			},
			"sub-section": {
				"location": "hello-world.html#sub-section",
				"title": "Sub-Section",
				"text": "Sub-section content.\n\n",
				"depth": 1,
				"section_tag": "sub-section"
			},
			"some-target": {
				"location": "hello-world.html#some-target",
				"title": "Some Target",
				"text": "This is more information about some-target.\n\n",
				"depth": 1,
				"section_tag": "sub-section"
			}
		}`,
	}),

	Entry("interesting content", Example{
		Input: `\title{Hello, world!}

How are you?

Here's a paragraph with \code{code}, a \link{link}{https://example.com},
and a \reference{sub-section}.

\table-of-contents

\list{
	Item 1

	Another line
}{
	Item 2
}

\ordered-list{
	Ordered Item 1
}{
	Ordered Item 2
}

\code{{
	line 1
	line 2
}}

\table{
	\table-row{a}{1}
}{
	\table-row{b}{2}
}

\definitions{
	\definition{a}{1}
}{
	\definition{b}{2}
}

\section{
	\title{Sub-Section}

	Sub-section content.
}
`,

		SearchIndex: `{
			"hello-world": {
				"location": "hello-world.html",
				"title": "Hello, world!",
				"text": "How are you?\n\nHere's a paragraph with code, a link, and a Sub-Section.\n\n* Item 1\n\n  Another line\n\n* Item 2\n\n1. Ordered Item 1\n\n2. Ordered Item 2\n\nline 1\nline 2\n| a | 1 |\n| b | 2 |\n\na: 1\nb: 2\n\n",
				"depth": 0,
				"section_tag": "hello-world"
			},
			"sub-section": {
				"location": "hello-world.html#sub-section",
				"title": "Sub-Section",
				"text": "Sub-section content.\n\n",
				"depth": 1,
				"section_tag": "sub-section"
			}
		}`,
	}),
)
