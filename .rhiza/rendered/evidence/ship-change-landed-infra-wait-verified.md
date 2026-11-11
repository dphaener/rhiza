# ship-change-landed-infra-wait-verified

evidence · `01M75YR8XR3SW3M4TX5FK45G9G` · actor [claude-code-8ee8366a](../actor/claude-code-8ee8366a.md)

## result

PASS

## statement

Verified at 20947c77. go test -trimpath ./internal/workflow/... ok. go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/run-ship-(landed-)?infra-wait' PASS for run-ship-infra-wait and run-ship-landed-infra-wait. rhiza workflow check ship-change is valid (builtin). workflow describe shows rebase-merge-gate-failed on.fail -> landed-infra-triage (rhiza ship check --infra) -> landed-infra-wait -> landed-green; the built-in wait runs until 'rhiza ship check --settled --changed --rerun-infra' with timeout 6h, interval 30s and max_visits 3; gate-fix-merge stays on.pass. Read the diff of run-ship-landed-infra-wait.txtar. The fake gh's check is Actions run 42 with logs naming no failing test. wait-landed spends one rerun and the head is still red. landed-green, conflict-triage-landed and rebase-merge-gate-failed then fail through to landed-infra-triage and landed-infra-wait. The wait's own 'gh run rerun 42 --failed' comes right after a --rerun-infra poll, and the script never reruns. A bin/rhiza shim moves red_at back so that rerun is due. The head turns green, and the run goes on through gate-rebase-merge, rebase-merge, a squash merge pinned to the green head, wait-main, main-green and close-out (asserted in run status and gh.log). This fixes both earlier findings: merge and close-out are now reached, and the rerun now comes from the wait.

## Links

### Out

- verifies [ship-change-landed-infra-wait](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-landed-infra-wait.md)

## History

- `01M75R7FN9PWNCC5XFXWV3NVQR` Fix A run that fails on a CI infrastructure outage stays failed until a person retries it after the outage clears. Evidence: nebari plan run 01M75F9WMZTCBH5MRN8BZNPY7A (change set 01M75F9X20W9WXEXJ1RVR7WEVR, PR #96, changes only .rhiza/) failed pr-green twice on 2026-11-11 ('PR #96 at c062647: failed: test (1/5)..(5/5) after 2 reruns'); every failed job was GitHub runners unable to resolve the Ubuntu mirrors (360 'Temporary failure resolving archive.ubuntu.com' lines in run 43037903471, apt step, before any spec ran); nebari main CI 43032838034 failed the same way. The fix for 01M6W9VN (plan workflow reruns a CI flake) spends its two reruns within minutes, so a sustained outage exhausts them, and the watchdog (01M6W96C) recovers stalled or green-unmerged work but not a run that ended failed on infra. What unstuck it: the orchestrator reran the failed jobs and ran rhiza run resume --retry 01M75F9W at 15:29Z and 16:00Z, then merged nebari #97 (apt .deb cache) into the branch with ship rebase; it went green and merged. Wanted: when a CI failure is classified as infra (no failing test recorded, e.g. runner, network or package-mirror errors; ci-triage already distinguishes this), the run waits and retries with backoff over a longer window (e.g. up to a few hours) instead of ending, or the watchdog retries runs that failed at pr-green/rebased-green/landed-green on infra once a later CI run on main or another PR passes; only after that window does it fail and reach the inbox. Applies to plan and ship-change workflows.
