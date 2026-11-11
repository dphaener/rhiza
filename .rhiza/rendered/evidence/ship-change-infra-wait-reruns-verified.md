# ship-change-infra-wait-reruns-verified

evidence · `01M75Y21E9A1D02C2F26AH4BDS` · actor [claude-code-b47efc9e](../actor/claude-code-b47efc9e.md)

## result

PASS

## statement

Verified at 17be6dbf (attempt 2). rhiza workflow check ship-change: valid (builtin), and the built-in's wait-ci-rerun until is 'rhiza ship check --settled --changed --rerun-infra' (internal/workflow/builtin/ship-change.toml:162, asserted in workflow_test.go:457). go test -trimpath ./internal/workflow/... ok. go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/run-ship-(infra-wait|rerun-fail)': both PASS. Read the diff: run-ship-infra-wait's fake gh now gives Actions run 42 (detailsUrl .../actions/runs/42/job/N), run view logs with no failing test, and turns green only after the second 'run rerun'; the script itself never reruns. A bin/rhiza shim only backdates red_at (the test clock is fixed) and logs --rerun-infra polls; the script asserts two 'run rerun 42 --failed' lines and that the second directly follows a 'rhiza ship check --settled --changed --rerun-infra' poll, so the wait issued it. Run status shows ci-triage fail -> wait-ci-rerun pass -> pr-green (visit 2) pass -> signoff. Attempt 1's finding (rerun simulated by the fake's queue) is fixed.

## Links

### Out

- verifies [ship-change-infra-wait-reruns](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-infra-wait-reruns.md)

## History

- `01M75R7FN9PWNCC5XFXWV3NVQR` Fix A run that fails on a CI infrastructure outage stays failed until a person retries it after the outage clears. Evidence: nebari plan run 01M75F9WMZTCBH5MRN8BZNPY7A (change set 01M75F9X20W9WXEXJ1RVR7WEVR, PR #96, changes only .rhiza/) failed pr-green twice on 2026-11-11 ('PR #96 at c062647: failed: test (1/5)..(5/5) after 2 reruns'); every failed job was GitHub runners unable to resolve the Ubuntu mirrors (360 'Temporary failure resolving archive.ubuntu.com' lines in run 43037903471, apt step, before any spec ran); nebari main CI 43032838034 failed the same way. The fix for 01M6W9VN (plan workflow reruns a CI flake) spends its two reruns within minutes, so a sustained outage exhausts them, and the watchdog (01M6W96C) recovers stalled or green-unmerged work but not a run that ended failed on infra. What unstuck it: the orchestrator reran the failed jobs and ran rhiza run resume --retry 01M75F9W at 15:29Z and 16:00Z, then merged nebari #97 (apt .deb cache) into the branch with ship rebase; it went green and merged. Wanted: when a CI failure is classified as infra (no failing test recorded, e.g. runner, network or package-mirror errors; ci-triage already distinguishes this), the run waits and retries with backoff over a longer window (e.g. up to a few hours) instead of ending, or the watchdog retries runs that failed at pr-green/rebased-green/landed-green on infra once a later CI run on main or another PR passes; only after that window does it fail and reach the inbox. Applies to plan and ship-change workflows.
