# ship-check-infra-rerun-backoff-qa-1

evidence · `01M7607ADF4QZ3Q20V8W2206F7` · actor [claude-code-e6d490f4](../actor/claude-code-e6d490f4.md)

## result

PASS

## statement

QA at c7b756f9 with the built binary (go build -trimpath ./cmd/rhiza) in a scratch repo with a bare origin and a fake gh (Actions run 42, apt-mirror log, no failing test); clock moved by hand through RHIZA_TEST_NOW, each call a fresh process. ship check --infra: exit 1 with no red recorded; after --settled spent the head's one rerun and ship check recorded the red, --infra exit 0 and --test-failures exit 1. --settled --changed --rerun-infra: at red_at+4m not settled, naming the infra failure and next rerun at red_at+5m, no rerun; past it it ran gh run rerun 42 --failed (rerun 1) with next at +10m; reruns then fell at 10m, 20m, 30m, 30m after the previous one (state file red_at, infra_reruns 5, infra_rerun_at), with each new failed job ID (10..13) still not settled. An IN_PROGRESS attempt was not settled; a SUCCESS attempt settled (exit 0) and plain ship check passed and cleared red_at/infra_* from the state. On a new head: an attempt red with --- FAIL: TestGarden in its log settled the wait (exit 0), and --infra then failed with 'the red head has failing tests'. A push of another head mid-wait settled past the red head (new head's CI got its own rerun). A new red head reset red_at. Scoped tests TestShipCheckInfra and TestShipCheckRerunInfraBackoff pass. Nits: the rerun note omits the run ID ('reran the failed jobs of run (rerun 1)', ship.go:971), and a passing --infra prints only '.'.

## Links

### Out

- verifies [ship-check-infra-rerun-backoff](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-infra-rerun-backoff.md)

## History

- `01M75R7FN9PWNCC5XFXWV3NVQR` Fix A run that fails on a CI infrastructure outage stays failed until a person retries it after the outage clears. Evidence: nebari plan run 01M75F9WMZTCBH5MRN8BZNPY7A (change set 01M75F9X20W9WXEXJ1RVR7WEVR, PR #96, changes only .rhiza/) failed pr-green twice on 2026-11-11 ('PR #96 at c062647: failed: test (1/5)..(5/5) after 2 reruns'); every failed job was GitHub runners unable to resolve the Ubuntu mirrors (360 'Temporary failure resolving archive.ubuntu.com' lines in run 43037903471, apt step, before any spec ran); nebari main CI 43032838034 failed the same way. The fix for 01M6W9VN (plan workflow reruns a CI flake) spends its two reruns within minutes, so a sustained outage exhausts them, and the watchdog (01M6W96C) recovers stalled or green-unmerged work but not a run that ended failed on infra. What unstuck it: the orchestrator reran the failed jobs and ran rhiza run resume --retry 01M75F9W at 15:29Z and 16:00Z, then merged nebari #97 (apt .deb cache) into the branch with ship rebase; it went green and merged. Wanted: when a CI failure is classified as infra (no failing test recorded, e.g. runner, network or package-mirror errors; ci-triage already distinguishes this), the run waits and retries with backoff over a longer window (e.g. up to a few hours) instead of ending, or the watchdog retries runs that failed at pr-green/rebased-green/landed-green on infra once a later CI run on main or another PR passes; only after that window does it fail and reach the inbox. Applies to plan and ship-change workflows.
