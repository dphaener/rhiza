# watchdog-stalled-step-qa-2

evidence · `01M6X118E4Q6VP9MND02HCN0C8` · actor [claude-code-d3420702](../actor/claude-code-d3420702.md)

## result

PASS

## statement

Built rhiza at the branch head (7420691c fix included) into a scratch root with fake rex and gh. A command step 'sleep 300' with stall = "5s", started as a session leader (setsid rhiza run start --foreground --autonomous): after 9s rhiza daemon --once logged 'step one silent since ...: stopped and resumed with --retry (budget 1 of 2)', and ps showed both the run and its 'sh -c exec 2>&1 sleep 300' child and sleep gone (cycle 1 left them orphaned). The next poll resumed it as a detached 'rhiza run resume --foreground <run> --retry'; that stalled again and the following poll stopped it, its sleep child gone too (budget 2 of 2). A wait step with timeout = "20s" under rhiza.daemon.watchdog.waitLimit = 5s was left alone past 9s (limit = timeout + 10m) and ended on its own timeout; the same wait step with stall = "5s" was stopped and retried.

## Links

### Out

- verifies [watchdog-stalled-step](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-stalled-step.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
