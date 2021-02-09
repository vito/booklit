there's a fancy new website at [https://booklit.page](https://booklit.page)!

## refinements

* everything is godoc'd! this should make it a bit easier to use booklit APIs when writing plugins.

* the error UI style has been updated to match the website

* template rendering errors now use the new error UI

* styled sections no longer require a `(name)-page.tmpl` template

* the error page for referencing an unknown tag now suggests similar tags

* triple-curly-braces can now be used for single-line arguments, `\like{{{this}}}`!

* when `\title` is called twice in one section, an error will now be raised suggesting that the second `\title` should in its own `\section{...}`.

## fixes

* when running in server mode, `/foo/index.html` will no longer be treated as `/index.html`

## tweaks

* the `chroma` plugin will now set the HTML content as `.Content` in the styled template, rather than passing it as a partial.

* `\reference` no longer tries to match sections by their title - only tags, which is the documented and intended behavior.

* a log line will now be printed for each registered plugin
