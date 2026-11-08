# watchdog-green-pr-unmerged

capability · `01M6WJ2D976WM21C1X3YVMCHY8` · parent [daemon-watchdog](../daemon-watchdog.md)

## acceptance

go test -trimpath ./internal/api -run TestWatchdogGreenUnmerged, with a fake gh, shows a resume --retry for an autonomous ship that ended failed at its merge step, nothing for a live run at its merge step, and an inbox item with no resume and no merge call both for a run that ended succeeded with its PR unmerged and when the ship stage needs the owner

## statement

When a run's PR is green but the run never merged it, the daemon acts by the run's state: a run still live (at its wait or merge step) is left alone; a run that ended failed at or before its merge step is resumed with --retry so the workflow's own merge step merges it, never merging outside a workflow; a run that ended succeeded with its PR unmerged (run resume refuses --retry on it), or whose workflow's merge needs the owner, gets an inbox item instead

## Links

### Out

- decided-by [watchdog-scope-and-defaults](../../../../../decision/watchdog-scope-and-defaults.md)

### In

- verifies [watchdog-green-pr-unmerged-verified](../../../../../evidence/watchdog-green-pr-unmerged-verified.md)
- verifies [pr-375-daemon-watchdog](../../../../../evidence/pr-375-daemon-watchdog.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
