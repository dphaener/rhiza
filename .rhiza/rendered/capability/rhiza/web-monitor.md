# web-monitor

capability · `01M6PFEKF6584S58PGZP17YAF2` · parent [rhiza](../rhiza.md)

## statement

The daemon serves a web UI on localhost for monitoring Rhiza: the daemon, runs, the inbox and the roadmap. It is read-only except for the safe actions behind the web token (answer an inbox item, approve a sign-off, pause or resume).

## Links

### Out

- decided-by [web-ui-in-daemon-read-only](../../decision/web-ui-in-daemon-read-only.md)

### In

- verifies [pr-574-safe-actions-in-the-web-ui](../../evidence/pr-574-safe-actions-in-the-web-ui.md)

## History

- `01M6PF31YCJ041HB94MDWQR331` I want to build out the first pass of a web UI that can be used to monitor the progress of Rhiza. This will run as part of the daemon process and be served on localhost. The user will be responsible for any port forwarding (if it's running on this VM for instance). First pass will be read only. For monitoring purposes only and we will build from there. The first requirement will be to obtain a design. You will generate a prompt that I will then feed to Claude design to generate a design system.
- `01M79NNP2BRD9R8XG15XCYACZ4` Safe actions in the web UI behind a required web token: answer an inbox item, approve a sign-off, pause or resume the daemon or a project; nothing else
