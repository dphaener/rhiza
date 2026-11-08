# watchdog-ci-rerun-qa-2

evidence · `01M6X118XDZ88BTXDMBWAH1AJM` · actor [claude-code-d3420702](../actor/claude-code-d3420702.md)

## result

PASS

## statement

An autonomous run whose command step 'rhiza ship check --settled' failed (fake gh answering pr view with a CLOSED PR) was polled by four rhiza daemon --once processes: the fake gh log showed one 'gh pr view' in total, none on later polls, and no gh run list or rerun (47741152: the PR is read once and a closed PR is not re-read while the journal is unchanged).

## Links

### Out

- verifies [watchdog-ci-rerun](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-ci-rerun.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
