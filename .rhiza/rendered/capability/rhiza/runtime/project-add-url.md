# project-add-url

capability · `01M6RN8501TDM99PPTXZCXH3Q1` · parent [runtime](../runtime.md)

## acceptance

A testscript adds a local bare repo by its file:// URL (no network) and finds the clone under the root and in project ls; a bad URL exits non-zero and leaves rhiza.toml unchanged.

## statement

rhiza project add <git URL> clones the repo into <root>/<name> (name from the URL or --name) and registers it as project add <path> does; a failed clone registers nothing.

## Links

### Out

- decided-by [m29-registered-projects](../../../decision/m29-registered-projects.md)

### In

- verifies [project-add-url-verified](../../../evidence/project-add-url-verified.md)
- verifies [project-add-url-qa-1](../../../evidence/project-add-url-qa-1.md)
- verifies [project-add-url-qa-2](../../../evidence/project-add-url-qa-2.md)
- verifies [pr-335-root-and-registry](../../../evidence/pr-335-root-and-registry.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
