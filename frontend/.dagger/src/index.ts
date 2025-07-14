/**
 * A Dagger module for Frontend (SvelteKit) functions
 *
 * This module provides build, test, lint, and packaging functions for the
 * SvelteKit frontend application.
 */
import { dag, Container, Directory, object, func } from "@dagger.io/dagger"

@object()
export class Frontend {
  /**
   * Returns a container with Node.js build environment set up
   */
  @func()
  buildEnvironment(): Container {
    return dag
      .container()
      .from("node:22-alpine")
      .withExec(["apk", "add", "--no-cache", "git"])
      .withWorkdir("/app")
      .withMountedCache("/root/.npm", dag.cacheVolume("npm-cache"))
  }

  /**
   * Install dependencies for the frontend application
   */
  @func()
  install(
    /**
     * Source directory containing the SvelteKit app
     * @defaultPath "."
     */
    source: Directory
  ): Container {
    return this.buildEnvironment()
      .withDirectory("/app", source)
      .withExec(["npm", "ci"])
  }

  /**
   * Build the frontend application
   */
  @func()
  build(
    /**
     * Source directory containing the SvelteKit app
     * @defaultPath "."
     */
    source: Directory
  ): Container {
    return this.install(source)
      .withExec(["npm", "run", "build"])
  }

  /**
   * Run frontend tests
   */
  @func()
  async test(
    /**
     * Source directory containing the SvelteKit app
     * @defaultPath "."
     */
    source: Directory
  ): Promise<string> {
    // For now, return a stub since SvelteKit tests need to be configured
    return "Frontend tests: SKIPPED (no tests configured yet)"
  }

  /**
   * Run linting checks on the frontend code
   */
  @func()
  async lint(
    /**
     * Source directory containing the SvelteKit app
     * @defaultPath "."
     */
    source: Directory
  ): Promise<string> {
    return this.install(source)
      .withExec(["npm", "run", "lint"])
      .stdout()
  }

  /**
   * Run type checking
   */
  @func()
  async check(
    /**
     * Source directory containing the SvelteKit app
     * @defaultPath "."
     */
    source: Directory
  ): Promise<string> {
    return this.install(source)
      .withExec(["npm", "run", "check"])
      .stdout()
  }

  /**
   * Package creates a production-ready container image
   */
  @func()
  package(
    /**
     * Source directory containing the SvelteKit app
     * @defaultPath "."
     */
    source: Directory
  ): Container {
    // Build the application
    const buildResult = this.build(source).directory("/app/build")

    // Create minimal production image with nginx
    return dag
      .container()
      .from("nginx:1.25-alpine")
      .withDirectory("/usr/share/nginx/html", buildResult)
      .withExposedPort(80)
  }

  /**
   * Publish builds and publishes the frontend container image
   * For now, this is stubbed out as requested
   */
  @func()
  async publish(
    /**
     * Source directory containing the SvelteKit app
     * @defaultPath "."
     */
    source: Directory,
    /**
     * Container registry to publish to
     * @default "ttl.sh/conflux-frontend"
     */
    registry: string = "ttl.sh/conflux-frontend"
  ): Promise<string> {
    // TODO: Implement actual publishing logic
    const _container = this.package(source)
    
    // Stub: Return the image reference that would be published
    return `${registry}:latest`
  }

  /**
   * Run all quality checks (lint + check)
   */
  @func()
  async allChecks(
    /**
     * Source directory containing the SvelteKit app
     * @defaultPath "."
     */
    source: Directory
  ): Promise<string> {
    // Run linting
    const lintResult = await this.lint(source)
    
    // Run type checking
    const checkResult = await this.check(source)
    
    return `Linting:\n${lintResult}\n\nType Checking:\n${checkResult}`
  }
}
