package main

import (
	"dagger.io/dagger"
	"github.com/vito/booklit/ci/pkgs"
)

func main() {
	dagger.ServeCommands(
		Build,
		Test,
	)
}

func Build(ctx dagger.Context, version string) (*dagger.Directory, error) {
	if version == "" {
		version = "dev"
	}

	return pkgs.GoBuild(ctx, Base(ctx), Code(ctx), pkgs.GoBuildOpts{
		Packages: []string{"./cmd/booklit"},
		Xdefs:    []string{"github.com/vito/booklit.Version=" + version},
		Static:   true,
	}), nil
}

func Test(ctx dagger.Context) (string, error) {
	return pkgs.GoTest(ctx, Base(ctx), Code(ctx)).Stdout(ctx)
}

func Base(ctx dagger.Context) *dagger.Container {
	return pkgs.Wolfi(ctx, []string{"go"})
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
