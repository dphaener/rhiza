# rebase-when-behind

capability · `01M5X3GCHZRJ9YPV8BA37J1AN3` · parent [workflow-runs](../workflow-runs.md)

## acceptance

Its child leaves' tests: run-ship-change-rebase.txtar, ship-check-settles-conflict.txtar and run-ship-conflict-fix.txtar

## statement

Before landing and again before merging, ship-change brings origin/main into a PR that main has moved past by merging it, never by rebasing or force-pushing: conflicts only under .rhiza/rendered/ are settled by re-rendering; a conflict in code goes to a conflict-fix agent step that resolves it; after the merge the run pushes plainly and waits for CI again.

## Links

### Out

- decided-by [ship-change-holds-up](../../../../decision/ship-change-holds-up.md)

### In

- verifies [rebase-when-behind-tests](../../../../evidence/rebase-when-behind-tests.md)
- verifies [pr-383-fix-ship-change-stops-on-a](../../../../evidence/pr-383-fix-ship-change-stops-on-a.md)

## History

- `01M5X3GBH4ZQTZVXFEHR966747` ship-change that holds up: the agent step passes only when the agent says it is done, land renders with the branch's own build in Rhiza's repo, a PR that main has moved past is rebased, and a failed run can be resumed from the step that failed.
- `01M6WMCCRC0KSCWEV17VY8KVQP` Fix Ship-change stops on a rebase conflict and waits for the owner instead of resolving it: run 01M6W0Z45JXE (change set 01M6VY7G6HPJ, Workflow dispatch, PR #367) stopped at wait-pr on 2026-11-08 ~01:20Z with 'rebasing cs/workflow-dispatch onto origin/main conflicts in describe-golden.txtar, internal/api/change.go, internal/api/run.go, internal/changeset/body.go, body_test.go, internal/daemon/requests.go, internal/store/read_test.go; the rebase was aborted; answer question rebase-conflict-workflow-dispatch'. With several runs in flight, main always moves under a long run, so this will recur. Instead, ship-change should merge origin/main into the branch (never rebase or force-push) and give the conflicts to a fixer agent session with both sides' commits and change sets as context. The fixer resolves them, regenerates golden files, runs the touched packages' tests, commits with trailers and pushes; the run records the push and goes on. It asks the owner (or the arbiter in an autonomous run) only if the fixer reports a real design clash. Unstuck by hand: an agent merged main in the run's worktree, answered the question and resumed the run.
