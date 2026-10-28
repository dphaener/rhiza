# leaf-failed-01m60atz-4

question · `01M61CGZBK1E1FB2ERQP0X5VQP`

## statement

Review check qa still fails after 4 cycles: finding f1 (major at internal/workflow/builtin/ship-change.toml:61) on leaf ship-change-file: The cycle-3 fix holds on the approve path, but rhiza run request-changes still pushes a branch that fails rhiza check. signoff's on.changes goes to review/revi…. Its last cycle: cycle 4: confirmed f1 (major). How should it be fixed, or should leaf ship-change-file change?

## Links

### Out

- blocks [ship-change-file](../capability/ship-through-rhiza/ship-change-file.md)

### In

- answers [render-before-open-pr](../decision/render-before-open-pr.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
