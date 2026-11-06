# project-ls-rm

capability · `01M6RN85109TG2YXD5Z6VNQ6DF` · parent [runtime](../runtime.md)

## acceptance

A testscript registers one initialised and one bare git checkout, sees initialised and not initialised in project ls --json, removes one with project rm and finds its folder still on disk.

## statement

rhiza project ls [--json] lists each registered project with its name, path, repo ID and whether it is initialised (has .rhiza/) or not initialised; rhiza project rm <name> unregisters a project without touching its files.

## Links

### Out

- decided-by [m29-unready-listed](../../../decision/m29-unready-listed.md)

### In

- verifies [project-ls-rm-verified](../../../evidence/project-ls-rm-verified.md)
- verifies [project-ls-rm-qa-1](../../../evidence/project-ls-rm-qa-1.md)
- verifies [project-ls-rm-qa-2](../../../evidence/project-ls-rm-qa-2.md)
- verifies [pr-335-root-and-registry](../../../evidence/pr-335-root-and-registry.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
