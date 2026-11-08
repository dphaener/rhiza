# watchdog-stuck-requests

capability · `01M6WJ2DHPKP2B9JRXXB8NCGPX` · parent [daemon-watchdog](../daemon-watchdog.md)

## acceptance

go test -trimpath ./internal/api -run TestWatchdogRequests shows a request whose plan run is gone resumed, a request whose plan run failed on a transient signature retried, and a request failed for any other reason left for the owner

## statement

A plan request stuck planning (its plan run stalled, or ended with no live run while the request still says planning) or failed on a transient signature is recovered the same way as any run: its plan run is nudged, retried or resumed within the budget

## Links

### Out

- decided-by [watchdog-scope-and-defaults](../../../../../decision/watchdog-scope-and-defaults.md)

### In

- verifies [watchdog-stuck-requests-verified](../../../../../evidence/watchdog-stuck-requests-verified.md)
- verifies [pr-375-daemon-watchdog](../../../../../evidence/pr-375-daemon-watchdog.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
