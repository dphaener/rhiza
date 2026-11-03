# run-start-through-daemon-qa-1

evidence · `01M6G7R991J8QR1NMVWBV4NKQX` · actor [claude-code-5fcf5cce](../actor/claude-code-5fcf5cce.md)

## result

PASS

## statement

At a5ee2b13 I built rhiza (go build) and ran it by hand in a scratch git repo with RHIZA_* and REX_* unset, a fake rex that logs its calls and runs the agent command, and a fake claude agent (held by a file). rhiza run start <cs> printed 'Queued a ship-change run of <cs> for the daemon: no daemon runs, so it waits for one (rhiza daemon)', wrote .git/rhiza/starts.json and created no run. A second start of the same change set exited 2 with 'already has a run queued'. The queued change set served no milestone, and rhiza daemon --once still started it headless. rhiza run start --foreground --workflow review drove the run in the calling process to its sign-off, and nothing was queued. go test ./internal/tui/roadmapapp passes (fake runner). I did not drive the TUI itself. Minor: the queued entry's 'at' is the zero time.

## Links

### Out

- verifies [run-start-through-daemon](../capability/rhiza/orchestration/daemon/run-start-through-daemon.md)

## History

- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
