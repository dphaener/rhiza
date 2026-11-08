# plan-run-reruns-ci-flakes-verified

evidence · `01M6XGYK3PSTC54F2P9WWDE7VE` · actor [claude-code-a9c472db](../actor/claude-code-a9c472db.md)

## result

PASS

## statement

At 36bd98c8 I read the diff: plan.toml (and plan-milestone.toml) change only wait-pr/wait-rebased until and pr-green/rebased-green run to add --reruns 2 to rhiza ship check, with on.pass/on.fail/max_visits unchanged. TestPlanRerunsCIFlakes loads the built-in plan workflows, asserts those four steps carry --reruns 2 with pass routes, no fail route and max_visits 3, and asserts ship-change has at least one ship check step and none carries --reruns. go test -trimpath ./internal/workflow/... passed. go test -trimpath ./cmd/rhiza -run 'TestScripts/plan-workflow' passed all 8 scripts: plan-workflow-ci-flake (the fake gh fails test (4/5) until two 'run rerun 42 --failed' calls; the run succeeds, gh.log has exactly 2 reruns and 1 pr merge, merge and close-out pass), plan-workflow-ci-flake-red (stays red; the run fails at pr-green, no pr merge, status detail says 'failed: test (4/5) after 2 reruns' and 'the PR changes only .rhiza/'), and plan-workflow-pr still passes.

## Links

### Out

- verifies [plan-run-reruns-ci-flakes](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-run-reruns-ci-flakes.md)

## History

- `01M6WHTZC20SXN9ETH9ENEZ1PN` Fix The plan workflow ends a run on a CI flake: its pr-green and rebased-green steps fail the run as soon as a check is red, with no rerun. Evidence: nebari plan run 01M6VY5APJP9 (PR #78) passed review, gate, pr-green and one rebase, then rebased-green failed at 18230d2 on 'test (4/5); failing tests: rspec ./spec/system/garden_location_spec.rb[1:1:11]', a spec that fails at random and passes on rerun (nebari runs 42427557598, 42427449690). Unstuck by gh run rerun --failed plus rhiza run resume --retry. A plan PR changes only .rhiza/, so a red app test can't be its fault: on red, the plan workflow should rerun the failed jobs once (budgeted) and only then fail, naming the test as unrelated to the plan. Ship-change has ci-triage for this (#350); plan has none.
