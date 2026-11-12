# ship-change-landed-infra-wait-verified-3

evidence · `01M761SRQ3VM3632WGTYYJC3Z7` · actor [claude-code-fa95a565](../actor/claude-code-fa95a565.md)

## result

PASS

## statement

Verifier at fa388f90 (attempt 3). Built the binary in a scratch dir: `rhiza workflow check ship-change` is valid (builtin), and `workflow describe ship-change` shows landed-green and rebased-merge-green on.fail -> conflict-triage-landed -> rebase-merge-gate-failed, whose on.fail -> landed-infra-triage, on.pass -> landed-infra-wait, on.pass -> landed-green. The built-in's landed-infra-wait is `until rhiza ship check --settled --changed --rerun-infra`, 6h timeout, 30s interval, max_visits 3. The earlier finding is fixed: the built-in landed-green now has max_visits 4 (1 + the wait's 3 returns). TestShipChangeLocalGate asserts that cap, and TestShipChangeLandedInfraWaitScript pins the script's copy of landed-green, landed-infra-triage and landed-infra-wait (caps, commands, routes) to the built-in. `go test -trimpath ./internal/workflow/...` ok. `go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/run-ship-(landed-)?infra-wait'`: both PASS. run-ship-landed-infra-wait drives a landed run with a fake gh: landed-green fails, goes through conflict-triage-landed and rebase-merge-gate-failed to landed-infra-triage (pass), then landed-infra-wait, whose own --rerun-infra poll issues `run rerun 42 --failed` before the green. It then passes at landed-green visit 2, goes through gate-rebase-merge, rebase-merge, rebased-merge-green, merge (`pr merge 7 --squash --match-head-commit`), main-green and close-out, and the run succeeds. Read the diff: the code change is the cap, its test and the script comment; HEAD fa388f90 is an empty review commit.

## Links

### Out

- verifies [ship-change-landed-infra-wait](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-landed-infra-wait.md)

## History

- `01M75R7FN9PWNCC5XFXWV3NVQR` Fix A run that fails on a CI infrastructure outage stays failed until a person retries it after the outage clears. Evidence: nebari plan run 01M75F9WMZTCBH5MRN8BZNPY7A (change set 01M75F9X20W9WXEXJ1RVR7WEVR, PR #96, changes only .rhiza/) failed pr-green twice on 2026-11-11 ('PR #96 at c062647: failed: test (1/5)..(5/5) after 2 reruns'); every failed job was GitHub runners unable to resolve the Ubuntu mirrors (360 'Temporary failure resolving archive.ubuntu.com' lines in run 43037903471, apt step, before any spec ran); nebari main CI 43032838034 failed the same way. The fix for 01M6W9VN (plan workflow reruns a CI flake) spends its two reruns within minutes, so a sustained outage exhausts them, and the watchdog (01M6W96C) recovers stalled or green-unmerged work but not a run that ended failed on infra. What unstuck it: the orchestrator reran the failed jobs and ran rhiza run resume --retry 01M75F9W at 15:29Z and 16:00Z, then merged nebari #97 (apt .deb cache) into the branch with ship rebase; it went green and merged. Wanted: when a CI failure is classified as infra (no failing test recorded, e.g. runner, network or package-mirror errors; ci-triage already distinguishes this), the run waits and retries with backoff over a longer window (e.g. up to a few hours) instead of ending, or the watchdog retries runs that failed at pr-green/rebased-green/landed-green on infra once a later CI run on main or another PR passes; only after that window does it fail and reach the inbox. Applies to plan and ship-change workflows.
