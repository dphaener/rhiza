# question-97efbc82-1

question · `01M6WMKQYEHGK1KKNGP2FQ4MH6`

## session

97efbc82-de6a-48f6-a9fa-01a1b9467030

## statement

Ship-change conflict fix (change set 01M6WMCC): how should the run bring origin/main into a change set's branch? Today `rhiza ship rebase` (steps rebase, rebase-merge) and `ship check --settled` (wait steps, on a CONFLICTING PR) rebase and force-push with a lease, and only merge when the branch already holds a merge. / When the merge conflicts in code, a fixer agent step resolves it. Who pushes the fixer's merge commit to the PR? / If the fixer finds a real design clash (both sides change the same behaviour in incompatible ways), what happens? / How many conflict-fix cycles may one ship-change run make before it fails?

## Links

### In

- answers [ship-merge-main-conflict-fixer](../decision/ship-merge-main-conflict-fixer.md)

## History

- `01M6WMCCRC0KSCWEV17VY8KVQP` Fix Ship-change stops on a rebase conflict and waits for the owner instead of resolving it: run 01M6W0Z45JXE (change set 01M6VY7G6HPJ, Workflow dispatch, PR #367) stopped at wait-pr on 2026-11-08 ~01:20Z with 'rebasing cs/workflow-dispatch onto origin/main conflicts in describe-golden.txtar, internal/api/change.go, internal/api/run.go, internal/changeset/body.go, body_test.go, internal/daemon/requests.go, internal/store/read_test.go; the rebase was aborted; answer question rebase-conflict-workflow-dispatch'. With several runs in flight, main always moves under a long run, so this will recur. Instead, ship-change should merge origin/main into the branch (never rebase or force-push) and give the conflicts to a fixer agent session with both sides' commits and change sets as context. The fixer resolves them, regenerates golden files, runs the touched packages' tests, commits with trailers and pushes; the run records the push and goes on. It asks the owner (or the arbiter in an autonomous run) only if the fixer reports a real design clash. Unstuck by hand: an agent merged main in the run's worktree, answered the question and resumed the run.
