# plan-run-reruns-ci-flakes-qa-1

evidence · `01M6XJKBS4X4C6QYNXA1GEFK1F` · actor [claude-code-14b9e333](../actor/claude-code-14b9e333.md)

## result

PASS

## statement

QA cycle 1: read internal/workflow/builtin/plan.toml (and plan-milestone.toml): wait-pr and wait-rebased until 'rhiza ship check --reruns 2 --settled', pr-green and rebased-green run 'rhiza ship check --reruns 2', routes and max_visits unchanged; ship-change untouched. Drove those exact commands by hand with the built binary against a fake gh: a red head was rerun twice and then failed with the reruns spent and the .rhiza-only note, which is the detail the failing step carries. Ran go test -trimpath ./cmd/rhiza -run 'TestScripts/plan-workflow' (ci-flake, ci-flake-red, pr) and ./internal/workflow/....

## Links

### Out

- verifies [plan-run-reruns-ci-flakes](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-run-reruns-ci-flakes.md)

## History

- `01M6WHTZC20SXN9ETH9ENEZ1PN` Fix The plan workflow ends a run on a CI flake: its pr-green and rebased-green steps fail the run as soon as a check is red, with no rerun. Evidence: nebari plan run 01M6VY5APJP9 (PR #78) passed review, gate, pr-green and one rebase, then rebased-green failed at 18230d2 on 'test (4/5); failing tests: rspec ./spec/system/garden_location_spec.rb[1:1:11]', a spec that fails at random and passes on rerun (nebari runs 42427557598, 42427449690). Unstuck by gh run rerun --failed plus rhiza run resume --retry. A plan PR changes only .rhiza/, so a red app test can't be its fault: on red, the plan workflow should rerun the failed jobs once (budgeted) and only then fail, naming the test as unrelated to the plan. Ship-change has ci-triage for this (#350); plan has none.
