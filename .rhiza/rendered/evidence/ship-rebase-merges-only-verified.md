# ship-rebase-merges-only-verified

evidence · `01M6X6KYJXNJA5F6EAFSEVJA19` · actor [claude-code-3a8a616d](../actor/claude-code-3a8a616d.md)

## result

PASS

## statement

Verified at d1a0eabb by a verifier session (claude-code, Opus 5.5). Ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-ship-change-rebase|ship-check-settles-conflict|text-refs-engine-writes|describe-golden)$' -count=1: ok. run-ship-change-rebase.txtar covers each acceptance point: case 1 main moves with an unrelated change set's commit, the pushed head has 3 words in rev-list --parents (main and the old head as parents), the push is plain so origin's branch only fast-forwards, and the rendered index.md conflict is settled by rendering ('Settled 1 rendered file') and pushed; case 2 a code conflict fails the rebase step with exit 1, MERGE_HEAD present and git status showing unmerged paths, no rebase-conflict question, ship check --json shows merge_conflict with feature.txt, @shipping, main's commit subject and its Change: 01M5PWGEF85Y185JZ2ZGCNH40M and the branch's change; a merge committed by hand is then pushed by ship rebase and recorded as pushed in .git/rhiza/ship/*.json, with merge_conflict cleared. go test -trimpath ./internal/api/ -run TestShip -count=1: ok (TestShipRebaseClean asserts no force/lease/rebase calls; CodeConflict, FixerMerge, RenderedConflict, PushFails, Refusals). Ran the whole internal/api package with -skip Budget: every test passes; the package reports FAIL only for the pre-existing guard on .rhiza-root/usage.live.json writes, which this commit does not touch. A run without the skip also failed only TestReadBudget/TestReadBudgetApply timing at load average 16. Read the diff of ship_rebase.go: git merge --no-edit, plain git push origin HEAD:refs/heads/<branch>, merge left in progress (no merge --abort) on a code conflict. Nit: the comment on shipSettleConflict in ship.go still says a code conflict raises the rebase-conflict question.

## Links

### Out

- verifies [ship-rebase-merges-only](../capability/rhiza/workflows/workflow-runs/rebase-when-behind/ship-rebase-merges-only.md)

## History

- `01M6WMCCRC0KSCWEV17VY8KVQP` Fix Ship-change stops on a rebase conflict and waits for the owner instead of resolving it: run 01M6W0Z45JXE (change set 01M6VY7G6HPJ, Workflow dispatch, PR #367) stopped at wait-pr on 2026-11-08 ~01:20Z with 'rebasing cs/workflow-dispatch onto origin/main conflicts in describe-golden.txtar, internal/api/change.go, internal/api/run.go, internal/changeset/body.go, body_test.go, internal/daemon/requests.go, internal/store/read_test.go; the rebase was aborted; answer question rebase-conflict-workflow-dispatch'. With several runs in flight, main always moves under a long run, so this will recur. Instead, ship-change should merge origin/main into the branch (never rebase or force-push) and give the conflicts to a fixer agent session with both sides' commits and change sets as context. The fixer resolves them, regenerates golden files, runs the touched packages' tests, commits with trailers and pushes; the run records the push and goes on. It asks the owner (or the arbiter in an autonomous run) only if the fixer reports a real design clash. Unstuck by hand: an agent merged main in the run's worktree, answered the question and resumed the run.
