// Package main provides CI/CD functions for the booklit project.
package main

import (
	"context"
	"fmt"

	"main/internal/dagger"
)

// Booklit represents the main Dagger module for the booklit project
type Booklit struct {
	// +private
	Source *dagger.Directory
}

func New(
	// +defaultPath="/"
	source *dagger.Directory,
) *Booklit {
	return &Booklit{
		Source: source,
	}
}

// Build builds the booklit binary, with an optional version.
// Returns the directory containing the built binary.
func (m *Booklit) Build(
	ctx context.Context,
	// +optional
	version string,
) (*dagger.Directory, error) {
	if version == "" {
		version = "dev"
	}

	fmt.Printf("Building booklit version: %s\n", version)

	// Build using Go container
	return dag.Container().
		From("golang:1.22-alpine").
		WithDirectory("/app", m.Source).
		WithWorkdir("/app").
		WithExec([]string{"go", "build", "-ldflags", fmt.Sprintf("-X github.com/vito/booklit.Version=%s", version), "-o", "booklit", "./cmd/booklit"}).
		Directory("/app"), nil
}

// Test runs all Go tests in the project.
func (m *Booklit) Test(ctx context.Context) (string, error) {
	return dag.Container().
		From("golang:1.22-alpine").
		WithDirectory("/app", m.Source).
		WithWorkdir("/app").
		WithExec([]string{"go", "test", "-v", "./..."}).
		Stdout(ctx)
}

// Lint runs golangci-lint against all Go code.
func (m *Booklit) Lint(ctx context.Context) (string, error) {
	return dag.Container().
		From("golangci/golangci-lint:latest").
		WithDirectory("/app", m.Source).
		WithWorkdir("/app").
		WithExec([]string{"golangci-lint", "run", "./..."}).
		Stdout(ctx)
}
