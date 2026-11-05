# plan-merge-loops-to-rebase

decision · `01M6MT9X3ZK7259TEX95FVYW4W`

## statement

A plan run's merge step that merges nothing because origin/main moved past the checked head goes back to the rebase step, rebases, waits for the new head's CI and merges again, as ship-change's merge goes back to rebase-merge; the rebase, wait-rebased, rebased-green and merge steps are capped at 3 visits each, ship-change's cap, after which the run ends failed

## wording

(assumed, would have escalated) the change set's intent asks for the loop within a cap; the cap is ship-change's

## Links

### Out

- derived-from [plan-merge-rebases-interview](../source/plan-merge-rebases-interview.md)

### In

- decided-by [plan-run-merge-rebases](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-run-merge-rebases.md)

## History

- `01M6MT67M8VH509KK24BP65QRG` Fix Plan runs fail at merge when main moves: the plan workflow's merge step ends the run with 'Not merged: origin/main has moved past PR's checked head; rebase and check again' (run 01M6MM3J, PR #272, 2026-11-05) instead of going back to rebase and wait for green as ship-change does, so any merge landing between rebased-green and merge strands a plan PR that then also conflicts on .rhiza/rendered. The merge step should loop back to rebase, within a cap.
