# wait-merges-conflicting-pr-qa-1

evidence · `01M6XFWV49PJMX1SEYYBFC6JWF` · actor [claude-code-b805a61d](../actor/claude-code-b805a61d.md)

## result

PASS

## statement

QA at a1136367 with the built binary and a fake gh reporting green checks and mergeable CONFLICTING while origin's branch lacks main: with the merge left in progress plain rhiza ship check exited 1 ('failed: merging origin/main into cs/ship-it is unresolved, conflicts in code.txt') though checks were green; after the resolving merge was committed plain ship check exited 0 and ship check --conflicts exited 1 ('is resolved'). With main moved cleanly, ship check --settled merged, pushed a two-parent merge and exited 1 (not settled); with main moved conflicting, ship check --settled --json exited 0 with settled true, conflicts [code.txt] and the merge left in progress, and no question was raised.

## Links

### Out

- verifies [wait-merges-conflicting-pr](../capability/rhiza/orchestration/shipping/wait-green-pr/wait-merges-conflicting-pr.md)

## History

- `01M6WMCCRC0KSCWEV17VY8KVQP` Fix Ship-change stops on a rebase conflict and waits for the owner instead of resolving it: run 01M6W0Z45JXE (change set 01M6VY7G6HPJ, Workflow dispatch, PR #367) stopped at wait-pr on 2026-11-08 ~01:20Z with 'rebasing cs/workflow-dispatch onto origin/main conflicts in describe-golden.txtar, internal/api/change.go, internal/api/run.go, internal/changeset/body.go, body_test.go, internal/daemon/requests.go, internal/store/read_test.go; the rebase was aborted; answer question rebase-conflict-workflow-dispatch'. With several runs in flight, main always moves under a long run, so this will recur. Instead, ship-change should merge origin/main into the branch (never rebase or force-push) and give the conflicts to a fixer agent session with both sides' commits and change sets as context. The fixer resolves them, regenerates golden files, runs the touched packages' tests, commits with trailers and pushes; the run records the push and goes on. It asks the owner (or the arbiter in an autonomous run) only if the fixer reports a real design clash. Unstuck by hand: an agent merged main in the run's worktree, answered the question and resumed the run.
