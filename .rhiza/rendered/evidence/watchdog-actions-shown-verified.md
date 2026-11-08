# watchdog-actions-shown-verified

evidence · `01M6WZ0K2NZ7PJTW0A9ZD1TK0B` · actor [claude-code-453dfa0b](../actor/claude-code-453dfa0b.md)

## result

PASS

## statement

At 21a63405: go test -trimpath ./cmd/rhiza -run 'TestScripts/(watchdog-health|describe-golden)$' -count=1 passed: rhiza health prints "watchdog's last actions, newest first:" with run r2 (00:02) above run r1 (00:01) and the timestamped transient-recovery line, and daemon status --json carries watchdog_actions. The journal there is a state fixture whose lines match the watchdog's real transient-recovery log format; go test -trimpath ./internal/api -run TestWatchdogLog -count=1 -v passed, faking a 502 transient failure through watchdogPasses over two polls: one log line per action, journalled newest first without the watchdog: prefix, with what it saw (502), what it did (resumed with --retry) and the run, a repeat replacing its entry, bounded at 20, and listed in health detail. Read the diff: daemonWatchdog wraps d.log with watchdogJournalLog, poll-behind lines are journalled too, the journal survives daemon restart (daemonStartState), and go vet is clean. The wider -run 'Watchdog|Health|DaemonStatus' run showed no test failures; it reported FAIL only because tests wrote to the enclosing ~/code/.rhiza-root, a known environment leak unrelated to this leaf.

## Links

### Out

- verifies [watchdog-actions-shown](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-actions-shown.md)

## History

- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
