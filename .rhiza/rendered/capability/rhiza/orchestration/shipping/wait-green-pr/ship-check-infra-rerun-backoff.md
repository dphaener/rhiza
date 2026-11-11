# ship-check-infra-rerun-backoff

capability · `01M75RHS7JP9D0WKRPAJAMTR38` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

Go tests in internal/api with a fake gh show --infra passing on a recorded infra red and failing on green, on a red with a failing test and with no red recorded; --settled --changed --rerun-infra rerunning the red head's failed jobs at 5m, 10m, 20m, 30m, 30m (driven by RHIZA_TEST_NOW, no real sleeps), staying not settled across infra-red attempts with the infra failure and next rerun named in its output, keeping the backoff across a fresh process from the state file, and settling on a green attempt, on an attempt red with a failing test and on a new head; run with go test -trimpath ./internal/api/... only

## statement

rhiza ship check gains an infra triage and an infra wait mode: --infra passes only when the head's latest failing ship check recorded a red head (red_head) with no failing test, and fails for a green head, a red with failing tests or any other failure; and --settled --changed --rerun-infra, while the red head stays red with no failing test, reruns its failed jobs itself with backoff (first after 5m, then doubling to a 30m cap, the last rerun's time and count kept in the ship state so the backoff survives polls and daemon restarts) and reports not settled with a line naming the infra failure and the next rerun's time; it settles when CI moved past the red head and finished green, finished red with a failing test, or a new head's CI finished, so one wait visit spans the whole outage

## Links

### Out

- decided-by [ci-infra-red-waits-in-run](../../../../../decision/ci-infra-red-waits-in-run.md)
- derived-from [ci-infra-outage-wait-plan](../../../../../source/ci-infra-outage-wait-plan.md)

### In

- depends-on [plan-ci-infra-wait](plan-ci-infra-wait.md)
- depends-on [ship-change-infra-wait-reruns](ship-change-infra-wait-reruns.md)
- depends-on [ship-change-landed-infra-wait](ship-change-landed-infra-wait.md)
- verifies [ship-check-infra-rerun-backoff-verified](../../../../../evidence/ship-check-infra-rerun-backoff-verified.md)
- verifies [ship-check-infra-rerun-backoff-qa-1](../../../../../evidence/ship-check-infra-rerun-backoff-qa-1.md)
- verifies [pr-441-fix-a-run-that-fails-on-a-ci](../../../../../evidence/pr-441-fix-a-run-that-fails-on-a-ci.md)

## History

- `01M75R7FN9PWNCC5XFXWV3NVQR` Fix A run that fails on a CI infrastructure outage stays failed until a person retries it after the outage clears. Evidence: nebari plan run 01M75F9WMZTCBH5MRN8BZNPY7A (change set 01M75F9X20W9WXEXJ1RVR7WEVR, PR #96, changes only .rhiza/) failed pr-green twice on 2026-11-11 ('PR #96 at c062647: failed: test (1/5)..(5/5) after 2 reruns'); every failed job was GitHub runners unable to resolve the Ubuntu mirrors (360 'Temporary failure resolving archive.ubuntu.com' lines in run 43037903471, apt step, before any spec ran); nebari main CI 43032838034 failed the same way. The fix for 01M6W9VN (plan workflow reruns a CI flake) spends its two reruns within minutes, so a sustained outage exhausts them, and the watchdog (01M6W96C) recovers stalled or green-unmerged work but not a run that ended failed on infra. What unstuck it: the orchestrator reran the failed jobs and ran rhiza run resume --retry 01M75F9W at 15:29Z and 16:00Z, then merged nebari #97 (apt .deb cache) into the branch with ship rebase; it went green and merged. Wanted: when a CI failure is classified as infra (no failing test recorded, e.g. runner, network or package-mirror errors; ci-triage already distinguishes this), the run waits and retries with backoff over a longer window (e.g. up to a few hours) instead of ending, or the watchdog retries runs that failed at pr-green/rebased-green/landed-green on infra once a later CI run on main or another PR passes; only after that window does it fail and reach the inbox. Applies to plan and ship-change workflows.
