# run-continue-through-daemon-qa-2

evidence · `01M6GNC5JTSVGRH5PPR7TZNB85` · actor [claude-code-b014619b](../actor/claude-code-b014619b.md)

## result

PASS

## statement

At fdbebe1e I built rhiza with go build into a scratch dir and ran it by hand in a scratch git repo, with RHIZA_* and REX_* unset and a fake rex that logs its calls (no agent ever answers). On a run the daemon had given up on (interrupted again after its one resume), rhiza run resume <run> with a live daemon printed 'Queued the continuation of run <run> for the daemon: daemon rhiza:<pid> continues it headless'. The daemon then ran a detached 'rhiza run resume --foreground <run>' child in its own session, with no rex window new for it. run ls showed the run running, and its log gained another 'continuing visit 1' line.

## Links

### Out

- verifies [run-continue-through-daemon](../capability/rhiza/orchestration/daemon/run-continue-through-daemon.md)

## History

- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
