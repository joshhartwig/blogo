# Blogo

Blogo is a lightweight blog engine written in Go.

## Architecture

- cmd/blogo contains CLI entrypoint
- internal/server handles HTTP concerns
- internal/content handles markdown parsing and post repositories
- internal/templates handles template caching
- internal/config handles TOML configuration

## Conventions

- Prefer standard library over third-party packages
- Use html/template
- Wrap errors with fmt.Errorf(... %w ...)
- Keep handlers thin
- Prefer explicit structs over map[string]any

## Content Structure

Posts are folder-based:

content/
  my-post/
    index.md
    cover.jpg
    images/

## Future Goals

- Theme support
- Static site generation
- RSS improvements
- Search
