# ship-rebase-merges-only-qa-1

evidence · `01M6XFWTYYSQSV6EHBJ8VTQEMD` · actor [claude-code-b805a61d](../actor/claude-code-b805a61d.md)

## result

PASS

## statement

QA at a1136367 with the built binary in a scratch repo (bare local origin, branch cs/ship-it, all RHIZA_* env cleared): main moved with a conflicting code.txt commit carrying a Change: trailer; rhiza ship rebase exited 1 printing 'conflicts in: code.txt ... left in progress for a fixer', MERGE_HEAD present and code.txt UU, the ship state recorded merge_conflict with files, onto, main's commit subject and its Change: change set and the branch's change set; no question node was written (rhiza ls --kind question: none). After committing a resolving merge by hand, rhiza ship rebase pushed it ('Pushed 79b25ca'), origin's branch fast-forwarded to the local head, the state recorded pushed and dropped merge_conflict. No force or lease anywhere in the ship code's git calls.

## Links

### Out

- verifies [ship-rebase-merges-only](../capability/rhiza/workflows/workflow-runs/rebase-when-behind/ship-rebase-merges-only.md)

## History

- `01M6WMCCRC0KSCWEV17VY8KVQP` Fix Ship-change stops on a rebase conflict and waits for the owner instead of resolving it: run 01M6W0Z45JXE (change set 01M6VY7G6HPJ, Workflow dispatch, PR #367) stopped at wait-pr on 2026-11-08 ~01:20Z with 'rebasing cs/workflow-dispatch onto origin/main conflicts in describe-golden.txtar, internal/api/change.go, internal/api/run.go, internal/changeset/body.go, body_test.go, internal/daemon/requests.go, internal/store/read_test.go; the rebase was aborted; answer question rebase-conflict-workflow-dispatch'. With several runs in flight, main always moves under a long run, so this will recur. Instead, ship-change should merge origin/main into the branch (never rebase or force-push) and give the conflicts to a fixer agent session with both sides' commits and change sets as context. The fixer resolves them, regenerates golden files, runs the touched packages' tests, commits with trailers and pushes; the run records the push and goes on. It asks the owner (or the arbiter in an autonomous run) only if the fixer reports a real design clash. Unstuck by hand: an agent merged main in the run's worktree, answered the question and resumed the run.
