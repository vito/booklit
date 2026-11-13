# booklit

[![Go Reference](https://pkg.go.dev/badge/github.com/vito/booklit.svg)](https://pkg.go.dev/github.com/vito/booklit)

Booklit is a tool for building static websites from semantic documents.

## documentation

[booklit.page](https://booklit.page)

## installation

grab the latest [release](https://github.com/vito/booklit/releases), or build
from source:

```bash
go get github.com/vito/booklit/cmd/booklit
go install github.com/vito/booklit/cmd/booklit
```

## usage

```bash
booklit -i foo.lit -o ./out
```

## example

clone this repo and build the Booklit website:

```bash
./scripts/build-docs
```

then browse the generated docs from `./docs/index.html`.

alternatively, run the docs in server mode:

```bash
./scripts/build-docs -s 3000
```

...and then browse to [localhost:3000](https://localhost:3000)

## development

tests, generated files and release builds all run with [Dagger](https://dagger.io):

```bash
dagger check      # tests, cross-compilation, and that generated files are up to date
dagger generate   # regenerate the parser, the stylesheets, and the docs site
```

`dagger generate` covers everything derived from a source file:

| generated                                    | from                              |
| -------------------------------------------- | --------------------------------- |
| `ast/booklit.peg.go`                         | `ast/booklit.peg` (`go generate`) |
| `errhtml/errors.css`, `docs/css/booklit.css` | `less/*.less`                     |
| `docs/*.html`, `docs/outputs/`               | `docs/lit/`, `docs/hello/`        |

to cut a release, push `<version>/name` and `<version>/notes.md` to the
`version` branch, then:

```bash
dagger call release --version v0.13.0 --github-token "cmd://gh auth token"
```
