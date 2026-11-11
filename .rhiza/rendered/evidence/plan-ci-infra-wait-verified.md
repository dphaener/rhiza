# plan-ci-infra-wait-verified

evidence · `01M75XPN2J2W8XDR723JCRBW0N` · actor [claude-code-b554c945](../actor/claude-code-b554c945.md)

## result

PASS

## statement

At be216e8c (leaf commits d898ca6d and be216e8c): rhiza workflow check plan and plan-milestone are valid. workflow describe for both shows pr-green on.fail -> pr-green-infra -> pr-green-infra-wait -> pr-green, and the same loop for rebased-green. The infra wait steps use rhiza ship check --settled --changed --rerun-infra with a 6h timeout, a 30s interval and max_visits 3. go test -trimpath ./internal/workflow/... passes, including TestPlanBuiltin and TestPlanRerunsCIFlakes, which assert this routing. go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/plan-workflow-ci' passes all four scripts. plan-workflow-ci-infra-wait runs the built-in plan workflow with a fake gh. wait-pr spends both --reruns 2 (run rerun 42 --failed appears twice in gh.log). The run view log has no failing test. pr-green fails into pr-green-infra, which passes into pr-green-infra-wait. The wait passes once the fake's new job is green, and pr-green passes on visit 2. The run then goes through rebase, rebased-green, merge (pr merge appears once in gh.log) and close-out, and succeeds. plan-workflow-ci-infra-timeout ends the run on a pr-green-infra-wait timeout, and run status names an infra failure. plan-workflow-ci-flake-red's log now has a FAIL line, so a red with a failing test still ends the run. The diff fixes attempt 1's finding f1. One limit: under the fixed test clock the wait's own backoff rerun is never due, so the green comes from the fake's new job, not from a rerun the wait issued.

## Links

### Out

- verifies [plan-ci-infra-wait](../capability/rhiza/orchestration/shipping/wait-green-pr/plan-ci-infra-wait.md)

## History

- `01M75R7FN9PWNCC5XFXWV3NVQR` Fix A run that fails on a CI infrastructure outage stays failed until a person retries it after the outage clears. Evidence: nebari plan run 01M75F9WMZTCBH5MRN8BZNPY7A (change set 01M75F9X20W9WXEXJ1RVR7WEVR, PR #96, changes only .rhiza/) failed pr-green twice on 2026-11-11 ('PR #96 at c062647: failed: test (1/5)..(5/5) after 2 reruns'); every failed job was GitHub runners unable to resolve the Ubuntu mirrors (360 'Temporary failure resolving archive.ubuntu.com' lines in run 43037903471, apt step, before any spec ran); nebari main CI 43032838034 failed the same way. The fix for 01M6W9VN (plan workflow reruns a CI flake) spends its two reruns within minutes, so a sustained outage exhausts them, and the watchdog (01M6W96C) recovers stalled or green-unmerged work but not a run that ended failed on infra. What unstuck it: the orchestrator reran the failed jobs and ran rhiza run resume --retry 01M75F9W at 15:29Z and 16:00Z, then merged nebari #97 (apt .deb cache) into the branch with ship rebase; it went green and merged. Wanted: when a CI failure is classified as infra (no failing test recorded, e.g. runner, network or package-mirror errors; ci-triage already distinguishes this), the run waits and retries with backoff over a longer window (e.g. up to a few hours) instead of ending, or the watchdog retries runs that failed at pr-green/rebased-green/landed-green on infra once a later CI run on main or another PR passes; only after that window does it fail and reach the inbox. Applies to plan and ship-change workflows.
