# plan-ci-infra-wait

capability · `01M75RHSENQJ88C8SBDYFDFMZY` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

rhiza workflow check passes on plan and plan-milestone and describe shows pr-green and rebased-green routing through their infra triage and wait back to themselves; testscript plan-workflow-ci-flake-red (or a new plan-workflow-ci-infra-wait with a fake gh) drives a plan run whose CI fails with no failing test past --reruns 2 into the infra wait, which reruns, and the run merges once a rerun is green; a second run whose wait times out ends failed naming the infra failure; run with go test -trimpath ./internal/workflow/... and go test -trimpath ./cmd/rhiza -run 'TestScripts/plan-workflow-ci' only

## statement

In the plan and plan-milestone workflows, pr-green and rebased-green each get on.fail to their own infra triage step (rhiza ship check --infra, max_visits 4) whose on.pass is an infra wait step (until rhiza ship check --settled --changed --rerun-infra, timeout 6h, interval 30s, max_visits 3) passing back to that same green step, and their max_visits goes from 3 to 4 so a green step can be visited again after the infra wait, so a plan PR red on an infra outage reruns with backoff until CI settles and continues to rebase or merge; a red with a failing test fails the triage and ends the run as today, and the wait's timeout ends the run with a line naming the infra failure, which reaches the inbox

## Links

### Out

- depends-on [ship-check-infra-rerun-backoff](ship-check-infra-rerun-backoff.md)
- decided-by [ci-infra-red-waits-in-run](../../../../../decision/ci-infra-red-waits-in-run.md)
- supersedes [plan-run-reruns-ci-flakes](../../../workflows/built-in-workflows/plan-workflow/plan-run-reruns-ci-flakes.md)
- derived-from [ci-infra-outage-wait-plan](../../../../../source/ci-infra-outage-wait-plan.md)

### In

- verifies [plan-ci-infra-wait-verified](../../../../../evidence/plan-ci-infra-wait-verified.md)
- verifies [plan-ci-infra-wait-qa-1](../../../../../evidence/plan-ci-infra-wait-qa-1.md)
- verifies [plan-ci-infra-wait-verified-2](../../../../../evidence/plan-ci-infra-wait-verified-2.md)
- verifies [pr-441-fix-a-run-that-fails-on-a-ci](../../../../../evidence/pr-441-fix-a-run-that-fails-on-a-ci.md)

## History

- `01M75R7FN9PWNCC5XFXWV3NVQR` Fix A run that fails on a CI infrastructure outage stays failed until a person retries it after the outage clears. Evidence: nebari plan run 01M75F9WMZTCBH5MRN8BZNPY7A (change set 01M75F9X20W9WXEXJ1RVR7WEVR, PR #96, changes only .rhiza/) failed pr-green twice on 2026-11-11 ('PR #96 at c062647: failed: test (1/5)..(5/5) after 2 reruns'); every failed job was GitHub runners unable to resolve the Ubuntu mirrors (360 'Temporary failure resolving archive.ubuntu.com' lines in run 43037903471, apt step, before any spec ran); nebari main CI 43032838034 failed the same way. The fix for 01M6W9VN (plan workflow reruns a CI flake) spends its two reruns within minutes, so a sustained outage exhausts them, and the watchdog (01M6W96C) recovers stalled or green-unmerged work but not a run that ended failed on infra. What unstuck it: the orchestrator reran the failed jobs and ran rhiza run resume --retry 01M75F9W at 15:29Z and 16:00Z, then merged nebari #97 (apt .deb cache) into the branch with ship rebase; it went green and merged. Wanted: when a CI failure is classified as infra (no failing test recorded, e.g. runner, network or package-mirror errors; ci-triage already distinguishes this), the run waits and retries with backoff over a longer window (e.g. up to a few hours) instead of ending, or the watchdog retries runs that failed at pr-green/rebased-green/landed-green on infra once a later CI run on main or another PR passes; only after that window does it fail and reach the inbox. Applies to plan and ship-change workflows.
