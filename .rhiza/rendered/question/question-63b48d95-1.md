# question-63b48d95-1

question · `01M75RCW4H03EE02G5Y5FFNEPJ`

## session

63b48d95-4c8c-4f2e-baf4-01a1cb842c3c

## statement

Change set 01M75R7 (CI infra outage leaves a run failed): ship-change already sends an infra-only red (no failing test) from ci-triage to wait-ci-rerun, but that wait only passes once someone else reruns CI or pushes, and plan/plan-milestone have no infra path at all (pr-green/rebased-green just end after --reruns 2). How should a run ride out a sustained outage? / Should ship-change's post-land green steps (landed-green, rebased-merge-green), which today send an infra red through conflict-triage-landed and rebase-merge-gate-failed to the end of the run, get the same infra wait in this change set?

## Links

### In

- answers [ci-infra-red-waits-in-run](../decision/ci-infra-red-waits-in-run.md)

## History

- `01M75R7FN9PWNCC5XFXWV3NVQR` Fix A run that fails on a CI infrastructure outage stays failed until a person retries it after the outage clears. Evidence: nebari plan run 01M75F9WMZTCBH5MRN8BZNPY7A (change set 01M75F9X20W9WXEXJ1RVR7WEVR, PR #96, changes only .rhiza/) failed pr-green twice on 2026-11-11 ('PR #96 at c062647: failed: test (1/5)..(5/5) after 2 reruns'); every failed job was GitHub runners unable to resolve the Ubuntu mirrors (360 'Temporary failure resolving archive.ubuntu.com' lines in run 43037903471, apt step, before any spec ran); nebari main CI 43032838034 failed the same way. The fix for 01M6W9VN (plan workflow reruns a CI flake) spends its two reruns within minutes, so a sustained outage exhausts them, and the watchdog (01M6W96C) recovers stalled or green-unmerged work but not a run that ended failed on infra. What unstuck it: the orchestrator reran the failed jobs and ran rhiza run resume --retry 01M75F9W at 15:29Z and 16:00Z, then merged nebari #97 (apt .deb cache) into the branch with ship rebase; it went green and merged. Wanted: when a CI failure is classified as infra (no failing test recorded, e.g. runner, network or package-mirror errors; ci-triage already distinguishes this), the run waits and retries with backoff over a longer window (e.g. up to a few hours) instead of ending, or the watchdog retries runs that failed at pr-green/rebased-green/landed-green on infra once a later CI run on main or another PR passes; only after that window does it fail and reach the inbox. Applies to plan and ship-change workflows.
