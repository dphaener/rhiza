# run-log-qa-2

evidence · `01M6GNCMK6WAGHQ2MZB72YV70P` · actor [claude-code-b014619b](../actor/claude-code-b014619b.md)

## result

PASS

## statement

At fdbebe1e I built rhiza with go build into a scratch dir and ran it by hand in a scratch git repo, with RHIZA_* and REX_* unset and a fake rex that logs its calls (no agent ever answers). Two daemon starts of the same change set (the first abandoned, the second queued) logged to <run1>.log and <run2>.log. rhiza run log <run1> --all showed only builder session 475ed046 and <run2> only 50641ced, so neither log showed the other's lines (this was the cycle-1 finding). A daemon resume and a queued continuation appended to <run2>.log. A request's plan run logged to <plan run>.log, and rhiza run log <plan run> --all printed its 'step interview: entered' lines. With the log file moved away, rhiza run log <run> said it has no log and exited 2.

## Links

### Out

- verifies [run-log](../capability/rhiza/orchestration/daemon/run-log.md)

## History

- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
