# daemon-headless-run-drivers-qa-2

evidence · `01M6GNAEQC5SM5FDC76ZK697TD` · actor [claude-code-b014619b](../actor/claude-code-b014619b.md)

## result

PASS

## statement

At fdbebe1e I built rhiza with go build into a scratch dir and ran it by hand in a scratch git repo, with RHIZA_* and REX_* unset and a fake rex that logs its calls (no agent ever answers). rhiza daemon --permission-mode manual --once printed 'Started a run of <cs>' and exited. ps showed the driver 'rhiza run start --foreground <cs> --permission-mode manual' with PPID 1 in its own session. rex.log had only 'session ls' and the builder's own window new, no window new for the driver. The log was .git/rhiza/logs/<run>.log (no start-N file left) and held the child's 'step build/build: entered' line, and run ls --json showed the run running and live after the daemon exited. Under a long-lived rhiza daemon I killed a driver it had started: the pid was gone at once (reaped, no zombie), the run turned interrupted, and the daemon resumed it once as a detached 'rhiza run resume --foreground <run>' child of its own. After a second kill it was reaped too and the inbox listed the run as gave-up, with one daemon-resumed entry. rhiza request 'Add a widget' made the daemon start 'rhiza run start --foreground --workflow plan --new Add a widget --request <request>' detached, logging to <plan run>.log.

## Links

### Out

- verifies [daemon-headless-run-drivers](../capability/rhiza/orchestration/daemon/daemon-headless-run-drivers.md)

## History

- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
