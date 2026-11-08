# pr-383-fix-ship-change-stops-on-a

evidence · `01M6YJKSPQVCPSRCS202RAR9GM` · actor [dphaener](../actor/dphaener.md)

## head

6d6c4d1e7b0d4736e1fa378149423fc413a1990f

## pr

https://github.com/dphaener/rhiza/pull/383

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #383 (https://github.com/dphaener/rhiza/pull/383) passed CI at 6d6c4d1

## status

verified

## Links

### Out

- verifies [rebase-when-behind](../capability/rhiza/workflows/workflow-runs/rebase-when-behind.md)
- verifies [ship-change-conflict-fix-step](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-conflict-fix-step.md)
- verifies [conflict-fix-design-clash](../capability/rhiza/orchestration/shipping/wait-green-pr/conflict-fix-design-clash.md)
- verifies [ship-rebase-merges-only](../capability/rhiza/workflows/workflow-runs/rebase-when-behind/ship-rebase-merges-only.md)
- verifies [wait-merges-conflicting-pr](../capability/rhiza/orchestration/shipping/wait-green-pr/wait-merges-conflicting-pr.md)
- verifies [ship-change-ci-fix-step](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-ci-fix-step.md)

## History

- `01M6WMCCRC0KSCWEV17VY8KVQP` Fix Ship-change stops on a rebase conflict and waits for the owner instead of resolving it: run 01M6W0Z45JXE (change set 01M6VY7G6HPJ, Workflow dispatch, PR #367) stopped at wait-pr on 2026-11-08 ~01:20Z with 'rebasing cs/workflow-dispatch onto origin/main conflicts in describe-golden.txtar, internal/api/change.go, internal/api/run.go, internal/changeset/body.go, body_test.go, internal/daemon/requests.go, internal/store/read_test.go; the rebase was aborted; answer question rebase-conflict-workflow-dispatch'. With several runs in flight, main always moves under a long run, so this will recur. Instead, ship-change should merge origin/main into the branch (never rebase or force-push) and give the conflicts to a fixer agent session with both sides' commits and change sets as context. The fixer resolves them, regenerates golden files, runs the touched packages' tests, commits with trailers and pushes; the run records the push and goes on. It asks the owner (or the arbiter in an autonomous run) only if the fixer reports a real design clash. Unstuck by hand: an agent merged main in the run's worktree, answered the question and resumed the run.
