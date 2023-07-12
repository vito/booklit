package main

import (
	"dagger.io/dagger"
	"github.com/dagger/dagger/universe/apkoenv"
	"github.com/dagger/dagger/universe/goenv"
)

func main() {
	ctx := dagger.DefaultContext()
	ctx.Client().Environment().
		WithCheck_(Unit).
		WithCommand_(Build).
		Serve(ctx)
}

func Build(ctx dagger.Context, version string) (*dagger.Directory, error) {
	if version == "" {
		version = "dev"
	}

	return goenv.Build(ctx, Base(ctx), Code(ctx), goenv.GoBuildOpts{
		Packages: []string{"./cmd/booklit"},
		Xdefs:    []string{"github.com/vito/booklit.Version=" + version},
		Static:   true,
	}), nil
}

func Unit(ctx dagger.Context) (string, error) {
	return goenv.Test(ctx, Base(ctx), Code(ctx)).Stdout(ctx)
}

func Base(ctx dagger.Context) *dagger.Container {
	return apkoenv.Wolfi(ctx, []string{"go"})
}

func Code(ctx dagger.Context) *dagger.Directory {
	return ctx.Client().Host().Directory(".", dagger.HostDirectoryOpts{
		Include: []string{
			"**/*.go",
			"**/go.mod",
			"**/go.sum",
			"**/testdata/**/*",
			"**/*.proto",
			"**/*.tmpl",
			"**/*.bass",
		},
		Exclude: []string{
			"ci/**/*",
		},
	})
}