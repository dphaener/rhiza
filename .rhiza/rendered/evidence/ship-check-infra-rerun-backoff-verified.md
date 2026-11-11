# ship-check-infra-rerun-backoff-verified

evidence · `01M75WAC8ZCM1SQ8Q11XZ9NM73` · actor [claude-code-f387f8bd](../actor/claude-code-f387f8bd.md)

## result

PASS

## statement

Verified commit 19f7def9 (verifier, attempt 1). Ran go test -trimpath ./internal/api/... -run 'TestShipCheckInfra|TestShipCheckRerunInfraBackoff' -v -count=1: both PASS. TestShipCheckInfra uses fake gh: --infra fails with no red recorded, passes on a recorded infra red (setup log, no failing test), fails on a red whose log has --- FAIL: TestGarden, fails on a green head, and refuses --infra with --settled. TestShipCheckRerunInfraBackoff drives ctx.Now (no sleeps): under --settled --changed --rerun-infra, no rerun 1s before each due time, and exactly one gh run rerun 42 --failed at 5m, 10m, 20m, 30m, 30m (5 total). Each infra-red attempt stays not settled, and Detail names the infra failure and the next rerun at. A later clock with the backoff read back from the state file does not rerun early. Settles and passes on a green attempt, on an attempt red with a failing test, and on a new head. Read the diff: RedAt/InfraReruns/InfraRerunAt are in the ship state, reset on a new red head and cleared on a recorded green. infraBackoff gives 5m doubling to a 30m cap. Other flag combinations are rejected. The whole internal/api package (go test -trimpath ./internal/api/ -v -timeout 900s) ran 935 passing tests including all 88 TestShip*, with no FAIL. It hit the timeout inside the unrelated read-timing test TestReadBudgetApply on a host at load avg ~11. describe-golden testscript passes with the new infra/rerun_infra flags.

## Links

### Out

- verifies [ship-check-infra-rerun-backoff](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-infra-rerun-backoff.md)

## History

- `01M75R7FN9PWNCC5XFXWV3NVQR` Fix A run that fails on a CI infrastructure outage stays failed until a person retries it after the outage clears. Evidence: nebari plan run 01M75F9WMZTCBH5MRN8BZNPY7A (change set 01M75F9X20W9WXEXJ1RVR7WEVR, PR #96, changes only .rhiza/) failed pr-green twice on 2026-11-11 ('PR #96 at c062647: failed: test (1/5)..(5/5) after 2 reruns'); every failed job was GitHub runners unable to resolve the Ubuntu mirrors (360 'Temporary failure resolving archive.ubuntu.com' lines in run 43037903471, apt step, before any spec ran); nebari main CI 43032838034 failed the same way. The fix for 01M6W9VN (plan workflow reruns a CI flake) spends its two reruns within minutes, so a sustained outage exhausts them, and the watchdog (01M6W96C) recovers stalled or green-unmerged work but not a run that ended failed on infra. What unstuck it: the orchestrator reran the failed jobs and ran rhiza run resume --retry 01M75F9W at 15:29Z and 16:00Z, then merged nebari #97 (apt .deb cache) into the branch with ship rebase; it went green and merged. Wanted: when a CI failure is classified as infra (no failing test recorded, e.g. runner, network or package-mirror errors; ci-triage already distinguishes this), the run waits and retries with backoff over a longer window (e.g. up to a few hours) instead of ending, or the watchdog retries runs that failed at pr-green/rebased-green/landed-green on infra once a later CI run on main or another PR passes; only after that window does it fail and reach the inbox. Applies to plan and ship-change workflows.
