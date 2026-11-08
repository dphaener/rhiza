# plan-run-reruns-ci-flakes

capability · `01M6WJ6ZN7V3CXMEZYFA58HAYV` · parent [plan-workflow](../plan-workflow.md)

## acceptance

A Go test in internal/workflow loads the built-in plan workflow and asserts wait-pr, wait-rebased (until) and pr-green, rebased-green (run) each call rhiza ship check with --reruns 2, with their on.pass, on.fail and max_visits unchanged, and that ship-change's ship check steps carry no --reruns; a testscript plan-workflow-ci-flake with a fake gh whose check fails twice on the PR head and passes on the second rerun shows the plan run merging and closing out with two 'run rerun' calls in gh.log, and a case red after both reruns fails the run with a detail saying the PR changes only .rhiza/; plan-workflow-pr still passes; run with go test -trimpath ./internal/workflow/... and go test -trimpath ./cmd/rhiza -run 'TestScripts/plan-workflow' only

## statement

The built-in plan workflow passes --reruns 2 to rhiza ship check in wait-pr, pr-green, wait-rebased and rebased-green, so a plan PR head that goes red gets its failed Actions jobs rerun twice (three CI attempts) before the step fails and ends the run with ship check's failure detail (the failing tests, the reruns spent and, for a .rhiza-only diff, that the PR changes only .rhiza/); no step, route or max_visits changes, and ship-change keeps the default single rerun

## Links

### Out

- depends-on [ship-check-rerun-budget](../../../orchestration/shipping/wait-green-pr/ship-check-rerun-budget.md)
- depends-on [ship-check-failure-names-reruns-and-rhiza-only-diff](../../../orchestration/shipping/wait-green-pr/ship-check-failure-names-reruns-and-rhiza-only-diff.md)
- decided-by [plan-ci-flake-reruns](../../../../../decision/plan-ci-flake-reruns.md)

### In

- verifies [plan-run-reruns-ci-flakes-verified](../../../../../evidence/plan-run-reruns-ci-flakes-verified.md)
- verifies [plan-run-reruns-ci-flakes-qa-1](../../../../../evidence/plan-run-reruns-ci-flakes-qa-1.md)
- verifies [pr-386-fix-the-plan-workflow-ends-a](../../../../../evidence/pr-386-fix-the-plan-workflow-ends-a.md)

## History

- `01M6WHTZC20SXN9ETH9ENEZ1PN` Fix The plan workflow ends a run on a CI flake: its pr-green and rebased-green steps fail the run as soon as a check is red, with no rerun. Evidence: nebari plan run 01M6VY5APJP9 (PR #78) passed review, gate, pr-green and one rebase, then rebased-green failed at 18230d2 on 'test (4/5); failing tests: rspec ./spec/system/garden_location_spec.rb[1:1:11]', a spec that fails at random and passes on rerun (nebari runs 42427557598, 42427449690). Unstuck by gh run rerun --failed plus rhiza run resume --retry. A plan PR changes only .rhiza/, so a red app test can't be its fault: on red, the plan workflow should rerun the failed jobs once (budgeted) and only then fail, naming the test as unrelated to the plan. Ship-change has ci-triage for this (#350); plan has none.
