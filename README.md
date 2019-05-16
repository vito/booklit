# booklit

[![GoDoc](https://godoc.org/github.com/vito/booklit?status.svg)](https://godoc.org/github.com/vito/booklit)
[![CI](https://wings.pivotal.io/api/v1/teams/vito/pipelines/booklit/jobs/unit/badge)](https://wings.pivotal.io/teams/vito/pipelines/booklit/jobs/unit)

## installation

grab the latest [release](https://github.com/vito/booklit/releases), or build
from source:

```bash
go get github.com/vito/booklit/cmd/booklit
```

## usage

```bash
booklit -i foo.lit -o ./out
```

## example

clone this repo and build its docs:

```bash
booklit \
  -i docs/lit/index.lit \
  -o docs \
  --html-templates docs/lit/html \
  --plugin github.com/vito/booklit/booklitdoc \
  --plugin github.com/vito/booklit/chroma/plugin
```

then browse the generated docs from `./docs/index.html`.

you can see the result at https://vito.github.io/booklit
