# plan-run-merge-rebases-qa-1

evidence · `01M6MVT8ZMFCTDPC3EVGGBKW9S` · actor [claude-code-bf74acc5](../actor/claude-code-bf74acc5.md)

## result

PASS

## statement

Built rhiza from the branch head (go build -trimpath) and ran rhiza run start --foreground --workflow plan by hand in a scratch repo with a bare origin and fake rex/claude/gh (plan-workflow-merge-main-moved.txtar's). With origin/main moved once after rebased-green: merge (visit 1) failed with 'Not merged: origin/main has moved past PR #7's checked head', went to rebase; rebase (visit 2) pushed a head on the new main, merge (visit 2) made the only gh pr merge call, pinned with --match-head-commit to that head (which contains the other commit), close-out passed and the run succeeded. With origin/main moving on every check: merge failed at visits 1-3 with no gh pr merge call, and the run ended failed with 'step "rebase" went over max_visits 3'. rhiza workflow describe plan shows merge on.fail -> rebase. TestPlanBuiltin and the plan-workflow-merge-main-moved, plan-workflow-pr and run-ship-merge-main-moved scripts pass. The rendered-only conflict path of ship rebase was not exercised.

## Links

### Out

- verifies [plan-run-merge-rebases](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-run-merge-rebases.md)

## History

- `01M6MT67M8VH509KK24BP65QRG` Fix Plan runs fail at merge when main moves: the plan workflow's merge step ends the run with 'Not merged: origin/main has moved past PR's checked head; rebase and check again' (run 01M6MM3J, PR #272, 2026-11-05) instead of going back to rebase and wait for green as ship-change does, so any merge landing between rebased-green and merge strands a plan PR that then also conflicts on .rhiza/rendered. The merge step should loop back to rebase, within a cap.
