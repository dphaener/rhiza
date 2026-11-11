# ship-change-infra-wait-reruns-qa-1

evidence · `01M7607AV6549FWN0XWEMPVERK` · actor [claude-code-e6d490f4](../actor/claude-code-e6d490f4.md)

## result

PASS

## statement

QA at c7b756f9 with the built binary: rhiza workflow check ship-change is valid; describe shows ci-triage on.fail -> wait-ci-rerun, whose until is rhiza ship check --settled --changed --rerun-infra; that command, run by hand against a fake gh with an infra red, reran run 42's failed jobs itself on a 5/10/20/30/30m backoff and settled once a rerun went green. Testscripts run-ship-infra-wait and run-ship-rerun-fail pass.

## Links

### Out

- verifies [ship-change-infra-wait-reruns](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-infra-wait-reruns.md)

## History

- `01M75R7FN9PWNCC5XFXWV3NVQR` Fix A run that fails on a CI infrastructure outage stays failed until a person retries it after the outage clears. Evidence: nebari plan run 01M75F9WMZTCBH5MRN8BZNPY7A (change set 01M75F9X20W9WXEXJ1RVR7WEVR, PR #96, changes only .rhiza/) failed pr-green twice on 2026-11-11 ('PR #96 at c062647: failed: test (1/5)..(5/5) after 2 reruns'); every failed job was GitHub runners unable to resolve the Ubuntu mirrors (360 'Temporary failure resolving archive.ubuntu.com' lines in run 43037903471, apt step, before any spec ran); nebari main CI 43032838034 failed the same way. The fix for 01M6W9VN (plan workflow reruns a CI flake) spends its two reruns within minutes, so a sustained outage exhausts them, and the watchdog (01M6W96C) recovers stalled or green-unmerged work but not a run that ended failed on infra. What unstuck it: the orchestrator reran the failed jobs and ran rhiza run resume --retry 01M75F9W at 15:29Z and 16:00Z, then merged nebari #97 (apt .deb cache) into the branch with ship rebase; it went green and merged. Wanted: when a CI failure is classified as infra (no failing test recorded, e.g. runner, network or package-mirror errors; ci-triage already distinguishes this), the run waits and retries with backoff over a longer window (e.g. up to a few hours) instead of ending, or the watchdog retries runs that failed at pr-green/rebased-green/landed-green on infra once a later CI run on main or another PR passes; only after that window does it fail and reach the inbox. Applies to plan and ship-change workflows.
