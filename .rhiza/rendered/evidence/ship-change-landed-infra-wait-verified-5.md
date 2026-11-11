# ship-change-landed-infra-wait-verified-5

evidence · `01M763J1C4HMCGND880QGF1EDR` · actor [claude-code-84ea1a58](../actor/claude-code-84ea1a58.md)

## result

PASS

## statement

Verified at ecec6df5. A build of this commit's rhiza ran `workflow check ship-change`: valid. `workflow describe ship-change` shows conflict-triage-landed on.fail -> rebase-merge-gate-failed, whose on.fail -> landed-infra-triage (rhiza ship check --infra), whose on.pass -> landed-infra-wait (until rhiza ship check --settled --changed --rerun-infra, 6h, 30s, max_visits 3), whose on.pass -> landed-green (max_visits 4). go test -trimpath ./internal/workflow/... passes, and TestShipChangeLocalGate pins landed-green's cap at 4 (>= 1 + the wait's cap). go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/run-ship-(landed-)?infra-wait' passes both scripts. GOTMPDIR was moved to the scratchpad and RHIZA_ROOT/CHANGE/RUN were unset, because the run's GOTMPDIR under the repo's .git makes rhiza init find the real repo. run-ship-landed-infra-wait drives a landed run through landed-green fail, conflict-triage-landed, rebase-merge-gate-failed, landed-infra-triage and landed-infra-wait, then landed-green visit 2 pass, gate-rebase-merge, rebase-merge, wait-rebased-merge, rebased-merge-green, merge, wait-main, main-green and close-out. It asserts the wait's own `run rerun 42 --failed` follows a --rerun-infra poll and comes before a `pr merge 7 --squash --match-head-commit <landed head>`, and that the branch is removed. The model's landed nodes ship-change-ci-fix-step and ship-change-conflict-fix-step now describe the post-land infra wait, and ship-change-landed-infra-wait supersedes both. Nit, not part of the acceptance: the comment above [step.conflict-triage-landed] (ship-change.toml:271) still says a red CI 'still ends the run there'.

## Links

### Out

- verifies [ship-change-landed-infra-wait](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-landed-infra-wait.md)

## History

- `01M75R7FN9PWNCC5XFXWV3NVQR` Fix A run that fails on a CI infrastructure outage stays failed until a person retries it after the outage clears. Evidence: nebari plan run 01M75F9WMZTCBH5MRN8BZNPY7A (change set 01M75F9X20W9WXEXJ1RVR7WEVR, PR #96, changes only .rhiza/) failed pr-green twice on 2026-11-11 ('PR #96 at c062647: failed: test (1/5)..(5/5) after 2 reruns'); every failed job was GitHub runners unable to resolve the Ubuntu mirrors (360 'Temporary failure resolving archive.ubuntu.com' lines in run 43037903471, apt step, before any spec ran); nebari main CI 43032838034 failed the same way. The fix for 01M6W9VN (plan workflow reruns a CI flake) spends its two reruns within minutes, so a sustained outage exhausts them, and the watchdog (01M6W96C) recovers stalled or green-unmerged work but not a run that ended failed on infra. What unstuck it: the orchestrator reran the failed jobs and ran rhiza run resume --retry 01M75F9W at 15:29Z and 16:00Z, then merged nebari #97 (apt .deb cache) into the branch with ship rebase; it went green and merged. Wanted: when a CI failure is classified as infra (no failing test recorded, e.g. runner, network or package-mirror errors; ci-triage already distinguishes this), the run waits and retries with backoff over a longer window (e.g. up to a few hours) instead of ending, or the watchdog retries runs that failed at pr-green/rebased-green/landed-green on infra once a later CI run on main or another PR passes; only after that window does it fail and reach the inbox. Applies to plan and ship-change workflows.
