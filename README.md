# booklit

[![Go Reference](https://pkg.go.dev/badge/github.com/vito/booklit.svg)](https://pkg.go.dev/github.com/vito/booklit)
[![CI](https://ci.concourse-ci.org/api/v1/teams/main/pipelines/booklit/jobs/unit/badge)](https://ci.concourse-ci.org/teams/main/pipelines/booklit/jobs/unit)

## documentation

[booklit.page](https://booklit.page)

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
