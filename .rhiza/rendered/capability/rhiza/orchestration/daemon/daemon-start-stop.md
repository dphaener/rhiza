# daemon-start-stop

capability · `01M6FES4J3E8VXG2KW0Y9Z1C96` · parent [daemon](../daemon.md)

## acceptance

Testscripts with a fake systemctl and a fake Rex: with the unit installed, start calls systemctl --user start and stop calls systemctl --user stop, and daemon status --json reports mode service; with none, start opens a rex window named 'rhiza daemon' running rhiza daemon, stop ends that daemon, and status reports mode tab

## statement

rhiza daemon start and rhiza daemon stop start or stop this repo's daemon the way its installation runs it: through the systemd unit when rhiza daemon install wrote one, else in the Rex tab 'rhiza daemon' in the repo's session, and daemon status says which mode runs

## Links

### Out

- decided-by [m24-service-or-tab](../../../../decision/m24-service-or-tab.md)

### In

- verifies [daemon-start-stop-verified](../../../../evidence/daemon-start-stop-verified.md)
- verifies [daemon-start-stop-qa-1](../../../../evidence/daemon-start-stop-qa-1.md)
- verifies [pr-200-the-daemon-runs-as-a-systemd](../../../../evidence/pr-200-the-daemon-runs-as-a-systemd.md)
- verifies [pr-295-autonomy-app-on-modal-shell](../../../../evidence/pr-295-autonomy-app-on-modal-shell.md)

## History

- `01M6FEMDQ1D2HWXT3GR4F3Z5H6` The daemon runs as a systemd service or in a Rex tab, as each installation chooses, and opens its runs in a fixed per-repo Rex session
- `01M6MYPA335J1S3952RJAGZJN9` Autonomy app: the plan, accept, ship and sign-off toggles (rhiza.autonomy.*), rebuilt on the modal shell over daemon.settings and daemon.set, replacing internal/tui/autonomyapp and retiring the capabilities of the old autonomy app it replaces
