# project-add-path

capability · `01M6RN84Y388QM22T4H109V2ZJ` · parent [runtime](../runtime.md)

## acceptance

A testscript adds a checkout and sees it in rhiza.toml and in rhiza project ls --json, and gets exit 2 for a path outside the root, a nested path, a non-git folder and a duplicate name.

## statement

rhiza project add <path> registers an existing git checkout that sits directly under the root, named by its folder (or --name), recording its repo ID when it has .rhiza/; a path outside the root or nested deeper, a non-git folder and a taken name are refused with exit 2 and a next.

## Links

### Out

- decided-by [m29-registered-projects](../../../decision/m29-registered-projects.md)
- decided-by [m29-worktrees-siblings](../../../decision/m29-worktrees-siblings.md)

### In

- verifies [project-add-path-verified](../../../evidence/project-add-path-verified.md)
- verifies [project-add-path-qa-1](../../../evidence/project-add-path-qa-1.md)
- verifies [project-add-path-qa-2](../../../evidence/project-add-path-qa-2.md)
- verifies [pr-335-root-and-registry](../../../evidence/pr-335-root-and-registry.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
