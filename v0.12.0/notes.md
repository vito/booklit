there's a fancy new website at [https://booklit.page](https://booklit.page)!

## refinements

* the error UI style has been updated to match the website

* template rendering errors now use the new error UI

* styled sections no longer require a `(name)-page.tmpl` template

* the error page for referencing an unknown tag now suggests similar tags

* triple-curly-braces can now be used for single-line arguments, `\like{{{this}}}`!

## fixes

* when running in server mode, `/foo/index.html` will no longer be treated as `/index.html`
