# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A personal blog platform (https://synic.dev) built with Go, Templ templates, HTMX, and Tailwind CSS.
Articles are written in Markdown with YAML frontmatter, converted to JSON at build time, embedded in the
binary, and served from memory.

## Build Commands

The build system uses [Mage](https://magefile.org/). Run targets with `mage <target>` (or
`go run github.com/magefile/mage@latest <target>` if mage isn't installed globally).

| Command                  | Description                                |
|--------------------------|--------------------------------------------|
| `mage`                   | Default: runs `dev` (hot-reload on :3000)  |
| `mage build:dev`         | Build debug binary to `bin/blog-debug`     |
| `mage build:release`     | Build release binary (tests first)         |
| `mage test`              | Run `go vet` then `go test -race ./...`    |
| `mage vet`               | Run `go vet ./...`                         |
| `mage codegen`           | Generate Templ components + Tailwind CSS   |
| `mage articles:convert`  | Convert changed markdown articles to JSON  |
| `mage articles:reconvert`| Force reconvert all articles               |
| `mage articles:create`   | Interactive: create new article (nvim)     |
| `mage pygmentize`        | Regenerate syntax highlighting CSS         |
| `mage clean`             | Remove built binaries                      |
| `mage container`         | Build Docker image                         |

To run a single test: `go test -run TestName ./internal/controller/`

## Architecture

**No web framework** — uses Go standard library `net/http` with `http.ServeMux`.

### Request Flow

`main.go` → embeds assets via `//go:embed` → loads articles into in-memory `ArticleRepository` → bundles
CSS/JS inline → registers routes on `http.ServeMux` → wraps with middleware (Logger → HTMX detection) →
serves on `:3000`.

### Package Layout (`internal/`)

- **`controller/`** — HTTP handlers. `ArticleController` handles index (with search/tag
  filter/pagination), single article, archive, and RSS feed.
- **`store/`** — `ArticleRepository` interface + `FSArticleRepository` (loads JSON articles from embedded
  FS into memory, indexed by slug).
- **`view/`** — Templ templates (`.templ` files) and rendering helpers. `BaseLayout` wraps all pages.
  `BundleStaticAssets` inlines CSS/JS into HTML. Generated `_templ.go` files should not be edited.
- **`model/`** — Data structures: `Article`, `ContextData`, `OpenGraphData`, `PageData`.
- **`converter/`** — Markdown-to-JSON pipeline: parses YAML frontmatter, renders Markdown via Goldmark
  with Chroma syntax highlighting.
- **`middleware/`** — Logger and HTMX detection middleware. `Wrap()` composes middleware.
- **`routes.go`** — Route registration. `static.go` — static file handler. `buildinfo.go` — build-time
  constants.

### Key Patterns

- **Functional options** — used for controller config (`WithPagination`) and render options (`WithStatus`,
  `WithOpenGraphData`).
- **HTMX partial rendering** — templates check if a request is an HTMX request to return partial HTML vs
  full page with layout.
- **Build-time article conversion** — Markdown → JSON happens at build time (or via
  `mage articles:convert`), not at runtime. JSON is embedded in the binary.

### Routes

```
GET  /                         → Article list (search: ?q=, tag: ?tag=, pagination: ?page=)
GET  /article/{date}/{slug}    → Single article
GET  /archive                  → Archive with tag counts
GET  /feed.xml                 → RSS feed
GET  /static/{path}            → Static assets
GET  /articles/{date}/{slug}   → 301 redirect to /article/...
```

## Article Format

Articles live in `articles/` as `YYYY-MM-DD_slug.md`. Frontmatter fields:

```yaml
---
title: Article Title
slug: article-slug
tags: [go, web]
publishedAt: 2024-01-01T00:00:00Z   # omit to keep unpublished
summary: |
  Short summary text
---
```

Articles without `publishedAt` are excluded from the index (unless debug mode is on) but are still
accessible via direct URL.

## Tech Stack

Go 1.24 · Templ (type-safe HTML templates) · HTMX · Tailwind CSS · Goldmark (Markdown) ·
Chroma (syntax highlighting) · Mage (build) · air (hot reload) · testify (test assertions)

Node deps (Tailwind, css-minify) are installed automatically by `mage dev` if missing.

## Datastar Directives

Always use the **colon-dot** form of Datastar attribute directives, never the underscore form.

- Correct: `data-on:click.prevent`
- Incorrect: `data-on-click__prevent`

The colon separates the directive from the event/target, and dots separate modifiers. For example:

- `data-on:click.prevent` (not `data-on-click__prevent`)
- `data-bind:value` (not `data-bind-value`)
- `data-on:submit.prevent` (not `data-on-submit__prevent`)
