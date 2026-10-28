# leaf-failed-01m60atz-1

question · `01M60EP81TAWBENX2YZVXBJ2XJ`

## statement

Review check qa still fails after 3 cycles: finding f1 (major at internal/workflow/builtin/ship-change.toml:19) on leaf ship-change-file: The render step commits only .rhiza/rendered. The engine writes each agent session's actor record into .rhiza/log in the worktree and leaves it uncommitted. A …. Its last cycle: cycle 3: confirmed f1 (major). How should it be fixed, or should leaf ship-change-file change?

## Links

### Out

- blocks [ship-change-file](../capability/ship-through-rhiza/ship-change-file.md)

### In

- answers [render-commits-actor-records](../decision/render-commits-actor-records.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
