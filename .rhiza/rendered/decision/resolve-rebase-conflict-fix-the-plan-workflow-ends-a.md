# resolve-rebase-conflict-fix-the-plan-workflow-ends-a

decision · `01M6XQV5H1FZZKJ72V4XTC5QKZ`

## reasoning

The conflict was purely additive: both sides appended new tests at the same spot after TestShipCheck's conflict test, and neither changes the other's. Main now merges main into a branch rather than rebasing (no force-push), so a merge commit matches the current ship flow. ship.go auto-merged; go test -trimpath ./internal/api -run 'Ship|Plan|Watchdog', internal/workflow tests and the describe-golden and plan-workflow-ci-flake testscripts pass on the merged tree.

## statement

Merge origin/main into cs/fix-the-plan-workflow-ends-a instead of rebasing, and resolve internal/api/ship_test.go by keeping both sides' tests: the branch's TestShipCheckFailureDetailNotes (with its merge-base/diff git constants) and main's withRunEnv helper with TestShipCheckNamesRun, TestShipLandNamesRun and TestShipRebaseNamesRun; re-render .rhiza/rendered/index.md

## wording

(assumed, would have escalated)

## Links

### Out

- answers [rebase-conflict-fix-the-plan-workflow-ends-a](../question/rebase-conflict-fix-the-plan-workflow-ends-a.md)

## History

- `01M6WHTZC20SXN9ETH9ENEZ1PN` Fix The plan workflow ends a run on a CI flake: its pr-green and rebased-green steps fail the run as soon as a check is red, with no rerun. Evidence: nebari plan run 01M6VY5APJP9 (PR #78) passed review, gate, pr-green and one rebase, then rebased-green failed at 18230d2 on 'test (4/5); failing tests: rspec ./spec/system/garden_location_spec.rb[1:1:11]', a spec that fails at random and passes on rerun (nebari runs 42427557598, 42427449690). Unstuck by gh run rerun --failed plus rhiza run resume --retry. A plan PR changes only .rhiza/, so a red app test can't be its fault: on red, the plan workflow should rerun the failed jobs once (budgeted) and only then fail, naming the test as unrelated to the plan. Ship-change has ci-triage for this (#350); plan has none.
