# ship-rebase-merges-only

capability · `01M6WMSB5S57K6QW44F90TEDEZ` · parent [rebase-when-behind](../rebase-when-behind.md)

## acceptance

run-ship-change-rebase.txtar covers: main moves with an unrelated commit, the step merges (the pushed head has main and the branch's old head as parents, origin's branch only fast-forwards); a rendered-only conflict is settled by rendering and pushed; a code conflict exits 1 with the merge still in progress, the conflicted files and main's commits and change sets in ship check --json, and no rebase-conflict question written; a merge committed by hand on top is pushed and recorded as pushed on the next ship rebase; plus the internal/api ship tests; run with go test -trimpath ./internal/api/... and go test -trimpath ./cmd/rhiza -run 'TestScripts/run-ship-change-rebase'

## statement

rhiza ship rebase always merges origin's main into the change set's branch (git merge, never git rebase) and pushes with a plain push (no --force, no lease), recording the pushed head; conflicts only under .rhiza/rendered are settled by rendering again; a conflict in code leaves the merge in progress in the worktree (not aborted), records in the ship state the conflicted files, origin/main's commits since the merge-base with their Change: change sets, and the branch's change set, shows them in rhiza ship check --json, raises no question and exits 1; when HEAD already holds origin/main and is ahead of origin's branch (a fixer's committed merge) it renders, commits what rendering changed, pushes plainly and records the pushed head

## Links

### Out

- decided-by [ship-merge-main-conflict-fixer](../../../../../decision/ship-merge-main-conflict-fixer.md)

### In

- depends-on [wait-merges-conflicting-pr](../../../orchestration/shipping/wait-green-pr/wait-merges-conflicting-pr.md)
- verifies [ship-rebase-merges-only-verified](../../../../../evidence/ship-rebase-merges-only-verified.md)
- verifies [ship-rebase-merges-only-qa-1](../../../../../evidence/ship-rebase-merges-only-qa-1.md)
- verifies [pr-383-fix-ship-change-stops-on-a](../../../../../evidence/pr-383-fix-ship-change-stops-on-a.md)

## History

- `01M6WMCCRC0KSCWEV17VY8KVQP` Fix Ship-change stops on a rebase conflict and waits for the owner instead of resolving it: run 01M6W0Z45JXE (change set 01M6VY7G6HPJ, Workflow dispatch, PR #367) stopped at wait-pr on 2026-11-08 ~01:20Z with 'rebasing cs/workflow-dispatch onto origin/main conflicts in describe-golden.txtar, internal/api/change.go, internal/api/run.go, internal/changeset/body.go, body_test.go, internal/daemon/requests.go, internal/store/read_test.go; the rebase was aborted; answer question rebase-conflict-workflow-dispatch'. With several runs in flight, main always moves under a long run, so this will recur. Instead, ship-change should merge origin/main into the branch (never rebase or force-push) and give the conflicts to a fixer agent session with both sides' commits and change sets as context. The fixer resolves them, regenerates golden files, runs the touched packages' tests, commits with trailers and pushes; the run records the push and goes on. It asks the owner (or the arbiter in an autonomous run) only if the fixer reports a real design clash. Unstuck by hand: an agent merged main in the run's worktree, answered the question and resumed the run.
