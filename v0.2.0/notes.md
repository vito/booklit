## breaking changes

* removed the `PluginFactory` interface in favor of just using a function type


## new features

* the `booklit` command now supports serving the content and rebuilding when it
  changes, by specifying `-s PORT`

* added `\styled` which allows sections to set their own "style", which the
  HTML renderer then uses to look up a template.


## fixes

* when an extra closing brace is present, a parse error will be returned,
  rather than blowing the heck up


## misc

* lots more documentation for booklit itself. pretty much everything should be
  covered now.

* renamed 'sentence' to 'line' in a bunch of places internally; i found when
  documenting the syntax that this terminology wasn't really appropriate
