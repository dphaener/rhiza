# root-init

capability · `01M6RN84TZ5W0WZ9Z08W2HSFMX` · parent [runtime](../runtime.md)

## acceptance

A testscript runs rhiza root init, then rhiza root show --json from the root, from a registered project and from a <project>.<slug> sibling worktree, and gets the same root each time; a second root init on the same directory exits 2 with a next.

## statement

rhiza root init [dir] makes a directory a Rhiza root by writing <dir>/rhiza.toml (its projects and machine-wide settings), and every rhiza command finds its root by walking up from the working directory to the nearest rhiza.toml.

## Links

### Out

- decided-by [m29-registered-projects](../../../decision/m29-registered-projects.md)

### In

- verifies [root-init-verified](../../../evidence/root-init-verified.md)
- verifies [root-init-qa-1](../../../evidence/root-init-qa-1.md)
- verifies [root-init-qa-2](../../../evidence/root-init-qa-2.md)
- verifies [pr-335-root-and-registry](../../../evidence/pr-335-root-and-registry.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
