# watchdog-files-bugs

capability · `01M6WJ2DW8TQBJPM5R59T3QGFK` · parent [daemon-watchdog](../daemon-watchdog.md)

## acceptance

go test -trimpath ./internal/api -run TestWatchdogBugs shows one bug request queued with the run ID, step and log lines, a second matching recovery queuing none while the first request is open, a CI flake queuing none, and fileBugs = false queuing none

## statement

For each recovery caused by a bug (a stall, a transient failure other than a CI flake, a PR its run left unmerged or without CI) the daemon queues one bug request with the evidence (run and request IDs, step, log lines), and does not queue another while an open request already names the same recovery kind and step; fileBugs = false (rhiza.daemon.watchdog.fileBugs, or rhiza.project.watchdog.fileBugs for one project) turns this off

## Links

### Out

- decided-by [watchdog-scope-and-defaults](../../../../../decision/watchdog-scope-and-defaults.md)

### In

- verifies [watchdog-files-bugs-verified](../../../../../evidence/watchdog-files-bugs-verified.md)
- verifies [watchdog-files-bugs-qa-1](../../../../../evidence/watchdog-files-bugs-qa-1.md)
- verifies [watchdog-files-bugs-qa-2](../../../../../evidence/watchdog-files-bugs-qa-2.md)
- verifies [pr-375-daemon-watchdog](../../../../../evidence/pr-375-daemon-watchdog.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
