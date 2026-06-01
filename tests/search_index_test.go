package tests

import (
	"testing"
)

func TestSearchIndex(t *testing.T) {
	for _, tt := range []struct {
		name    string
		example Example
	}{
		{
			name: "sections",
			example: Example{
				Input: `# Hello, world!

How are you?

Here's another paragraph.

## How I'm doing

Good, thanks! And you?

## Their Reply

Good, thanks!
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
			},
		},
		{
			name: "targets",
			example: Example{
				Input: `# Hello, world!

How are you?

Here's another paragraph.

## Sub-Section

Sub-section content.

<Target tag="some-target">Some Target</Target>
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
				"text": "Sub-section content.\n\n",
				"depth": 1,
				"section_tag": "sub-section"
			}
		}`,
			},
		},
		{
			name: "interesting content",
			example: Example{
				Input: "# Hello, world!\n" +
					"\n" +
					"How are you?\n" +
					"\n" +
					"Here's a paragraph with `code`, a [link](https://example.com),\n" +
					"and a <Reference tag=\"sub-section\"/>.\n" +
					"\n" +
					"<TableOfContents/>\n" +
					"\n" +
					"- Item 1\n" +
					"\n" +
					"  Another line\n" +
					"\n" +
					"- Item 2\n" +
					"\n" +
					"1. Ordered Item 1\n" +
					"\n" +
					"2. Ordered Item 2\n" +
					"\n" +
					"```\n" +
					"line 1\n" +
					"line 2\n" +
					"```\n" +
					"\n" +
					"| a | 1 |\n" +
					"| --- | --- |\n" +
					"| b | 2 |\n" +
					"\n" +
					"<Definitions>\n" +
					"<Definition term=\"a\">1</Definition>\n" +
					"<Definition term=\"b\">2</Definition>\n" +
					"</Definitions>\n" +
					"\n" +
					"## Sub-Section\n" +
					"\n" +
					"Sub-section content.\n",

				SearchIndex: `{
			"hello-world": {
				"location": "hello-world.html",
				"title": "Hello, world!",
				"text": "How are you?\n\nHere's a paragraph with code, a link, and a Sub-Section.\n\nItem 1\n\nAnother line\n\nItem 2\n\nOrdered Item 1\n\nOrdered Item 2\n\nline 1\nline 2\na1b2a: 1\nb: 2\n\n",
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
			},
		},
	} {
		t.Run(tt.name, tt.example.Run)
	}
}
