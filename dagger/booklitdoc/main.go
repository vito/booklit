// Booklitdoc provides docs-site helpers for the Booklit documentation as
// Dagger functions. Each function returns Booklit content serialized with the
// contentjson wire format (as JSON), which Booklit decodes back into native
// content when the function is called from a {expr} interpolation.
package main

import (
	"context"

	"dagger/booklitdoc/internal/dagger"
)

type Booklitdoc struct{}

// LitSyntax highlights Booklit source code and returns it as serialized
// Booklit content. `\function` references in the code are linkified into
// Booklit references, which resolve against the section the content is decoded
// into. The result is the JSON wire format; Booklit turns it back into native
// content.
func (m *Booklitdoc) LitSyntax(
	ctx context.Context,
	// Booklit source code to highlight.
	code string,
	// Tree-sitter language to use.
	// +optional
	// +default="lit"
	language string,
) (dagger.JSON, error) {
	if language == "" {
		language = "lit"
	}

	out, err := dag.Container().
		From("golang:1.26").
		WithDirectory("/src", dag.CurrentModule().Source()).
		WithWorkdir("/src").
		WithEnvVariable("CGO_ENABLED", "1").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("booklitdoc-go-mod")).
		WithMountedCache("/root/.cache/go-build", dag.CacheVolume("booklitdoc-go-build")).
		WithExec([]string{"go", "run", "./cmd/lit-syntax", "--code", code, "--language", language}).
		Stdout(ctx)
	if err != nil {
		return "", err
	}
	return dagger.JSON(out), nil
}
