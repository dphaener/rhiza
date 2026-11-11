# pr-441-fix-a-run-that-fails-on-a-ci

evidence · `01M766KDXX5AACC754ZJNZNC6D` · actor [dphaener](../actor/dphaener.md)

## head

0232e60433f830f62b3ad3af5d59d169001839e2

## pr

https://github.com/dphaener/rhiza/pull/441

## result

10 checks passed: lint SUCCESS, go-install SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS

## statement

Pull request #441 (https://github.com/dphaener/rhiza/pull/441) passed CI at 0232e60

## status

verified

## Links

### Out

- verifies [ship-check-infra-rerun-backoff](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-infra-rerun-backoff.md)
- verifies [plan-ci-infra-wait](../capability/rhiza/orchestration/shipping/wait-green-pr/plan-ci-infra-wait.md)
- verifies [ship-change-infra-wait-reruns](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-infra-wait-reruns.md)
- verifies [ship-change-landed-infra-wait](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-landed-infra-wait.md)
- verifies [plan-run-reruns-ci-flakes](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-run-reruns-ci-flakes.md)
- verifies [ship-change-ci-fix-step](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-ci-fix-step.md)
- verifies [ship-change-conflict-fix-step](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-conflict-fix-step.md)

## History

- `01M75R7FN9PWNCC5XFXWV3NVQR` Fix A run that fails on a CI infrastructure outage stays failed until a person retries it after the outage clears. Evidence: nebari plan run 01M75F9WMZTCBH5MRN8BZNPY7A (change set 01M75F9X20W9WXEXJ1RVR7WEVR, PR #96, changes only .rhiza/) failed pr-green twice on 2026-11-11 ('PR #96 at c062647: failed: test (1/5)..(5/5) after 2 reruns'); every failed job was GitHub runners unable to resolve the Ubuntu mirrors (360 'Temporary failure resolving archive.ubuntu.com' lines in run 43037903471, apt step, before any spec ran); nebari main CI 43032838034 failed the same way. The fix for 01M6W9VN (plan workflow reruns a CI flake) spends its two reruns within minutes, so a sustained outage exhausts them, and the watchdog (01M6W96C) recovers stalled or green-unmerged work but not a run that ended failed on infra. What unstuck it: the orchestrator reran the failed jobs and ran rhiza run resume --retry 01M75F9W at 15:29Z and 16:00Z, then merged nebari #97 (apt .deb cache) into the branch with ship rebase; it went green and merged. Wanted: when a CI failure is classified as infra (no failing test recorded, e.g. runner, network or package-mirror errors; ci-triage already distinguishes this), the run waits and retries with backoff over a longer window (e.g. up to a few hours) instead of ending, or the watchdog retries runs that failed at pr-green/rebased-green/landed-green on infra once a later CI run on main or another PR passes; only after that window does it fail and reach the inbox. Applies to plan and ship-change workflows.
