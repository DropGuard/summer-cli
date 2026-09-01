# Summer CLI

The official command-line interface for the [Summer Framework](https://github.com/DropGuard/Summer).

Summer CLI provides a frictionless way to scaffold and manage Summer projects, bypassing complex Maven setups and archetype downloads. It is written in Go and compiles to a single, dependency-free binary.

## Features

- **Blazing Fast Scaffolding**: Embeds the Summer archetype directly in the binary. `summer create` works in under 0.1 seconds, completely offline.
- **Modern Developer Experience**: Wraps Maven execution (`summer dev`, `summer check`, `summer build`) for cleaner workflows.
- **Zero Configuration**: No need to fiddle with `settings.xml` or GitHub Packages tokens just to try the framework.

## Installation

```bash
# macOS / Linux
curl -fsSL https://raw.githubusercontent.com/DropGuard/summer-cli/main/install.sh | bash
```

```powershell
# Windows (PowerShell)
irm https://raw.githubusercontent.com/DropGuard/summer-cli/main/install.ps1 | iex
```

### Uninstallation
Because Summer CLI is a clean, single-file binary without any hidden background services or registry keys, uninstalling is as simple as deleting the executable:
```bash
# macOS / Linux
sudo rm /usr/local/bin/summer
```
```powershell
# Windows (PowerShell)
Remove-Item "$env:LOCALAPPDATA\Programs\summer" -Recurse -Force
```

## Commands

### `summer create <project-name>`

Scaffold a new Summer framework project.

```bash
summer create my-app --group-id com.example --package com.example.app
```

### `summer dev`

Start the application in development mode with Jandex indexing and hot-reload. (Wraps `mvn summer:dev`). Maven arguments are passed through.

```bash
cd my-app
summer dev
# for example: summer dev -Dsummer.engine=runtime
```

### `summer check`

Run compilation, Jandex discovery, dependency resolution, and AOT wiring without packaging a JAR. It starts clean by default, so the result is not affected by stale output.

```bash
summer check
summer check -DskipTests
```

Use `--incremental` only after a successful normal build has created the incremental source state:

```bash
summer check --incremental
```

### `summer build`

Compile, test, and package the application into a production-ready Uber JAR using the AOT engine. The safe default wraps `mvn clean package` and removes all previous output.

```bash
summer build
```

For a faster, state-checked build that reuses `target/`:

```bash
summer build --incremental
```

`--clean` is accepted as an explicit spelling of the default safe mode. All other arguments are passed to Maven:

```bash
summer build --clean -DskipTests
summer build --incremental -DskipTests
```

## Architecture

- **Cobra**: Built on the industry-standard `spf13/cobra` framework.
- **//go:embed**: The original `summer-archetype` is zipped and embedded directly into the Go binary at compile time, eliminating network calls during project creation.
- **E2E Tested**: The CLI suite includes a black-box E2E test that compiles the binary, scaffolds a project, and runs a real Maven build to verify structural integrity.
