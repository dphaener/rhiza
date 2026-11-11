# ci-infra-red-waits-in-run

decision · `01M75RFG35Z6HW5XA7HP2YNRTC`

## reasoning

The change set's intent asks that an infra-classified red wait and retry with backoff over a longer window (up to a few hours) instead of ending, for plan and ship-change, and names pr-green, rebased-green and landed-green explicitly, so the post-land green steps belong in this change set; leaving them out would keep the exact failure mode the intent reports on the landed path. Of the two mechanisms the intent offers, the in-run wait is the smaller, more direct one: ship-change already has the shape (ci-triage routes a no-failing-test red to wait-ci-rerun with a 6h timeout and ship check --settled --changed), and its only gap is that it waits for someone else to rerun, so the fix is for the wait to rerun failed jobs itself with backoff. The watchdog path (watchdogCIRerun in internal/api/daemon_watchdog.go) reruns once per head and leaves a second red to the owner, and resuming an ended run with --retry is the recovery the orchestrator had to do by hand; extending it would make an outage still end and restart runs and spread retry state across the daemon memory. Reusing ci-triage's existing infra classification (no failing test recorded) keeps one definition of infra across workflows, and a bounded window keeps a permanently broken CI reaching the inbox.

## statement

A run rides out a CI infrastructure outage inside the run, not through the watchdog: whenever a green step's red head has no failing test recorded (ship check --test-failures fails), the workflow goes to an infra wait step that itself reruns the head's failed jobs with backoff (minutes at first, growing to tens of minutes) and passes as soon as CI moves past the red head and settles, back to the green step it came from; only once a bounded window (about 6h, the existing wait-ci-rerun timeout) expires does the run fail, naming the infra failure, and reach the inbox. A red with failing tests keeps its current path (ci-fix, or ending the run). This applies to every green step in this change set: plan's and plan-milestone's pr-green and rebased-green, ship-change's pr-green (wait-ci-rerun gains the reruns) and ship-change's post-land landed-green and rebased-merge-green, whose infra red goes to the same wait instead of through conflict-triage-landed and rebase-merge-gate-failed to the end of the run.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-63b48d95-1](../question/question-63b48d95-1.md)
- derived-from [ci-infra-outage-wait-plan](../source/ci-infra-outage-wait-plan.md)

### In

- decided-by [ship-check-infra-rerun-backoff](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-infra-rerun-backoff.md)
- decided-by [plan-ci-infra-wait](../capability/rhiza/orchestration/shipping/wait-green-pr/plan-ci-infra-wait.md)
- decided-by [ship-change-infra-wait-reruns](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-infra-wait-reruns.md)
- decided-by [ship-change-landed-infra-wait](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-landed-infra-wait.md)

## History

- `01M75R7FN9PWNCC5XFXWV3NVQR` Fix A run that fails on a CI infrastructure outage stays failed until a person retries it after the outage clears. Evidence: nebari plan run 01M75F9WMZTCBH5MRN8BZNPY7A (change set 01M75F9X20W9WXEXJ1RVR7WEVR, PR #96, changes only .rhiza/) failed pr-green twice on 2026-11-11 ('PR #96 at c062647: failed: test (1/5)..(5/5) after 2 reruns'); every failed job was GitHub runners unable to resolve the Ubuntu mirrors (360 'Temporary failure resolving archive.ubuntu.com' lines in run 43037903471, apt step, before any spec ran); nebari main CI 43032838034 failed the same way. The fix for 01M6W9VN (plan workflow reruns a CI flake) spends its two reruns within minutes, so a sustained outage exhausts them, and the watchdog (01M6W96C) recovers stalled or green-unmerged work but not a run that ended failed on infra. What unstuck it: the orchestrator reran the failed jobs and ran rhiza run resume --retry 01M75F9W at 15:29Z and 16:00Z, then merged nebari #97 (apt .deb cache) into the branch with ship rebase; it went green and merged. Wanted: when a CI failure is classified as infra (no failing test recorded, e.g. runner, network or package-mirror errors; ci-triage already distinguishes this), the run waits and retries with backoff over a longer window (e.g. up to a few hours) instead of ending, or the watchdog retries runs that failed at pr-green/rebased-green/landed-green on infra once a later CI run on main or another PR passes; only after that window does it fail and reach the inbox. Applies to plan and ship-change workflows.
