# Summer CLI

The official command-line interface for the [Summer Framework](https://github.com/DropGuard/Summer).

Summer CLI provides a frictionless way to scaffold and manage Summer projects, bypassing complex Maven setups and archetype downloads. It is written in Go and compiles to a single, dependency-free binary.

## Features

- **Blazing Fast Scaffolding**: Embeds the Summer archetype directly in the binary. `summer create` works in under 0.1 seconds, completely offline.
- **Modern Developer Experience**: Wraps Maven execution (`summer dev`, `summer build`) for cleaner workflows.
- **Zero Configuration**: No need to fiddle with `settings.xml` or GitHub Packages tokens just to try the framework.

## Installation

```bash
# macOS / Linux
curl -fsSL https://raw.githubusercontent.com/DropGuard/summer-cli/main/install.sh | bash
```

```powershell
# Windows (Run as Administrator)
Invoke-WebRequest -Uri "https://github.com/DropGuard/summer-cli/releases/latest/download/summer-windows-amd64.exe" -OutFile "$env:SystemRoot\system32\summer.exe"
```

## Commands

### `summer create <project-name>`

Scaffold a new Summer framework project.

```bash
summer create my-app --group-id com.example --package com.example.app
```

### `summer dev`

Start the application in development mode with Jandex indexing and hot-reload. (Wraps `mvn summer:dev`).

```bash
cd my-app
summer dev
```

### `summer build`

Compile, test, and package the application into a production-ready Uber JAR using the AOT engine. (Wraps `mvn clean package`).

```bash
summer build
```

## Architecture

- **Cobra**: Built on the industry-standard `spf13/cobra` framework.
- **//go:embed**: The original `summer-archetype` is zipped and embedded directly into the Go binary at compile time, eliminating network calls during project creation.
- **E2E Tested**: The CLI suite includes a black-box E2E test that compiles the binary, scaffolds a project, and runs a real Maven build to verify structural integrity.
