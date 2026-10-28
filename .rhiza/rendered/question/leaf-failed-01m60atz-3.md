# leaf-failed-01m60atz-3

question · `01M60EP9ESXQF8WC6EJJJ1V5ZR`

## statement

Review check model still fails after 3 cycles: finding f2 (major at internal/workflow/builtin/ship-change.toml:19) on leaf ship-change-file: Commit 179a96f added a `render` command step between build and review. It runs rhiza render and commits .rhiza/rendered on the branch. The leaf's statement lis…. Its last cycle: cycle 3: confirmed f1 (major), f2 (major). How should it be fixed, or should leaf ship-change-file change?

## Links

### Out

- blocks [ship-change-file](../capability/ship-through-rhiza/ship-change-file.md)

### In

- answers [ship-change-file-names-render](../decision/ship-change-file-names-render.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
