# ship-change-landed-infra-wait

capability · `01M75RHSWYX0S5ZRZ3NEZTC8W8` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

rhiza workflow check ship-change is valid and describe shows rebase-merge-gate-failed's on.fail reaching the post-land infra triage and wait, passing to landed-green; a testscript (run-ship-infra-wait extended, or a new run-ship-landed-infra-wait with a fake gh) drives a landed ship-change run whose landed head's CI fails with no failing test into that wait, which reruns, and the run merges and closes out once a rerun is green; run with go test -trimpath ./internal/workflow/... and go test -trimpath ./cmd/rhiza -run 'TestScripts/run-ship-(landed-)?infra-wait' only

## statement

After land, an infra red on ship-change's landed-green or rebased-merge-green no longer ends the run: rebase-merge-gate-failed gets on.fail to an infra triage step (rhiza ship check --infra) whose on.pass is an infra wait step (until rhiza ship check --settled --changed --rerun-infra, timeout 6h, interval 30s, max_visits 3) passing to landed-green, whose max_visits goes from the default 1 to 4 (one entry from wait-landed plus the wait's 3 returns), which re-checks the head and goes on through gate-rebase-merge and rebase-merge to merge; a conflict, a failing gate, a red with a failing test or any other failure keeps its path today, and the wait's timeout ends the run naming the infra failure

## Links

### Out

- depends-on [ship-check-infra-rerun-backoff](ship-check-infra-rerun-backoff.md)
- decided-by [ci-infra-red-waits-in-run](../../../../../decision/ci-infra-red-waits-in-run.md)
- supersedes [ship-change-ci-fix-step](ship-change-ci-fix-step.md)
- supersedes [ship-change-conflict-fix-step](ship-change-conflict-fix-step.md)
- derived-from [ci-infra-outage-wait-plan](../../../../../source/ci-infra-outage-wait-plan.md)

### In

- verifies [ship-change-landed-infra-wait-verified](../../../../../evidence/ship-change-landed-infra-wait-verified.md)
- verifies [ship-change-landed-infra-wait-qa-1](../../../../../evidence/ship-change-landed-infra-wait-qa-1.md)
- verifies [ship-change-landed-infra-wait-verified-3](../../../../../evidence/ship-change-landed-infra-wait-verified-3.md)
- verifies [ship-change-landed-infra-wait-verified-5](../../../../../evidence/ship-change-landed-infra-wait-verified-5.md)
- verifies [pr-441-fix-a-run-that-fails-on-a-ci](../../../../../evidence/pr-441-fix-a-run-that-fails-on-a-ci.md)

## History

- `01M75R7FN9PWNCC5XFXWV3NVQR` Fix A run that fails on a CI infrastructure outage stays failed until a person retries it after the outage clears. Evidence: nebari plan run 01M75F9WMZTCBH5MRN8BZNPY7A (change set 01M75F9X20W9WXEXJ1RVR7WEVR, PR #96, changes only .rhiza/) failed pr-green twice on 2026-11-11 ('PR #96 at c062647: failed: test (1/5)..(5/5) after 2 reruns'); every failed job was GitHub runners unable to resolve the Ubuntu mirrors (360 'Temporary failure resolving archive.ubuntu.com' lines in run 43037903471, apt step, before any spec ran); nebari main CI 43032838034 failed the same way. The fix for 01M6W9VN (plan workflow reruns a CI flake) spends its two reruns within minutes, so a sustained outage exhausts them, and the watchdog (01M6W96C) recovers stalled or green-unmerged work but not a run that ended failed on infra. What unstuck it: the orchestrator reran the failed jobs and ran rhiza run resume --retry 01M75F9W at 15:29Z and 16:00Z, then merged nebari #97 (apt .deb cache) into the branch with ship rebase; it went green and merged. Wanted: when a CI failure is classified as infra (no failing test recorded, e.g. runner, network or package-mirror errors; ci-triage already distinguishes this), the run waits and retries with backoff over a longer window (e.g. up to a few hours) instead of ending, or the watchdog retries runs that failed at pr-green/rebased-green/landed-green on infra once a later CI run on main or another PR passes; only after that window does it fail and reach the inbox. Applies to plan and ship-change workflows.
