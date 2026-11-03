# run-log-qa-1

evidence · `01M6G7SC6R4D0D4XVD28MDFBBJ` · actor [claude-code-5fcf5cce](../actor/claude-code-5fcf5cce.md)

## result

FAIL

## statement

At a5ee2b13 I built rhiza (go build) and ran it by hand in a scratch git repo with RHIZA_* and REX_* unset, a fake rex that logs its calls and runs the agent command, and a fake claude agent (held by a file). rhiza run log <run> printed the daemon-started run's driver log. --limit 2 --json returned 2 lines and next_cursor 2. --cursor 2 returned the next page and a 'next: ... --cursor 4' hint. --all printed everything. A foreground-only run exited 2 with 'has no log: only a run the daemon started has one'. --layer with REX_SESSION sent rex do with a new_layer running 'rhiza run log <run> --all | less +G', and without a Rex session it failed with exit 1. go test ./internal/tui/liveapp passes. FAIL: a second daemon start of the same change set appends to the same logs/<change set>.log. rhiza run log <new run> then opened with the earlier run's lines (its builder session 9f66a989), and the earlier run's log took on the new run's lines, so run log does not print that run's own log.

## Links

### Out

- verifies [run-log](../capability/rhiza/orchestration/daemon/run-log.md)

## History

- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
