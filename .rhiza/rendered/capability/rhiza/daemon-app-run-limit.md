# daemon-app-run-limit

capability · `01M6MZC8G897NPM4BYCBAFE5MP` · parent [rhiza](../rhiza.md)

## acceptance

Go tests in the daemon app's package with a fake runner: the field starts at maxRuns from daemon.settings, a positive integer is written with daemon.set rhiza.daemon.maxRuns and the app quits, a non-number or 0 is refused by the field's validator and calls nothing.

## statement

The same Form shows the run limit as a Text field prefilled from rhiza.daemon.maxRuns; submitting a valid new limit writes it through daemon.set and closes the app.

## Links

### Out

- decided-by [tui-five-apps](../../decision/tui-five-apps.md)
- decided-by [tui-rexuikit-components](../../decision/tui-rexuikit-components.md)

### In

- verifies [daemon-app-run-limit-verified](../../evidence/daemon-app-run-limit-verified.md)
- verifies [daemon-app-run-limit-qa-1](../../evidence/daemon-app-run-limit-qa-1.md)
- verifies [pr-293-daemon-app-on-the-modal-shell](../../evidence/pr-293-daemon-app-on-the-modal-shell.md)

## History

- `01M6MYNQTAMMV8YEZHDS1G5HMD` Daemon app: pause or resume the daemon and set its run limit, rebuilt on the modal shell over daemon.settings and daemon.set (rhiza.daemon.paused, rhiza.daemon.maxRuns)
