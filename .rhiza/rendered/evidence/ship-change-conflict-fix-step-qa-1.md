# ship-change-conflict-fix-step-qa-1

evidence · `01M6XFWVCBXJYMRGT6ZQXZCDPM` · actor [claude-code-b805a61d](../actor/claude-code-b805a61d.md)

## result

PASS

## statement

QA at a1136367 with the built binary: rhiza workflow check ship-change is valid and workflow describe shows rebase on.fail rebase-triage, rebase-triage on.pass conflict-fix, pr-green and rebased-green on.fail conflict-triage, conflict-triage on.pass conflict-fix / on.fail ci-triage, landed-green and rebase-merge on.fail conflict-triage-landed, conflict-triage-landed on.pass conflict-fix-landed on.pass rebase-merge. In a scratch repo, three conflicts each fixed by a merge commit with the same Review: <run> trailer and pushed by ship rebase: ship check --conflicts (RHIZA_RUN set) exited 0 on the 2nd and 3rd conflicts and exited 1 on the 4th naming 'the run has already made 3 conflict fixes'. go test -trimpath ./cmd/rhiza -run 'TestScripts/run-ship-conflict-fix$' ok (205s).

## Links

### Out

- verifies [ship-change-conflict-fix-step](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-conflict-fix-step.md)

## History

- `01M6WMCCRC0KSCWEV17VY8KVQP` Fix Ship-change stops on a rebase conflict and waits for the owner instead of resolving it: run 01M6W0Z45JXE (change set 01M6VY7G6HPJ, Workflow dispatch, PR #367) stopped at wait-pr on 2026-11-08 ~01:20Z with 'rebasing cs/workflow-dispatch onto origin/main conflicts in describe-golden.txtar, internal/api/change.go, internal/api/run.go, internal/changeset/body.go, body_test.go, internal/daemon/requests.go, internal/store/read_test.go; the rebase was aborted; answer question rebase-conflict-workflow-dispatch'. With several runs in flight, main always moves under a long run, so this will recur. Instead, ship-change should merge origin/main into the branch (never rebase or force-push) and give the conflicts to a fixer agent session with both sides' commits and change sets as context. The fixer resolves them, regenerates golden files, runs the touched packages' tests, commits with trailers and pushes; the run records the push and goes on. It asks the owner (or the arbiter in an autonomous run) only if the fixer reports a real design clash. Unstuck by hand: an agent merged main in the run's worktree, answered the question and resumed the run.
