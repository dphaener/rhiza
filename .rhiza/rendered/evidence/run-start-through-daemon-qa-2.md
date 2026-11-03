# run-start-through-daemon-qa-2

evidence · `01M6GNBFVMPYCANWG5ARYQ8T8C` · actor [claude-code-b014619b](../actor/claude-code-b014619b.md)

## result

PASS

## statement

At fdbebe1e I built rhiza with go build into a scratch dir and ran it by hand in a scratch git repo, with RHIZA_* and REX_* unset and a fake rex that logs its calls (no agent ever answers). After the first run of the change set was abandoned, rhiza run start <cs> --permission-mode manual printed 'Queued a ship-change run of <cs> for the daemon: no daemon runs, so it waits for one (rhiza daemon)' and wrote starts.json. A long-lived rhiza daemon then started it headless as a detached 'rhiza run start --foreground <cs>' child, drained the queue (starts: null) and opened no Rex tab for the driver.

## Links

### Out

- verifies [run-start-through-daemon](../capability/rhiza/orchestration/daemon/run-start-through-daemon.md)

## History

- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
