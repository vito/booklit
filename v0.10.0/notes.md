## breaking changes

* `--tag` has been renamed to `--section-tag`

## new features

* `--section-path` has been added; if specified, the `--in` section will be
  loaded and then the given section path will be loaded and rendered with the
  `--in` section as its context.

  this is useful for things like release notes where the release notes aren't
  included into the main docs section but need it for resolving references.

* a basic text renderer has been added, using `text/template` instead of
  `html/template` and only including a few non-markup-specific templates. this
  is useful as a base for Markdown renderers.

* the `joinLines` function previously added in v0.9.0 has been moved to the
  new text renderer.

## misc

* the 'parsing section' log line has been switched to debug level.
