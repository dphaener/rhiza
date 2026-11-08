# watchdog-stuck-requests-verified

evidence · `01M6WWQRF22MPND8E3YSSFMGTK` · actor [claude-code-db5f5583](../actor/claude-code-db5f5583.md)

## result

PASS

## statement

At e89dc1ca, ran go test -trimpath -count=1 ./internal/api -run 'TestWatchdogRequests$' -v: PASS. The test drives watchdogRequests over fake journals and a fake resume: request r1 whose plan run is unended and not live (gone, still planning) is resumed; r3 whose plan run failed on a git index.lock transient signature is retried (resume --retry, logged with the signature); r4 failed with 'exit 1' is not resumed (left for the owner); the live run, a milestone request and an unstarted request are untouched; a second poll with RetryBudget 1 resumes nothing. Read the diff: the pass is wired into watchdogPasses via daemon.ReadRequests, shares the per-run budget memory and the continue queue check, skips runs the failed-transient pass already owns, and leaves planned (succeeded) requests alone. Also ran the whole -run TestWatchdog set: every test passes; the package then reports the known pre-existing usage.live.json leak into ~/code/.rhiza-root, which TestWatchdogRequests alone does not trigger.

## Links

### Out

- verifies [watchdog-stuck-requests](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-stuck-requests.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
