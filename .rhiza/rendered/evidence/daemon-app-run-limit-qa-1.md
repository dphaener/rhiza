# daemon-app-run-limit-qa-1

evidence · `01M6NB2SXEQR6K2D9NKDDP0CZ6` · actor [claude-code-3fad0677](../actor/claude-code-3fad0677.md)

## result

PASS

## statement

At f06d3d29, same built binary and harness against a scratch git repo: the Run limit field opened at maxRuns=2 (daemon.settings' default); Enter, Backspace, 7, Enter ran 'daemon set -- rhiza.daemon.maxRuns 7', closed the app and git config read 7; the next open started at 7. Flipping pause and typing 12 wrote both settings (paused true, maxRuns 12). Typing 0, x or an empty value and Enter showed '✗ must be a whole number of at least 1', kept the app open, ran no daemon set and left git config at 12.

## Links

### Out

- verifies [daemon-app-run-limit](../capability/rhiza/daemon-app-run-limit.md)

## History

- `01M6MYNQTAMMV8YEZHDS1G5HMD` Daemon app: pause or resume the daemon and set its run limit, rebuilt on the modal shell over daemon.settings and daemon.set (rhiza.daemon.paused, rhiza.daemon.maxRuns)
