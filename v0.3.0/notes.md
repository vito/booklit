## breaking changes

* styled templates now execute with the `booklit.Styled` itself as `.`, with
  the unfortunate outcome being that any existing templates will cause an
  infinite loop upon rendering. make sure to change them all to `.Content` for
  the previous behavior!


## new features

* you can now attach other `booklit.Content` to a `booklit.Styled` via a
  `Partials` field. this allows for templates to be executed with multiple
  pieces of content, e.g. a `Title` partial to accompany the primary content.

* when `\styled` is set on a section, it will also result in checking for a
  `(style name)-page.tmpl` in place of the normal `page.tmpl` when the section
  is rendered at the top level


## fixes

* improved error propagation when using plugins and/or running the server
