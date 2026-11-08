# conflict-fix-design-clash-qa-1

evidence · `01M6XFWVMXGN3SJQKWB8TV0KTB` · actor [claude-code-b805a61d](../actor/claude-code-b805a61d.md)

## result

PASS

## statement

QA at a1136367: rhiza workflow describe ship-change shows conflict-fix and conflict-fix-landed prompts telling the fixer on a design clash to commit nothing, leave the merge in progress, ask one blocking question naming both change sets, the clashing files and behaviours, and never end the session with the merge in progress. Their until '! git rev-parse --verify --quiet MERGE_HEAD' run by hand in the scratch repo fails (exit 1) while the merge is in progress and holds once it is gone. go test -trimpath ./cmd/rhiza -run 'TestScripts/run-ship-conflict-fix$' ok (205s), covering the autonomous arbiter case, the non-autonomous blocked case and the done-without-commit failure.

## Links

### Out

- verifies [conflict-fix-design-clash](../capability/rhiza/orchestration/shipping/wait-green-pr/conflict-fix-design-clash.md)

## History

- `01M6WMCCRC0KSCWEV17VY8KVQP` Fix Ship-change stops on a rebase conflict and waits for the owner instead of resolving it: run 01M6W0Z45JXE (change set 01M6VY7G6HPJ, Workflow dispatch, PR #367) stopped at wait-pr on 2026-11-08 ~01:20Z with 'rebasing cs/workflow-dispatch onto origin/main conflicts in describe-golden.txtar, internal/api/change.go, internal/api/run.go, internal/changeset/body.go, body_test.go, internal/daemon/requests.go, internal/store/read_test.go; the rebase was aborted; answer question rebase-conflict-workflow-dispatch'. With several runs in flight, main always moves under a long run, so this will recur. Instead, ship-change should merge origin/main into the branch (never rebase or force-push) and give the conflicts to a fixer agent session with both sides' commits and change sets as context. The fixer resolves them, regenerates golden files, runs the touched packages' tests, commits with trailers and pushes; the run records the push and goes on. It asks the owner (or the arbiter in an autonomous run) only if the fixer reports a real design clash. Unstuck by hand: an agent merged main in the run's worktree, answered the question and resumed the run.
