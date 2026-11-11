# ship-change-infra-wait-reruns

capability · `01M75RHSQEVEXEYYNDCVJ5BKHS` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

rhiza workflow check ship-change is valid and wait-ci-rerun's until carries --rerun-infra; testscript run-ship-infra-wait drives a ship-change run whose CI fails with no failing test into wait-ci-rerun, where the wait itself reruns the failed jobs (no rerun by the script) and the run passes pr-green to signoff once a rerun is green; run with go test -trimpath ./internal/workflow/... and go test -trimpath ./cmd/rhiza -run 'TestScripts/run-ship-(infra-wait|rerun-fail)' only

## statement

ship-change's wait-ci-rerun (reached from ci-triage when pr-green or rebased-green went red with no failing test) waits with rhiza ship check --settled --changed --rerun-infra, so it reruns the red head's failed jobs itself with backoff over its 6h window instead of waiting for a person or a push, and still passes on to pr-green once CI settles

## Links

### Out

- depends-on [ship-check-infra-rerun-backoff](ship-check-infra-rerun-backoff.md)
- decided-by [ci-infra-red-waits-in-run](../../../../../decision/ci-infra-red-waits-in-run.md)
- derived-from [ci-infra-outage-wait-plan](../../../../../source/ci-infra-outage-wait-plan.md)

### In

- verifies [ship-change-infra-wait-reruns-verified](../../../../../evidence/ship-change-infra-wait-reruns-verified.md)
- verifies [ship-change-infra-wait-reruns-qa-1](../../../../../evidence/ship-change-infra-wait-reruns-qa-1.md)
- verifies [pr-441-fix-a-run-that-fails-on-a-ci](../../../../../evidence/pr-441-fix-a-run-that-fails-on-a-ci.md)

## History

- `01M75R7FN9PWNCC5XFXWV3NVQR` Fix A run that fails on a CI infrastructure outage stays failed until a person retries it after the outage clears. Evidence: nebari plan run 01M75F9WMZTCBH5MRN8BZNPY7A (change set 01M75F9X20W9WXEXJ1RVR7WEVR, PR #96, changes only .rhiza/) failed pr-green twice on 2026-11-11 ('PR #96 at c062647: failed: test (1/5)..(5/5) after 2 reruns'); every failed job was GitHub runners unable to resolve the Ubuntu mirrors (360 'Temporary failure resolving archive.ubuntu.com' lines in run 43037903471, apt step, before any spec ran); nebari main CI 43032838034 failed the same way. The fix for 01M6W9VN (plan workflow reruns a CI flake) spends its two reruns within minutes, so a sustained outage exhausts them, and the watchdog (01M6W96C) recovers stalled or green-unmerged work but not a run that ended failed on infra. What unstuck it: the orchestrator reran the failed jobs and ran rhiza run resume --retry 01M75F9W at 15:29Z and 16:00Z, then merged nebari #97 (apt .deb cache) into the branch with ship rebase; it went green and merged. Wanted: when a CI failure is classified as infra (no failing test recorded, e.g. runner, network or package-mirror errors; ci-triage already distinguishes this), the run waits and retries with backoff over a longer window (e.g. up to a few hours) instead of ending, or the watchdog retries runs that failed at pr-green/rebased-green/landed-green on infra once a later CI run on main or another PR passes; only after that window does it fail and reach the inbox. Applies to plan and ship-change workflows.
