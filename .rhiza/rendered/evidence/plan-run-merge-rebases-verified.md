# plan-run-merge-rebases-verified

evidence · `01M6MVHMXM0D3KH4GQM3TD64BA` · actor [claude-code-27115e41](../actor/claude-code-27115e41.md)

## result

PASS

## statement

At 22cb5a04 read the diff: plan.toml's merge step gains on.fail = rebase, and rebase, wait-rebased, rebased-green and merge each get max_visits = 3. go test -trimpath ./internal/workflow/... passes; TestPlanBuiltin asserts the merge fail->rebase transition and cap 3 on all four steps. go test -trimpath ./cmd/rhiza -run TestScripts/(plan-workflow-merge-main-moved|plan-workflow|plan-workflow-pr|run-ship-merge-main-moved|describe-golden) -count=1 passes. The new plan-workflow-merge-main-moved.txtar (modelled on run-ship-merge-main-moved) has a fake gh move origin/main once after the first rebased-green; rhiza run status shows merge (visit 1): fail, then rebase: exit 1; rebase (visit 2) pass; rebased-green (visit 2) pass; merge (visit 2): pass, then close-out; close-out pass, no visit 3; gh.log holds exactly one pr merge, 'pr merge 7 --squash --match-head-commit <merged head>', and the merged head descends from the other change set's commit, so rebase visit 2 pushed the rebased head; the run prints succeeded. The commit's .rhiza files are only the builder's own actor record.

## Links

### Out

- verifies [plan-run-merge-rebases](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-run-merge-rebases.md)

## History

- `01M6MT67M8VH509KK24BP65QRG` Fix Plan runs fail at merge when main moves: the plan workflow's merge step ends the run with 'Not merged: origin/main has moved past PR's checked head; rebase and check again' (run 01M6MM3J, PR #272, 2026-11-05) instead of going back to rebase and wait for green as ship-change does, so any merge landing between rebased-green and merge strands a plan PR that then also conflicts on .rhiza/rendered. The merge step should loop back to rebase, within a cap.
