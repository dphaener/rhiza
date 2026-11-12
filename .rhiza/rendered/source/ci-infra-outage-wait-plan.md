# ci-infra-outage-wait-plan

source · `01M75RGXB7N4XFBQ0MWDR39Z90`

## quotes

```json
[]
```

## raw_sha256

436a911ce5eff1f34ab32b720b5f0f2ec60fd4e359dcd3fbeda1a47d94d580c7

## session

63b48d95-4c8c-4f2e-baf4-01a1cb842c3c

## summary

Autonomous plan run for change set 01M75R7 (a run failed on a CI infra outage stays failed). The planner read the existing paths: ship-change's ci-triage already routes an infra-only red (no failing test) to wait-ci-rerun (6h), but that wait only passes once someone else reruns CI or pushes; plan and plan-milestone end at pr-green/rebased-green after --reruns 2; ship-change's landed-green and rebased-merge-green send an infra red through conflict-triage-landed and rebase-merge-gate-failed to the end of the run; the watchdog reruns CI once per head. Two questions (approach, and whether post-land green steps are in scope) went to the arbiter, which ruled (Decision ci-infra-red-waits-in-run): ride out the outage inside the run with an infra wait that reruns the red head's failed jobs itself with backoff over about 6h, returning to the green step, for every green step in plan, plan-milestone and ship-change including post-land; not through the watchdog. The planner then split the work into four leaves: the ship check rerun-with-backoff wait mode and infra triage flag, plan and plan-milestone infra paths, ship-change's pre-land wait-ci-rerun using it, and ship-change's post-land infra path.

## title

CI infra outage wait plan, autonomous run with the arbiter

## Links

### In

- derived-from [ci-infra-red-waits-in-run](../decision/ci-infra-red-waits-in-run.md)
- derived-from [ship-check-infra-rerun-backoff](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-infra-rerun-backoff.md)
- derived-from [plan-ci-infra-wait](../capability/rhiza/orchestration/shipping/wait-green-pr/plan-ci-infra-wait.md)
- derived-from [ship-change-infra-wait-reruns](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-infra-wait-reruns.md)
- derived-from [ship-change-landed-infra-wait](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-landed-infra-wait.md)

## History

- `01M75R7FN9PWNCC5XFXWV3NVQR` Fix A run that fails on a CI infrastructure outage stays failed until a person retries it after the outage clears. Evidence: nebari plan run 01M75F9WMZTCBH5MRN8BZNPY7A (change set 01M75F9X20W9WXEXJ1RVR7WEVR, PR #96, changes only .rhiza/) failed pr-green twice on 2026-11-11 ('PR #96 at c062647: failed: test (1/5)..(5/5) after 2 reruns'); every failed job was GitHub runners unable to resolve the Ubuntu mirrors (360 'Temporary failure resolving archive.ubuntu.com' lines in run 43037903471, apt step, before any spec ran); nebari main CI 43032838034 failed the same way. The fix for 01M6W9VN (plan workflow reruns a CI flake) spends its two reruns within minutes, so a sustained outage exhausts them, and the watchdog (01M6W96C) recovers stalled or green-unmerged work but not a run that ended failed on infra. What unstuck it: the orchestrator reran the failed jobs and ran rhiza run resume --retry 01M75F9W at 15:29Z and 16:00Z, then merged nebari #97 (apt .deb cache) into the branch with ship rebase; it went green and merged. Wanted: when a CI failure is classified as infra (no failing test recorded, e.g. runner, network or package-mirror errors; ci-triage already distinguishes this), the run waits and retries with backoff over a longer window (e.g. up to a few hours) instead of ending, or the watchdog retries runs that failed at pr-green/rebased-green/landed-green on infra once a later CI run on main or another PR passes; only after that window does it fail and reach the inbox. Applies to plan and ship-change workflows.
