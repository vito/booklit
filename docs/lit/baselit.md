\use-plugin{booklitdoc}
\use-plugin{chroma}

# Basic Functions {#baselit}

Booklit comes with a default plugin. It's called `baselit`, but you don't
need to know that, since all sections have it by default!

The default plugin provides the basic set of functions needed for authoring
Booklit documents, plus many common functions useful for writing prose.

\table-of-contents

## Defining Sections {#sections}

\define{\title{text}{tags...}}{
  Set the title of the current section, along with optional *tags* as
  repeated arguments.

  For example, the `title` invocation for this section is:

  \lit-syntax{{{
  \title{Basic Functions}{baselit}
  }}}

  To specify multiple tags, pass multiple arguments:

  \lit-syntax{{{
  \title{I'm a taggy section!}{tag-1}{tag-1}
  }}}

  If no tags are specified, the section's tag defaults to a sanitized form of
  the title (e.g. *I'm a fancy title!* becomes
  `im-a-fancy-title`).
}

\define{\aux{content}}{
  Denote auxiliary content which can be stripped in certain contexts without
  losing meaning.

  Used within a title declaration to provide content that will show up on the
  section page itself but will be omitted when referencing the section. This
  is handy for sub-titles that you don't care to show anywhere but the
  section's page itself.

  Example:

  \lit-syntax{{{
  \title{Booklit\aux{: a pretty lit authoring system}}
  }}}

  This section, when referenced, would only show *Booklit*, but its
  header would include the *content*.
}

\define{\section{content}}{
  Introduce a sub-section, inheriting plugins of the outer section.

  Each sub-section should conventionally begin with a call to
  \reference{title} to set its title.

  For example, here's a full section containing a sub-section:

  \lit-syntax{{{
  \title{I'm a parent section!}

  Hello, world!

  \section{
    \title{And I'm a child section!}

    Waaah! Waaaah!
  }
  }}}

  Sections can be nested arbitrarily deep, however it is recommended to keep
  a maximum depth of two on a single page.

  Sections can be rendered on separate pages by invoking
  \reference{split-sections}.
}

\define{\include-section{path}}{
  Load the Booklit document located at *path* (relative to the current
  section's file) and set it as a child section.

  The included section **does not** inherit the plugins of the parent
  section. Instead, it should explicitly call \reference{use-plugin} on its
  own, so that it's self-contained.
}

\define{\split-sections}{
  Configure the renderer to generate a separate page for each immediate
  sub-section rather than inlining them under smaller headings.
}

\define{\single-page}{
  When declared in a section, it overrules any \reference{split-sections} in
  the section and any child sections (recursively) in order to force them all
  on to one page. Each section's header sizing is preserved, however.

  This is useful for having all of your content which is normally split
  across many pages joined in to an additional "single-page" format for
  quick skimming and searching.
}

\define{\table-of-contents}{
  Generate a block element that displays the table of contents from this
  section downward upon rendering. Often used in combination with
  \reference{split-sections}.
}

\define{\omit-children-from-table-of-contents}{
  Configure the section to omit its children from table of contents listings.
  This is appropriate when the sub-sections within a section are not quite
  standalone; they may be brief and meant to be consumed all at once, so
  navigating to them individually would not make sense.
}

## Customizing Sections

\define{\use-plugin{name}}{
  Register the plugin identified by *name* in the section. The plugin
  must be specified by `--plugin` on the command-line. See
  \reference{plugins} for more information.
}

\define{\styled{name}}{
  Set the template's style to *name*. The renderer may then use this
  to present the section in a different way. See \reference{styled-sections}.
}

\define{\set-partial{name}{content}}{
  Define the partial *name* in the section with *content* as
  its content.

  This is useful for communicating content to either \reference{plugins} or
  custom templates given to the \reference{html-renderer}{HTML renderer}.
}

## Links & References

\define{\link{display}{target}}{
  Link to *target* (i.e. a URL), with *display* as the link's
  text.

  For example, \syntax{lit}{\\link{Example}{https://example.com}} becomes
  [Example](https://example.com).

  Note that the argument order is the reverse of \reference{reference};
  writing tends to flow more naturally this way without a big URL
  interrupting the sentence.

  You can also use standard Markdown link syntax: `[Example](https://example.com)`.
}

\define{\reference{tag}{display?}}{
  Link to the target associated with *tag*. If the optional
  *display* argument is specified, it will be used as the link's
  content. Otherwise, the tag's configured display will be rendered.

  For example, \syntax{lit}{\\reference{links-and-references}} becomes
  \reference{links-and-references}.
}

\define{\target{tag}{display?}}{
  Generate a target element that can be \reference{reference}d as
  *tag*. If *display* is specified, references will default to
  showing it as their link. Otherwise, *tag* itself will be the
  default.

  As an example, we'll create a target element in the following paragraph,
  with the tag *some-tag* and display *I'm just some tag!*:

  \target{some-tag}{I'm just some tag!} I'm a targetable paragraph.

  Then, we'll reference it with \syntax{lit}{\\reference{some-tag}}:

  \reference{some-tag}

  Clicking the above link should take you to the paragraph.
}

## Flow Content

*Flow* content is anything that forms a *sentence*, i.e. a
string of words or inline elements.

\define{\code{text}}{
  Present *text* in a monospace font upon rendering.

  If *text* is a single line, it is presented as inline code. If it is
  in paragraph form, it is presented as its own block. See
  \reference{function-syntax} for more information.

  This is often used with `{{two}}` braces to preserve whitespace,
  or `{{{three}}}` to produce verbatim content (in addition to preserving
  whitespace). See \reference{function-syntax} for more information.

  You can also use Markdown backticks for inline code: `` `code bits` ``.

  For example:

  \lit-syntax{{{
  I'm a sentence with some `code bits` in it.
  }}}

  ...renders as:

  I'm a sentence with some `code bits` in it.

  While, this example:

  \lit-syntax{{{
  \code{{
    This is a code block.
  }}
  }}}

  ...renders as:

  \code{{
  This is a code block.
  }}
}

\define{\italic{text}}{
  Present *text* in *italics* upon rendering.

  You can also use Markdown syntax: `*text*`.
}

\define{\bold{text}}{
  Present *text* in **bold** upon rendering.

  You can also use Markdown syntax: `**text**`.
}

\define{\larger{text}}{
  Present *text* ~20% \larger{larger} upon rendering.
}

\define{\smaller{text}}{
  Present *text* ~20% \smaller{smaller} upon rendering.
}

\define{\strike{text}}{
  Present *text* with a \strike{strike through it} upon rendering.
}

\define{\superscript{text}}{
  Present *text* in \superscript{superscript} upon rendering.
}

\define{\subscript{text}}{
  Present *text* in \subscript{subscript} upon rendering.
}

\define{\image{path}}{
  Renders the image at *path* inline.

  Currently there is no "magic" that will do anything with the file specified
  by *path* - if it's a local path, you should make sure it's present
  in the directory that your documents are being generated into.

  You can also use Markdown syntax: `![alt text](path)`.
}

## Block Content

*Block* content is anything that forms a *paragraph*, i.e. a
block of text or any element that is standalone.

\define{\inset{content}}{
  Render *content* indented a bit.

  \inset{
    Like this!
  }

  You can also use Markdown blockquote syntax: `> content`.
}

\define{\aside{content}}{
  Render *content* in some way that conveys that it's a side-note.

  \aside{
    Here I am!
  }

  Depending on your screen size, you should either see it to the right or
  above this line.

  This is largely up to how you style them, though. You may want them to just
  look something like \reference{inset} if you don't have a "margin" in your
  design.
}

\define{\list{items...}}{
  Render an unordered list of *items*.

  \list{one}{two}{three!}

  You can also use Markdown syntax: `- item`.
}

\define{\ordered-list{items...}}{
  Render an ordered list of *items*.

  \ordered-list{one}{two}{three!}

  You can also use Markdown syntax: `1. item`.
}

\define{\table{rows...}}{
  Render a table with *rows* as its content.

  \target{table-row}{`\\`**table-row**} The value for each row should
  be produced by using \reference{table-row} like so:

  \lit-syntax{{{
  \table{
    \table-row{a}{b}{c}
  }{
    \table-row{1}{2}{3}
  }
  }}}

  The above example renders as:

  \table{
    \table-row{a}{b}{c}
  }{
    \table-row{1}{2}{3}
  }
}

\define{\definitions{entries...}}{
  Render a definition list with *entries* as its entries.

  \target{definition}{`\\`**definition**} The value for each entry
  should be produced by using \reference{definition} like so:

  \lit-syntax{{{
  \definitions{
    \definition{a}{1}
  }{
    \definition{b}{2}
  }
  }}}

  The above example renders as:

  \definitions{
    \definition{a}{1}
  }{
    \definition{b}{2}
  }
}
