# conflict-fix-design-clash

capability · `01M6WMRS6BFQHDFXKF5RNPBF0G` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

run-ship-conflict-fix.txtar gains a case in an autonomous run where the fake fixer blocks on a question about a design clash: the run journals an asked entry for conflict-fix and adds the question to the run's change set, the fake arbiter decides it, the answer is delivered to the same fixer session, which commits the merge, and the run pushes and goes on to merge with one fixer session for that conflict; a non-autonomous case shows rhiza run status --json blocked at conflict-fix with kind question and no push while the session waits; a case where the fake fixer says done without committing shows conflict-fix failing because its until does not hold, with no push; run with go test -trimpath ./cmd/rhiza -run 'TestScripts/run-ship-conflict-fix' and go test -trimpath ./internal/run/...

## statement

When the conflict fixer finds a real design clash (both sides change the same behaviour incompatibly) it commits nothing and does not end its session: it leaves the merge in progress and asks one blocking question in its session, naming both change sets, the clashing files and the incompatible behaviours, using the engine's existing in-session block (no engine change, no exit 75): in an autonomous run the engine's ask path adds the question to the run's change set, has the arbiter decide it and delivers the answer to the same waiting session; otherwise the run shows blocked at conflict-fix and the owner answers in the session's tab. The same fixer session then resolves the merge following the answer, commits it with its trailers and says done, so the step's until holds and the run goes on to render-pr; a fixer that ends without committing the merge fails the step's until (the merge still in progress), and the run fails rather than passing on

## Links

### Out

- depends-on [ship-change-conflict-fix-step](ship-change-conflict-fix-step.md)
- decided-by [ship-merge-main-conflict-fixer](../../../../../decision/ship-merge-main-conflict-fixer.md)

### In

- verifies [conflict-fix-design-clash-verified](../../../../../evidence/conflict-fix-design-clash-verified.md)
- verifies [conflict-fix-design-clash-qa-1](../../../../../evidence/conflict-fix-design-clash-qa-1.md)
- verifies [pr-383-fix-ship-change-stops-on-a](../../../../../evidence/pr-383-fix-ship-change-stops-on-a.md)

## History

- `01M6WMCCRC0KSCWEV17VY8KVQP` Fix Ship-change stops on a rebase conflict and waits for the owner instead of resolving it: run 01M6W0Z45JXE (change set 01M6VY7G6HPJ, Workflow dispatch, PR #367) stopped at wait-pr on 2026-11-08 ~01:20Z with 'rebasing cs/workflow-dispatch onto origin/main conflicts in describe-golden.txtar, internal/api/change.go, internal/api/run.go, internal/changeset/body.go, body_test.go, internal/daemon/requests.go, internal/store/read_test.go; the rebase was aborted; answer question rebase-conflict-workflow-dispatch'. With several runs in flight, main always moves under a long run, so this will recur. Instead, ship-change should merge origin/main into the branch (never rebase or force-push) and give the conflicts to a fixer agent session with both sides' commits and change sets as context. The fixer resolves them, regenerates golden files, runs the touched packages' tests, commits with trailers and pushes; the run records the push and goes on. It asks the owner (or the arbiter in an autonomous run) only if the fixer reports a real design clash. Unstuck by hand: an agent merged main in the run's worktree, answered the question and resumed the run.
