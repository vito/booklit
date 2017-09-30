## new features

* ambiguous references (those which can be satisfied by multiple tags) now
  result in an error, just like broken references (those matching no tags). to
  suppress the error, pass `--allow-broken-references`, as with broken
  references.

* broken and ambiguous references will result in a warning during compilation;
  this makes it easier to hunt them down and fix them while using
  `--allow-broken-references`.

* sped up build time when building with plugins, by using `go install` instead
  of `go build`
