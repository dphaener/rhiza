# rebase-conflict-failed-steps-keep-output

question · `01M69QKMCRRGFSP7R4RDY2EQ83`

## statement

Rebasing branch cs/failed-steps-keep-output onto origin/main conflicts outside the rendered files, so ship-change stopped. How should these conflicts be resolved: internal/api/run.go?

## Links

### In

- answers [failed-steps-merge-resolution](../decision/failed-steps-merge-resolution.md)

## History

- `01M686RRCJSQ57M82FENWT2HNM` When a run's command step fails, its journal detail keeps the tail of the step's output, shown by run status and the inbox, so an unattended failure says why
