# wait-merges-conflicting-pr-verified

evidence · `01M6X9WF3NXXJ6V71PXF736M0F` · actor [claude-code-bda14bfb](../actor/claude-code-bda14bfb.md)

## result

PASS

## statement

Verifier at e12f061d (attempt 2). Ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(ship-check-settles-conflict|describe-golden)$' -count=1 -v: both PASS. Read the script and diff: case 1 (rendered-only conflict, fake gh CONFLICTING until origin's branch holds main) has ship check --settled exit non-zero, origin/cs/ship-it holding origin/main and the ship record's pushed equal to the new head; case 2 (code conflict) has ship check --settled --json exit 0 with settled true, code.txt and the merge left in progress (MERGE_HEAD present, origin branch unchanged), no question, ship check --json exit=1 with passed false naming code.txt although the fake gh reports checks SUCCESS, and plain ship check failing with 'unresolved, conflicts in code.txt'. After the resolving merge is committed, plain ship check exits 0 and ship check --conflicts exits 1 ('is resolved'); after ship check --settled pushes, merge_conflict is dropped from the ship state, plain check exits 0 and --conflicts exits 1 ('no merge conflict is recorded'). That closes attempt 1's finding f1: the new --conflicts flag (ShipCheckRequest.Conflicts, runShipConflicts in internal/api/ship.go) reads only the state file and git, and refuses to combine with other modes. go test -trimpath ./internal/api -run Ship -v: every TestShip* PASS, among them ShipRebaseCodeConflict, ShipRebaseFixerMerge and ShipCheckTestFailures. The package reported FAIL only through its guard against writes into the enclosing ~/code root's state folder, an effect of running under the root and not of this change.

## Links

### Out

- verifies [wait-merges-conflicting-pr](../capability/rhiza/orchestration/shipping/wait-green-pr/wait-merges-conflicting-pr.md)

## History

- `01M6WMCCRC0KSCWEV17VY8KVQP` Fix Ship-change stops on a rebase conflict and waits for the owner instead of resolving it: run 01M6W0Z45JXE (change set 01M6VY7G6HPJ, Workflow dispatch, PR #367) stopped at wait-pr on 2026-11-08 ~01:20Z with 'rebasing cs/workflow-dispatch onto origin/main conflicts in describe-golden.txtar, internal/api/change.go, internal/api/run.go, internal/changeset/body.go, body_test.go, internal/daemon/requests.go, internal/store/read_test.go; the rebase was aborted; answer question rebase-conflict-workflow-dispatch'. With several runs in flight, main always moves under a long run, so this will recur. Instead, ship-change should merge origin/main into the branch (never rebase or force-push) and give the conflicts to a fixer agent session with both sides' commits and change sets as context. The fixer resolves them, regenerates golden files, runs the touched packages' tests, commits with trailers and pushes; the run records the push and goes on. It asks the owner (or the arbiter in an autonomous run) only if the fixer reports a real design clash. Unstuck by hand: an agent merged main in the run's worktree, answered the question and resumed the run.
