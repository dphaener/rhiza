# m29-root-registry

intent · `01M6RN378KG7NSXMFC90YZ5QFW` · parent [milestone-29](../milestone-29.md)

## statement

A machine has one root directory whose config registers its projects, each added by path or git URL, and every Rhiza command inside a project knows its root and project.

## Links

### Out

- decided-by [m29-registered-projects](../../decision/m29-registered-projects.md)
- decided-by [m29-worktrees-siblings](../../decision/m29-worktrees-siblings.md)
- decided-by [m29-root-required](../../decision/m29-root-required.md)
- decided-by [m29-unready-listed](../../decision/m29-unready-listed.md)

### In

- serves change set `01M6RN3M8EC8FEDEJ8X6YVCPZ4` (Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project)

## History

- `01M6RMJ858AVSV2VA3K0TTV352` Rhiza can manage and build multiple projects on a single machine, from a single root directory.
