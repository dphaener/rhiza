# plan-run-reruns-ci-flakes

capability · `01M6WJ6ZN7V3CXMEZYFA58HAYV` · parent [plan-workflow](../plan-workflow.md)

## acceptance

A Go test in internal/workflow loads the built-in plan workflow and asserts wait-pr, wait-rebased (until) and pr-green, rebased-green (run) each call rhiza ship check with --reruns 2, with their on.pass and max_visits as plan-ci-infra-wait sets them and their on.fail to their infra triage, and that ship-change's ship check steps carry no --reruns; a testscript plan-workflow-ci-flake with a fake gh whose check fails twice on the PR head and passes on the second rerun shows the plan run merging and closing out with two 'run rerun' calls in gh.log, and a case red with a failing test after both reruns fails the run with a detail saying the PR changes only .rhiza/; plan-workflow-pr still passes; run with go test -trimpath ./internal/workflow/... and go test -trimpath ./cmd/rhiza -run 'TestScripts/plan-workflow' only

## statement

The built-in plan workflow passes --reruns 2 to rhiza ship check in wait-pr, pr-green, wait-rebased and rebased-green, so a plan PR head that goes red gets its failed Actions jobs rerun twice (three CI attempts) before the step fails with ship check's failure detail (the failing tests, the reruns spent and, for a .rhiza-only diff, that the PR changes only .rhiza/); a red with a failing test then ends the run, while a red with no failing test goes to the infra triage and wait of plan-ci-infra-wait; ship-change keeps the default single rerun

## Links

### Out

- depends-on [ship-check-rerun-budget](../../../orchestration/shipping/wait-green-pr/ship-check-rerun-budget.md)
- depends-on [ship-check-failure-names-reruns-and-rhiza-only-diff](../../../orchestration/shipping/wait-green-pr/ship-check-failure-names-reruns-and-rhiza-only-diff.md)
- decided-by [plan-ci-flake-reruns](../../../../../decision/plan-ci-flake-reruns.md)

### In

- supersedes [plan-ci-infra-wait](../../../orchestration/shipping/wait-green-pr/plan-ci-infra-wait.md)
- verifies [plan-run-reruns-ci-flakes-verified](../../../../../evidence/plan-run-reruns-ci-flakes-verified.md)
- verifies [plan-run-reruns-ci-flakes-qa-1](../../../../../evidence/plan-run-reruns-ci-flakes-qa-1.md)
- verifies [pr-386-fix-the-plan-workflow-ends-a](../../../../../evidence/pr-386-fix-the-plan-workflow-ends-a.md)
- verifies [pr-441-fix-a-run-that-fails-on-a-ci](../../../../../evidence/pr-441-fix-a-run-that-fails-on-a-ci.md)

## History

- `01M6WHTZC20SXN9ETH9ENEZ1PN` Fix The plan workflow ends a run on a CI flake: its pr-green and rebased-green steps fail the run as soon as a check is red, with no rerun. Evidence: nebari plan run 01M6VY5APJP9 (PR #78) passed review, gate, pr-green and one rebase, then rebased-green failed at 18230d2 on 'test (4/5); failing tests: rspec ./spec/system/garden_location_spec.rb[1:1:11]', a spec that fails at random and passes on rerun (nebari runs 42427557598, 42427449690). Unstuck by gh run rerun --failed plus rhiza run resume --retry. A plan PR changes only .rhiza/, so a red app test can't be its fault: on red, the plan workflow should rerun the failed jobs once (budgeted) and only then fail, naming the test as unrelated to the plan. Ship-change has ci-triage for this (#350); plan has none.
- `01M75R7FN9PWNCC5XFXWV3NVQR` Fix A run that fails on a CI infrastructure outage stays failed until a person retries it after the outage clears. Evidence: nebari plan run 01M75F9WMZTCBH5MRN8BZNPY7A (change set 01M75F9X20W9WXEXJ1RVR7WEVR, PR #96, changes only .rhiza/) failed pr-green twice on 2026-11-11 ('PR #96 at c062647: failed: test (1/5)..(5/5) after 2 reruns'); every failed job was GitHub runners unable to resolve the Ubuntu mirrors (360 'Temporary failure resolving archive.ubuntu.com' lines in run 43037903471, apt step, before any spec ran); nebari main CI 43032838034 failed the same way. The fix for 01M6W9VN (plan workflow reruns a CI flake) spends its two reruns within minutes, so a sustained outage exhausts them, and the watchdog (01M6W96C) recovers stalled or green-unmerged work but not a run that ended failed on infra. What unstuck it: the orchestrator reran the failed jobs and ran rhiza run resume --retry 01M75F9W at 15:29Z and 16:00Z, then merged nebari #97 (apt .deb cache) into the branch with ship rebase; it went green and merged. Wanted: when a CI failure is classified as infra (no failing test recorded, e.g. runner, network or package-mirror errors; ci-triage already distinguishes this), the run waits and retries with backoff over a longer window (e.g. up to a few hours) instead of ending, or the watchdog retries runs that failed at pr-green/rebased-green/landed-green on infra once a later CI run on main or another PR passes; only after that window does it fail and reach the inbox. Applies to plan and ship-change workflows.
