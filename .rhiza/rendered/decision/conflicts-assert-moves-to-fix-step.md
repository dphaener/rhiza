# conflicts-assert-moves-to-fix-step

decision · `01M6WNWM30M9EFMCNGM55HQN9H`

## reasoning

The critic confirmed ShipCheckRequest has no --conflicts flag today and ship-change-conflict-fix-step introduces it while depending on wait-merges-conflicting-pr, so the assertion cannot run when wait-merges-conflicting-pr is built and verified. Moving the flag earlier would reverse the order and push the triage gate (its 3-fix cap counted from Review: trailers) into a leaf about the wait step, which is the fix step's concern. Moving only the assertion keeps each leaf testable when it is built: wait-merges-conflicting-pr still defines when a recorded conflict counts as resolved (read from git), which plain ship check already exercises, and the leaf that adds --conflicts tests it, including the resolved case it does not cover now (its acceptance only has 'no recorded conflict' and '3 fixes'). The change is the smaller edit, and the conflict-fix design (ship-merge-main-conflict-fixer) is unchanged.

## statement

Move the ship check --conflicts assertion out of wait-merges-conflicting-pr and into ship-change-conflict-fix-step; keep the leaf order as is. wait-merges-conflicting-pr drops 'and ship check --conflicts exits 1' from its acceptance and 'and ship check --conflicts fails with no conflict' from its statement: after a resolving merge it asserts only that plain ship check exits 0 on green checks (before and after the push) and that the next push drops the resolved record. ship-change-conflict-fix-step, which introduces the flag, adds to its internal/api Go test that ship check --conflicts fails once a recorded conflict's resolving merge is committed (MERGE_HEAD gone and HEAD holding the recorded origin/main commit), and to ship-check-settles-conflict.txtar's resolving-merge step that ship check --conflicts exits 1 before and after the head is pushed.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6wmcatnzdh5gh9egga6efbr-r3-f1](../question/review-01m6wmcatnzdh5gh9egga6efbr-r3-f1.md)

## History

- `01M6WMCCRC0KSCWEV17VY8KVQP` Fix Ship-change stops on a rebase conflict and waits for the owner instead of resolving it: run 01M6W0Z45JXE (change set 01M6VY7G6HPJ, Workflow dispatch, PR #367) stopped at wait-pr on 2026-11-08 ~01:20Z with 'rebasing cs/workflow-dispatch onto origin/main conflicts in describe-golden.txtar, internal/api/change.go, internal/api/run.go, internal/changeset/body.go, body_test.go, internal/daemon/requests.go, internal/store/read_test.go; the rebase was aborted; answer question rebase-conflict-workflow-dispatch'. With several runs in flight, main always moves under a long run, so this will recur. Instead, ship-change should merge origin/main into the branch (never rebase or force-push) and give the conflicts to a fixer agent session with both sides' commits and change sets as context. The fixer resolves them, regenerates golden files, runs the touched packages' tests, commits with trailers and pushes; the run records the push and goes on. It asks the owner (or the arbiter in an autonomous run) only if the fixer reports a real design clash. Unstuck by hand: an agent merged main in the run's worktree, answered the question and resumed the run.
