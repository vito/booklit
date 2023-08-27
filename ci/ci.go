package main

import "context"

func main() {
	dag.Environment().
		WithCheck(Unit).
		WithCheck(Lint).
		WithCommand(Build).
		Serve()
}

// Build builds the booklit binary, with an optional version.
func Build(ctx context.Context, version string) (*Directory, error) {
	if version == "" {
		version = "dev"
	}

	return dag.Go().Build(Base(), Code(), GoBuildOpts{
		Packages: []string{"./cmd/booklit"},
		Xdefs:    []string{"github.com/vito/booklit.Version=" + version},
		Static:   true,
	}), nil
}

// Unit runs all Go tests.
func Unit(ctx context.Context) *EnvironmentCheck {
	return dag.Go().Test(Base(), Code(), GoTestOpts{
		Verbose: true,
	})
}

// Lint runs golangci-lint against all Go code.
func Lint(ctx context.Context) (string, error) {
	return dag.Go().GolangCilint(Base(), Code()).Stdout(ctx)
}

func Base() *Container {
	return dag.Apko().Wolfi([]string{"go", "golangci-lint"})
}

func Code() *Directory {
	return dag.Host().Directory(".", HostDirectoryOpts{
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
