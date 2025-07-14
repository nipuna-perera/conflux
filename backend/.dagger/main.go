package main

import (
	"context"
	"dagger/backend/internal/dagger"
	"fmt"
)

type Backend struct{}

// BuildEnvironment returns a container with Go build environment set up
func (m *Backend) BuildEnvironment() *dagger.Container {
	return dag.Container().
		From("golang:1.24-alpine").
		WithExec([]string{"apk", "add", "--no-cache", "git", "ca-certificates"}).
		WithWorkdir("/app").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-cache")).
		WithMountedCache("/root/.cache/go-build", dag.CacheVolume("go-build-cache"))
}

// Build compiles the backend application
func (m *Backend) Build(
	// +defaultPath="."
	source *dagger.Directory,
) *dagger.Container {
	return m.BuildEnvironment().
		WithDirectory("/app", source).
		WithExec([]string{"go", "mod", "download"}).
		WithExec([]string{"go", "build", "-o", "server", "./cmd/server"})
}

// Test runs backend tests
func (m *Backend) Test(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
) (string, error) {
	return m.BuildEnvironment().
		WithDirectory("/app", source).
		WithExec([]string{"go", "mod", "download"}).
		WithExec([]string{"go", "test", "-v", "./..."}).
		Stdout(ctx)
}

// Lint runs linting checks on the backend code
func (m *Backend) Lint(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
) (string, error) {
	// Install golangci-lint
	golangciLint := m.BuildEnvironment().
		WithExec([]string{"go", "install", "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"}).
		File("/go/bin/golangci-lint")

	return m.BuildEnvironment().
		WithFile("/usr/local/bin/golangci-lint", golangciLint).
		WithDirectory("/app", source).
		WithExec([]string{"golangci-lint", "run", "--verbose"}).
		Stdout(ctx)
}

// Package creates a production-ready container image
func (m *Backend) Package(
	// +defaultPath="."
	source *dagger.Directory,
) *dagger.Container {
	// Build the binary
	binary := m.Build(source).File("/app/server")

	// Create minimal production image
	return dag.Container().
		From("alpine:3.19").
		WithExec([]string{"apk", "add", "--no-cache", "ca-certificates"}).
		WithWorkdir("/app").
		WithFile("/app/server", binary).
		WithExposedPort(8080).
		WithEntrypoint([]string{"/app/server"})
}

// Publish builds and publishes the backend container image
// For now, this is stubbed out as requested
func (m *Backend) Publish(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
	// Container registry to publish to
	// +default="ttl.sh/conflux-backend"
	registry string,
) (string, error) {
	// TODO: Implement actual publishing logic
	_ = m.Package(source)

	// Stub: Return the image reference that would be published
	return fmt.Sprintf("%s:latest", registry), nil
}

// AllChecks runs all quality checks (test + lint)
func (m *Backend) AllChecks(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
) (string, error) {
	// Run tests
	testResult, err := m.Test(ctx, source)
	if err != nil {
		return "", fmt.Errorf("tests failed: %w", err)
	}

	// Run linting
	lintResult, err := m.Lint(ctx, source)
	if err != nil {
		return "", fmt.Errorf("linting failed: %w", err)
	}

	return fmt.Sprintf("Tests:\n%s\n\nLinting:\n%s", testResult, lintResult), nil
}
