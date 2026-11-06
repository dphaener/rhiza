# root-show

capability · `01M6RN84VVBPTANB7RAV73FSNX` · parent [runtime](../runtime.md)

## acceptance

Testscript cases in the root-init script cover the root, a project, a sibling worktree, the root itself (no project) and a directory outside any root (exit 2, next names root init).

## statement

rhiza root show [--json] prints the root, its config file and the project the working directory belongs to (by its folder or a <project>.<slug> sibling worktree), or none at the root itself; outside any root it exits 2 with next rhiza root init.

## Links

### Out

- decided-by [m29-registered-projects](../../../decision/m29-registered-projects.md)
- decided-by [m29-worktrees-siblings](../../../decision/m29-worktrees-siblings.md)

### In

- verifies [root-show-verified](../../../evidence/root-show-verified.md)
- verifies [root-show-qa-1](../../../evidence/root-show-qa-1.md)
- verifies [root-show-qa-2](../../../evidence/root-show-qa-2.md)
- verifies [pr-335-root-and-registry](../../../evidence/pr-335-root-and-registry.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
