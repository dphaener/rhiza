# autonomy-app-toggles

capability · `01M6MZCNKTRF4KEY30SV1J4YHD` · parent [rhiza](../rhiza.md)

## acceptance

Go tests in internal/tui/autonomyapp with a fake runner: the four fields start at the values daemon.settings gives, submitting calls daemon.set once per changed toggle and none for unchanged ones, then quits; Esc calls nothing; the tui-autonomy-app testscript's --once still prints what rhiza daemon settings --json prints.

## statement

The Autonomy app is one rexuikit Form of the plan, accept, ship and sign-off toggles, each auto or owner, prefilled from daemon.settings; submitting writes each changed toggle (rhiza.autonomy.plan, accept, ship, signoff) through daemon.set and closes the app.

## Links

### Out

- decided-by [tui-five-apps](../../decision/tui-five-apps.md)
- decided-by [tui-rexuikit-components](../../decision/tui-rexuikit-components.md)

### In

- verifies [autonomy-app-toggles-verified](../../evidence/autonomy-app-toggles-verified.md)
- verifies [autonomy-app-toggles-qa-1](../../evidence/autonomy-app-toggles-qa-1.md)
- verifies [autonomy-app-toggles-verified-2](../../evidence/autonomy-app-toggles-verified-2.md)
- verifies [pr-295-autonomy-app-on-modal-shell](../../evidence/pr-295-autonomy-app-on-modal-shell.md)

## History

- `01M6MYPA335J1S3952RJAGZJN9` Autonomy app: the plan, accept, ship and sign-off toggles (rhiza.autonomy.*), rebuilt on the modal shell over daemon.settings and daemon.set, replacing internal/tui/autonomyapp and retiring the capabilities of the old autonomy app it replaces
