// Conflux - Full-stack CI/CD pipeline for the configuration management platform
//
// This module orchestrates the build, test, lint, and deployment pipeline for both
// the Go backend and SvelteKit frontend applications. It demonstrates a complete
// multi-module Dagger workflow that replaces Docker Compose with code-based
// container orchestration.

package main

import (
	"context"
	"dagger/conflux/internal/dagger"
	"fmt"
)

type Conflux struct{}

// BuildBackend builds the backend application only
func (m *Conflux) BuildBackend(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
) (*dagger.Container, error) {
	// Build backend with source option
	backendContainer := dag.Backend().Build(dagger.BackendBuildOpts{
		Source: source.Directory("backend"),
	})

	return backendContainer, nil
}

// TestBackend runs tests for the backend only
func (m *Conflux) TestBackend(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
) (string, error) {
	// Run backend tests
	backendTests, err := dag.Backend().Test(ctx, dagger.BackendTestOpts{
		Source: source.Directory("backend"),
	})
	if err != nil {
		return "", fmt.Errorf("backend tests failed: %w", err)
	}

	return backendTests, nil
}

// LintBackend runs linting for the backend only
func (m *Conflux) LintBackend(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
) (string, error) {
	// Run backend linting
	backendLint, err := dag.Backend().Lint(ctx, dagger.BackendLintOpts{
		Source: source.Directory("backend"),
	})
	if err != nil {
		return "", fmt.Errorf("backend linting failed: %w", err)
	}

	return backendLint, nil
}

// BuildAll builds both backend and frontend applications
func (m *Conflux) BuildAll(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
) (string, error) {
	// Build backend
	backendContainer := dag.Backend().Build(dagger.BackendBuildOpts{
		Source: source.Directory("backend"),
	})

	// Build frontend - try without options first
	frontendContainer := dag.Frontend().Build(source.Directory("frontend"))

	// Verify both builds completed successfully
	backendBinary := backendContainer.File("/app/server")
	frontendBuild := frontendContainer.Directory("/app/build")

	// Check that files exist by getting their sizes
	backendSize, err := backendBinary.Size(ctx)
	if err != nil {
		return "", fmt.Errorf("backend build failed: %w", err)
	}

	frontendFiles, err := frontendBuild.Entries(ctx)
	if err != nil {
		return "", fmt.Errorf("frontend build failed: %w", err)
	}

	return fmt.Sprintf("✅ Build successful!\nBackend binary: %d bytes\nFrontend files: %d",
		backendSize, len(frontendFiles)), nil
}

// TestAll runs tests for both backend and frontend
func (m *Conflux) TestAll(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
) (string, error) {
	// Run backend tests
	backendTests, err := dag.Backend().Test(ctx, dagger.BackendTestOpts{
		Source: source.Directory("backend"),
	})
	if err != nil {
		return "", fmt.Errorf("backend tests failed: %w", err)
	}

	// Run frontend tests - try without options
	frontendTests, err := dag.Frontend().Test(ctx, source.Directory("frontend"))
	if err != nil {
		return "", fmt.Errorf("frontend tests failed: %w", err)
	}

	return fmt.Sprintf("Backend Tests:\n%s\n\nFrontend Tests:\n%s",
		backendTests, frontendTests), nil
}

// LintAll runs linting for both backend and frontend
func (m *Conflux) LintAll(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
) (string, error) {
	// Run backend linting
	backendLint, err := dag.Backend().Lint(ctx, dagger.BackendLintOpts{
		Source: source.Directory("backend"),
	})
	if err != nil {
		return "", fmt.Errorf("backend linting failed: %w", err)
	}

	// Run frontend linting
	frontendLint, err := dag.Frontend().AllChecks(ctx, source.Directory("frontend"))
	if err != nil {
		return "", fmt.Errorf("frontend linting failed: %w", err)
	}

	return fmt.Sprintf("Backend Linting:\n%s\n\nFrontend Linting:\n%s",
		backendLint, frontendLint), nil
}

// CI runs the complete CI pipeline: test, lint, and build
func (m *Conflux) CI(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
) (string, error) {
	// Run tests
	testResults, err := m.TestAll(ctx, source)
	if err != nil {
		return "", fmt.Errorf("tests failed: %w", err)
	}

	// Run linting
	lintResults, err := m.LintAll(ctx, source)
	if err != nil {
		return "", fmt.Errorf("linting failed: %w", err)
	}

	// Build applications
	buildResults, err := m.BuildAll(ctx, source)
	if err != nil {
		return "", fmt.Errorf("build failed: %w", err)
	}

	return fmt.Sprintf("🎉 CI Pipeline Successful!\n\n%s\n\n%s\n\n%s",
		testResults, lintResults, buildResults), nil
}

// PackageBackend creates production container image for backend
func (c *Conflux) PackageBackend(source *dagger.Directory) *dagger.Container {
	return dag.Backend().Package(dagger.BackendPackageOpts{
		Source: source.Directory("backend"),
	})
}

// PackageFrontend creates production container image for frontend
func (c *Conflux) PackageFrontend(source *dagger.Directory) *dagger.Container {
	return dag.Frontend().Package(source.Directory("frontend"))
}

// PublishBackend builds and publishes backend container image
func (m *Conflux) PublishBackend(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
	// Container registry prefix to publish to
	// +default="ttl.sh/conflux-backend"
	registryPrefix string,
) (string, error) {
	// Publish backend
	backendRef, err := dag.Backend().Publish(ctx, dagger.BackendPublishOpts{
		Source:   source.Directory("backend"),
		Registry: registryPrefix,
	})
	if err != nil {
		return "", fmt.Errorf("backend publish failed: %w", err)
	}

	return fmt.Sprintf("🚀 Backend published successfully: %s", backendRef), nil
}

// Deploy is a stub for deployment functionality
func (m *Conflux) Deploy(
	ctx context.Context,
	// +defaultPath="."
	source *dagger.Directory,
	// Deployment environment
	// +default="staging"
	environment string,
) (string, error) {
	// TODO: Implement actual deployment logic
	// This could integrate with Kubernetes, Docker Swarm, or other deployment targets

	// For now, just run the backend pipeline and return deployment info
	_, err := m.TestBackend(ctx, source)
	if err != nil {
		return "", fmt.Errorf("backend tests failed before deployment: %w", err)
	}

	_, err = m.LintBackend(ctx, source)
	if err != nil {
		return "", fmt.Errorf("backend linting failed before deployment: %w", err)
	}

	_, err = m.BuildBackend(ctx, source)
	if err != nil {
		return "", fmt.Errorf("backend build failed before deployment: %w", err)
	}

	return fmt.Sprintf("🚀 Deployment to %s completed successfully (stub)", environment), nil
}
