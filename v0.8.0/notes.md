## breaking changes

this release re-does how Booklit plugins are loaded. instead of generating and
compiling Go code on-the-fly, plugins are now loaded as [Go
plugins](https://golang.org/pkg/plugin/).

the main difference is that plugins will have to be `package main`, and that
plugins are only loaded once on start. so if you've changed a plugin you'll
need to restart the Booklit server.

because Go plugins only work on Linux and OS X, support for Booklit plugins on
Windows is now dependent on upstream Go changes. sorry about that.

## new features

* server mode will now render *only* the requested section, rather than
  rendering all sections and serving the requested file. this dramatically
  shortens the feedback cycle for large websites.

* section parsing is now cached based on file modification time. this also
  speeds things up quite a bit.

* HTML templates are *also* cached. you guessed it - things go faster.

## misc

generally, the switch from compiling and re-exec'ing to Go plugins cleaned
things up quite a bit and made the above optimizations possible to implement.
