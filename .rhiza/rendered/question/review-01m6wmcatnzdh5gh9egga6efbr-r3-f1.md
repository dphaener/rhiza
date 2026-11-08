# review-01m6wmcatnzdh5gh9egga6efbr-r3-f1

question · `01M6WNV511KX736205MYN2MVQN`

## statement

Plan review round 3 left finding f1 (missing-acceptance on wait-merges-conflicting-pr): Its acceptance (and statement) check that `ship check --conflicts` exits 1 after the resolving merge, but the --conflicts flag is introduced by ship-change-conflict-fix-step, which depends-on wait-merges-conflicting-pr and so is built after it. When this leaf is built and verified the flag does not exist, so that part of its acceptance cannot run. Either move the --conflicts assertion to ship-change-conflict-fix-step's acceptance, or move the flag into this leaf (or ship-rebase-merges-only) and reverse the ordering. The critic confirmed it: internal/api/ship.go ShipCheckRequest has no --conflicts flag today (only settled, main, pr-open, test-failures); ship-change-conflict-fix-step introduces it and depends-on wait-merges-conflicting-pr, so it is built later. wait-merges-conflicting-pr's acceptance asserts 'ship check --conflicts exits 1' after the resolving merge, which cannot run when that leaf is built and verified. How should it be settled?

## Links

### Out

- blocks [wait-merges-conflicting-pr](../capability/rhiza/orchestration/shipping/wait-green-pr/wait-merges-conflicting-pr.md)

### In

- answers [conflicts-assert-moves-to-fix-step](../decision/conflicts-assert-moves-to-fix-step.md)

## History

- `01M6WMCCRC0KSCWEV17VY8KVQP` Fix Ship-change stops on a rebase conflict and waits for the owner instead of resolving it: run 01M6W0Z45JXE (change set 01M6VY7G6HPJ, Workflow dispatch, PR #367) stopped at wait-pr on 2026-11-08 ~01:20Z with 'rebasing cs/workflow-dispatch onto origin/main conflicts in describe-golden.txtar, internal/api/change.go, internal/api/run.go, internal/changeset/body.go, body_test.go, internal/daemon/requests.go, internal/store/read_test.go; the rebase was aborted; answer question rebase-conflict-workflow-dispatch'. With several runs in flight, main always moves under a long run, so this will recur. Instead, ship-change should merge origin/main into the branch (never rebase or force-push) and give the conflicts to a fixer agent session with both sides' commits and change sets as context. The fixer resolves them, regenerates golden files, runs the touched packages' tests, commits with trailers and pushes; the run records the push and goes on. It asks the owner (or the arbiter in an autonomous run) only if the fixer reports a real design clash. Unstuck by hand: an agent merged main in the run's worktree, answered the question and resumed the run.
