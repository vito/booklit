# booklit


[![GoDoc](https://godoc.org/github.com/vito/booklit?status.svg)](https://godoc.org/github.com/vito/booklit)
[![CI](https://wings.concourse.ci/api/v1/teams/vito/pipelines/booklit/jobs/unit/badge)](https://wings.concourse.ci/teams/vito/pipelines/booklit/jobs/unit)

## installation

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
booklit -i ./docs/index.lit -o ./docs
```

you can see the result at https://vito.github.io/docs
