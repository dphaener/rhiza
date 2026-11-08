# watchdog-transient-failure-verified

evidence · `01M6WWZSBZMF8MMJXH22TQG8QH` · actor [claude-code-c055f7a6](../actor/claude-code-c055f7a6.md)

## result

PASS

## statement

At 8cb18c39, go test -trimpath -count=1 ./internal/api -run TestWatchdog -v passed: TestWatchdogTransient shows the go $WORK, git index.lock and GitHub 502 failures each resumed once with --retry, a non-matching failure (FAIL: TestThing) and an attended run left to the owner, and a live run with a 503 output left alone; a second poll inside the budget of 1 resumes nothing. The new TestWatchdogTransientBudgetAcrossPolls runs whole polls through watchdogPasses (the stalled-step pass included, as daemonWatchdog now does) and shows a failed run on HTTP 502 resumed once over five polls with Spent 1, no resume or charge while its continue is already queued, and its memory forgotten once the run is gone. Read the diff: watchdogPasses carries prev memory over for runs that are not live, fixing attempt 1's finding f1 (budget dropped each poll), and the transient and CI passes skip a queued run before charging. go vet ./internal/api clean. A whole-package run hung only in TestReadBudget (roadmap read, plan.go link probe on the SQLite index, load ~20), code this change set does not touch; left to CI.

## Links

### Out

- verifies [watchdog-transient-failure](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-transient-failure.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
