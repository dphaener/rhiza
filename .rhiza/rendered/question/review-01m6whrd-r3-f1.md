# review-01m6whrd-r3-f1

question · `01M6WK4B1BFDQVKA33C7PNE4KA`

## statement

Plan review round 3 left finding f1 (contradiction on watchdog-transient-failure): The model's daemon-resume-once (01M61R7MP3HGEB60GJV4M3CWQA) says the daemon leaves alone every failed run (the inbox lists it as failed) and never resumes a run interrupted again after its one resume. The watchdog leaves watchdog-transient-failure, watchdog-ci-rerun, watchdog-green-pr-unmerged, watchdog-stuck-requests and watchdog-stalled-step resume failed or re-stalled runs with --retry. The plan doesn't amend or narrow daemon-resume-once, so after it lands the two capabilities contradict each other. daemon-resume-once should say failed and given-up runs are left to the watchdog's policy. The critic confirmed it: daemon-resume-once (01M61R7MP3HGEB60GJV4M3CWQA, unchanged in the what-if view of this change set) says the daemon leaves alone every failed run and any run interrupted again after its one resume. watchdog-transient-failure resumes failed runs with --retry, and the other watchdog leaves resume given-up runs. The plan writes no op on daemon-resume-once and rhiza overlaps shows nothing, so once it lands the two capabilities contradict each other. How should it be settled?

## Links

### Out

- blocks [watchdog-transient-failure](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-transient-failure.md)

### In

- answers [daemon-resume-once-defers-to-watchdog](../decision/daemon-resume-once-defers-to-watchdog.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
