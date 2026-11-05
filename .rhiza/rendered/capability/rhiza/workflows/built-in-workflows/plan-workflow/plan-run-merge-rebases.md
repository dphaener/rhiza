# plan-run-merge-rebases

capability · `01M6MTA2NY5E0GWS0NDWVNEA6D` · parent [plan-workflow](../plan-workflow.md)

## acceptance

A Go workflow test shows the built-in plan workflow's merge step routing a fail to rebase and rebase, wait-rebased, rebased-green and merge each having max_visits 3; a testscript modelled on run-ship-merge-main-moved.txtar runs a plan run whose fake origin/main gains a commit once after rebased-green: the first merge makes no gh pr merge call and fails, the run's log shows merge (visit 1) fail then rebase, rebase (visit 2) pushes the rebased head, and merge (visit 2) calls gh pr merge pinned to the new head and the run reaches close-out and ends passed

## statement

When a plan run's merge finds origin/main has moved past the PR's checked head, the run goes back to rebase onto the new main (settling rendered-only conflicts by rendering again), waits for the rebased head's CI and merges it, up to 3 times, instead of ending the run with 'Not merged'

## Links

### Out

- decided-by [plan-merge-loops-to-rebase](../../../../../decision/plan-merge-loops-to-rebase.md)

### In

- verifies [plan-run-merge-rebases-verified](../../../../../evidence/plan-run-merge-rebases-verified.md)
- verifies [plan-run-merge-rebases-qa-1](../../../../../evidence/plan-run-merge-rebases-qa-1.md)
- verifies [pr-283-fix-plan-runs-fail-at-merge](../../../../../evidence/pr-283-fix-plan-runs-fail-at-merge.md)

## History

- `01M6MT67M8VH509KK24BP65QRG` Fix Plan runs fail at merge when main moves: the plan workflow's merge step ends the run with 'Not merged: origin/main has moved past PR's checked head; rebase and check again' (run 01M6MM3J, PR #272, 2026-11-05) instead of going back to rebase and wait for green as ship-change does, so any merge landing between rebased-green and merge strands a plan PR that then also conflicts on .rhiza/rendered. The merge step should loop back to rebase, within a cap.
