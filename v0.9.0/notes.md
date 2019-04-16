## new features

* `--out` is now optional; if not specified, the section will be rendered to
  `stdout`.

* `--tag` has been added; if specified, the section will be loaded and then
  only the section specified by the given tag will be rendered.

## misc

* the 'unknown template for (TYPE)' error condition is now returned earlier and
  returns an error specifying the name of the missing template.

* added a `joinLines` function, which is useful for implementing Markdown-style
  syntax where all lines but the first have to be inteded (e.g. lists).
