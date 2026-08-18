# template-bun

![build](https://img.shields.io/github/actions/workflow/status/mkvlrn/template-bun/checks.yml?branch=main&style=flat&logo=github)
[![template](https://img.shields.io/badge/template-use_this_template-2ea44f?style=flat&logo=github)](https://github.com/mkvlrn/template-bun/generate)
[![mise](https://mise-versions.jdx.dev/badge.svg)](https://mise.jdx.dev)
![license](https://img.shields.io/github/license/mkvlrn/template-bun?style=flat)

A sane, opinionated bun template.

> [!CAUTION]
> This template requires [mise](https://mise.jdx.dev) to manage runtimes, tools, and tasks in a single workflow, providing a lightweight alternative to devcontainers.
>
> You must [install mise](https://mise.jdx.dev/installing-mise.html) before starting. If you prefer a less opinionated setup, this template isn't for you.

Uses, among other tools/packages:

- [biome](https://github.com/biomejs/biome) for code linting and formatting
- [lefthook](https://github.com/evilmartians/lefthook) for git hooks
- [cocogitto](https://github.com/cocogitto/cocogitto) for commit message linting

## requirements and dependencies

As noted at the top, you need [mise](https://mise.jdx.dev) to get started with this template. Run `mise install` in the project root to fetch the pinned versions of Node and other tools locally.

This is _by far_ the easiest way to keep your environment consistent across different machines and team members, no matter the frequency of version updates.

Once the tooling is installed, you can install the Node dependencies with `bun install`.

> [!NOTE]
> Git hooks are in place to make sure both the tooling managed by mise and the project dependencies are synced with each checkout and merge.

## running

### `bun dev`

Runs the project in watch mode.

### `bun test`

Runs tests.

### `bun lint-fix`

Runs biome in fix mode to lint and format the project.

### `bun typecheck`

Runs type checking using tsc.

## ci

This repository uses GitHub Actions for CI. The workflow is defined in `.github/workflows/checks.yml`.

It automates:

- **Linting & Formatting**: Running Biome.
- **Type Checking**: Running TypeScript type checking.
- **Testing**: Running Bun test with code coverage.

## vscode

You might want to install the recommended extensions in vscode. Search for **@recommended** in the extensions tab, they'll show up as _"workspace recommendations"_.

If you have been using eslint and prettier and their extensions, you might want to disable eslint entirely and keep prettier as the formatter only for certain types of files.

This is done by the `.vscode/settings.json` file.

Debug configuration is also included for running the source directly with node.

## license

MIT
