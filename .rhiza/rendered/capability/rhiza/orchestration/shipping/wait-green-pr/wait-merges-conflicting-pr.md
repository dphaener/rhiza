# wait-merges-conflicting-pr

capability · `01M6WMSB7FA6TA1SMY3YW9AC1Q` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

ship-check-settles-conflict.txtar with a fake gh whose pr view reports mergeable CONFLICTING until a new head is pushed: with origin/main conflicting only under .rhiza/rendered, ship check --settled exits non-zero, origin's branch fast-forwards to a merge holding origin/main and the ship record's pushed is the new head; with a conflict in a code file it exits 0, the merge is still in progress, ship check --json lists the conflicts and no question is raised, and plain ship check then exits 1 naming the conflicted files although the fake gh reports the head's checks green; after a resolving merge is committed on top, plain ship check exits 0 on green checks and ship check --conflicts exits 1, before and after the head is pushed; run with go test -trimpath ./cmd/rhiza -run 'TestScripts/ship-check-settles-conflict'

## statement

rhiza ship check --settled reads the PR's mergeable state and, when GitHub reports it CONFLICTING, merges origin/main as rhiza ship rebase does: a clean or rendered-only merge is pushed plainly, recorded as pushed and reported not settled; a conflict in code leaves the merge in progress, records the conflicts in the ship state and reports settled (exit 0) with the conflict in its detail; it never exits 75 and raises no question on a conflict. Plain rhiza ship check (the green steps pr-green, rebased-green, landed-green, rebased-merge-green) fails, whatever the head's checks say, while the ship state records an unresolved merge conflict, naming the conflicted files in its detail, so the green step after the wait routes the run to the conflict triage and never on to signoff or land with a merge in progress. A recorded conflict counts as unresolved only while the run worktree still has the merge in progress (MERGE_HEAD present) or HEAD does not yet contain the origin/main commit the record names; once the fixer has committed the merge it counts as resolved, read from git with no write needed, so plain ship check passes on green checks and ship check --conflicts fails with no conflict, and the next ship check --settled, ship rebase or open-pr that pushes the head drops the resolved record from the ship state

## Links

### Out

- depends-on [ship-rebase-merges-only](../../../workflows/workflow-runs/rebase-when-behind/ship-rebase-merges-only.md)
- decided-by [ship-merge-main-conflict-fixer](../../../../../decision/ship-merge-main-conflict-fixer.md)

### In

- depends-on [ship-change-conflict-fix-step](ship-change-conflict-fix-step.md)
- blocks [review-01m6wmcatnzdh5gh9egga6efbr-r3-f1](../../../../../question/review-01m6wmcatnzdh5gh9egga6efbr-r3-f1.md)
- verifies [wait-merges-conflicting-pr-verified](../../../../../evidence/wait-merges-conflicting-pr-verified.md)
- verifies [wait-merges-conflicting-pr-qa-1](../../../../../evidence/wait-merges-conflicting-pr-qa-1.md)
- verifies [pr-383-fix-ship-change-stops-on-a](../../../../../evidence/pr-383-fix-ship-change-stops-on-a.md)

## History

- `01M6WMCCRC0KSCWEV17VY8KVQP` Fix Ship-change stops on a rebase conflict and waits for the owner instead of resolving it: run 01M6W0Z45JXE (change set 01M6VY7G6HPJ, Workflow dispatch, PR #367) stopped at wait-pr on 2026-11-08 ~01:20Z with 'rebasing cs/workflow-dispatch onto origin/main conflicts in describe-golden.txtar, internal/api/change.go, internal/api/run.go, internal/changeset/body.go, body_test.go, internal/daemon/requests.go, internal/store/read_test.go; the rebase was aborted; answer question rebase-conflict-workflow-dispatch'. With several runs in flight, main always moves under a long run, so this will recur. Instead, ship-change should merge origin/main into the branch (never rebase or force-push) and give the conflicts to a fixer agent session with both sides' commits and change sets as context. The fixer resolves them, regenerates golden files, runs the touched packages' tests, commits with trailers and pushes; the run records the push and goes on. It asks the owner (or the arbiter in an autonomous run) only if the fixer reports a real design clash. Unstuck by hand: an agent merged main in the run's worktree, answered the question and resumed the run.
