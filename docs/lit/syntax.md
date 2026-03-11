\use-plugin{booklitdoc}

# Syntax {#syntax}

Booklit documents are Markdown files extended with a special syntax for
[function calls](#function-syntax). Standard Markdown formatting
(emphasis, bold, links, code, lists, etc.) is supported natively, and
everything else is either text or an `\invoke` call.

\table-of-contents

## Prose Syntax {#prose-syntax}

Booklit builds on top of standard Markdown, so the prose rules are familiar:

\list{
  The top-level of a document is a series of *paragraphs*, separated
  by one or more blank lines.
}{
  A *paragraph* is a series of *lines*, separated by linebreaks.
  Adjacent lines within a paragraph are joined (soft line breaks become
  spaces).
}{
  *Emphasis* is written with `*asterisks*` and **bold** with
  `**double asterisks**`.
}{
  *Inline code* is written with `` `backticks` ``.
}{
  *Links* are written as `[display text](url)`.
}{
  *Images* are written as `![alt text](path)`.
}{
  *Headings* can be written with `#` prefix, which maps to
  [#title].
}{
  In addition to Markdown formatting, [function
  calls](#function-syntax) can be used inline or at the block level.
}

## Comment Syntax {#comment-syntax}

Comments are delimited by `{-` and `-}`. They can be multi-line,
appear in between words, and they can also be nested. This makes commenting
out large blocks of content easy:

\lit-syntax{{{
Hi, I'm{- a comment -} in the middle of a sentence!

{-
  I'm hanging out at the top level,

  {- being nested and stuff -}

  with multiple lines.
-}
}}}

## Function Syntax {#function-syntax}

Function calls are denoted by `\` followed by a series of alphanumeric
characters and hyphens (`foo-bar`), forming the function *name*.

To produce a literal `\` character, use `\\` (standard Markdown backslash
escape).

Following the name, there may be any number of *arguments*, which can
come in a few different forms:

\definitions{
  \definition{`{line}`}{
    With no linebreak after the `{`, the argument forms a single
    line. Markdown formatting is applied within the argument.
  }
}{
  \definition{`{word wrapped line}`}{
    As above, but soft line breaks are converted into a single
    space, as if it were written as `{word wrapped line}`.
  }
}{
  \definition{\lit-syntax{{{
{
  paragraph 1

  paragraph 2
}
  }}}}{
    With a linebreak after the `{`, the argument forms a block of
    paragraphs with full Markdown and function call support.
  }
}{
  \definition{\lit-syntax{{{
{{
  paragraph 1

    indented paragraph 2

  \with{syntax}
}}
  }}}}{
    With doubled-up curly braces, whitespace is preserved in the content,
    rather than being parsed into paragraphs. Function calls (`\invoke`)
    are still recognized within preformatted blocks.

    Note that the first line of the content determines an indentation level
    that is then stripped for all lines. It is the only whitespace that is
    ignored.
  }
}{
  \definition{\lit-syntax{{
{{{
  paragraph 1

    indented {paragraph} 2

  \\not-parsed{no-syntax}
}}}
  }}}{
    Tripled-up curly braces form a verbatim argument. Similar to
    preformatted, whitespace is preserved. In addition, there is no
    interpreting or parsing of function calls or Markdown within.
    This is useful for large code blocks where the content may contain
    special characters like `\`, `{`, or `}`.
  }
}
