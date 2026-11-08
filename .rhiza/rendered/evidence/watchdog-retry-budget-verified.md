# watchdog-retry-budget-verified

evidence · `01M6WXC14YCXXHBSBCZWHYPWWD` · actor [claude-code-42e85804](../actor/claude-code-42e85804.md)

## result

PASS

## statement

At 071e243b ran `go test -trimpath ./internal/api -run TestWatchdogBudget -v -count=1`: PASS. The test polls a run that keeps failing on a transient gh 502, round-tripping the watchdog memory and holds through daemon.State JSON between polls (a daemon restart). It shows each attempt journalled in WatchdogRun.Attempts with its time and action, Spent counted across restarts, exactly two resumes with budget 2 and none after over three more polls, and one inbox item (watchdogHoldItems) that names the spent budget, the 502 and both attempt times, with the same ID on the next poll. Read the diff: every pass that charged Spent++ (stalled, transient, CI rerun, no-CI, requests, green-unmerged) now goes through watchdogDeps.charge, which also appends the attempt. The budget sits in daemon state per run, separate from step retries. watchdogBudgetHolds swaps any other hold on a spent run for one item listing every attempt. Also ran `go test -trimpath ./internal/api -run Watchdog` and `go test -trimpath ./internal/daemon/...`, both ok; go vet and gofmt are clean.

## Links

### Out

- verifies [watchdog-retry-budget](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-retry-budget.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
