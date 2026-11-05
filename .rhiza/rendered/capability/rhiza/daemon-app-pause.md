# daemon-app-pause

capability · `01M6MZC0J2KTTEAQ16DF31DZBV` · parent [rhiza](../rhiza.md)

## acceptance

Go tests in the daemon app's package with a fake runner: the Confirm starts at the paused value daemon.settings gives, flipping it and submitting calls daemon.set with rhiza.daemon.paused and the new value and quits, submitting unchanged or Esc calls no daemon.set.

## statement

The Daemon app's rexuikit Form shows whether the daemon is paused as a Confirm, prefilled from daemon.settings; submitting a change writes rhiza.daemon.paused through daemon.set and closes the app.

## Links

### Out

- decided-by [tui-five-apps](../../decision/tui-five-apps.md)
- decided-by [tui-rexuikit-components](../../decision/tui-rexuikit-components.md)

### In

- verifies [daemon-app-pause-verified](../../evidence/daemon-app-pause-verified.md)
- verifies [daemon-app-pause-qa-1](../../evidence/daemon-app-pause-qa-1.md)
- verifies [pr-293-daemon-app-on-the-modal-shell](../../evidence/pr-293-daemon-app-on-the-modal-shell.md)

## History

- `01M6MYNQTAMMV8YEZHDS1G5HMD` Daemon app: pause or resume the daemon and set its run limit, rebuilt on the modal shell over daemon.settings and daemon.set (rhiza.daemon.paused, rhiza.daemon.maxRuns)
