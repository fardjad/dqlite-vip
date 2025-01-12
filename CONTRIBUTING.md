# Contributing to dqlite-vip

This guide outlines how to set up a development environment and perform common
development tasks.

## Setting Up the Development Environment

The project includes a [`devcontainer.json`](./.devcontainer/devcontainer.json)
configuration file to set up a consistent development environment using Docker.

1. **Install Prerequisites:**
   - [Docker](https://www.docker.com/) (or a compatible container runtime)
   - A devcontainer-compatible editor, such as
     [Visual Studio Code](https://code.visualstudio.com/)

2. **Open the Project:**
   - Open the project directory in your editor.
   - When prompted, select **Reopen in Container**. If not prompted, open the
     command palette (`Ctrl+Shift+P` or `Cmd+Shift+P`), search for
     `Dev Containers: Reopen in Container`, and select it.

Once the container is ready, you'll have all required dependencies
pre-installed.

## Task Runner: `just`

The project uses [`just`](https://github.com/casey/just) as a task runner to
simplify development tasks. Run `just` without arguments to view a list of
available recipes:

```bash
just
```

## Common Development Tasks

**Build the binary:**

```bash
just build-static # for a statically linked binary
just build-dynamic # for a dynamically linked binary
```

**Run the tests:**

```bash
just test
```

**Generating Mocks:**

If you make changes to the interfaces, you must regenerate the mocks:

```bash
just generate-mocks
```
