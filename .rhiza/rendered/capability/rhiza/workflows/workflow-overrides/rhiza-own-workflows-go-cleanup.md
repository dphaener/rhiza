# rhiza-own-workflows-go-cleanup

capability · `01M7642MCX86XXCXQ0FG85CM1T` · parent [workflow-overrides](../workflow-overrides.md)

## acceptance

rhiza workflow check passes on .rhiza/workflows and make self-check passes; a run of the cleanup command against a fake cache in a scratch dir removes only the stale entries.

## statement

Rhiza's own repo carries workflows in .rhiza/workflows that set GOTMPDIR to the run's temp folder and whose cleanup phase trims Go build cache entries unused for three hours, go-build temp folders older than two hours and rhiza test temp folders (rhiza-*) older than an hour.

## Links

### In

- verifies [rhiza-own-workflows-go-cleanup-verified](../../../../evidence/rhiza-own-workflows-go-cleanup-verified.md)
- verifies [rhiza-own-workflows-go-cleanup-qa-1](../../../../evidence/rhiza-own-workflows-go-cleanup-qa-1.md)
- verifies [pr-475-rhiza-s-go-behaviour-moves](../../../../evidence/pr-475-rhiza-s-go-behaviour-moves.md)

## History

- `01M763YX47S5XY0W3P8FVJ54S8` Rhiza's Go behaviour moves into its own workflows: workflows can set environment variables (Rhiza's set GOTMPDIR to the run's temp folder), Rhiza's cleanup phase trims the Go build cache and stale go-build and rhiza test temp folders, Rhiza's ship-change workflow builds rhiza from the PR head in a step; trimGoCache, gosweep's daemon use, the product's GOTMPDIR and ship's is-Rhiza-repo check leave the product.
