# daemon-service-install

capability · `01M6FES1Z3GF9HJGVS4Y1S687B` · parent [daemon](../daemon.md)

## acceptance

A testscript with fake systemctl and loginctl on PATH: install writes $HOME/.config/systemd/user/rhiza-<repo>.service with ExecStart running rhiza daemon, the main checkout as WorkingDirectory, Restart=always and After=rex.service, then calls systemctl --user daemon-reload, systemctl --user enable --now and loginctl enable-linger; a second install changes nothing; uninstall calls systemctl --user disable --now and removes the file; with no systemctl on PATH both exit 2 naming systemd

## statement

rhiza daemon install writes and enables a systemd user unit that runs this repo's daemon (in the main checkout, Restart=always, after rex.service) and turns on lingering, and rhiza daemon uninstall stops, disables and removes it; where systemd is missing both say so and change nothing

## Links

### Out

- decided-by [m24-server-daemon-no-tab](../../../../decision/m24-server-daemon-no-tab.md)

### In

- verifies [daemon-service-install-verified](../../../../evidence/daemon-service-install-verified.md)
- verifies [daemon-service-install-qa-1](../../../../evidence/daemon-service-install-qa-1.md)
- verifies [daemon-service-install-verified-2](../../../../evidence/daemon-service-install-verified-2.md)
- verifies [pr-200-the-daemon-runs-as-a-systemd](../../../../evidence/pr-200-the-daemon-runs-as-a-systemd.md)

## History

- `01M6FEMDQ1D2HWXT3GR4F3Z5H6` The daemon runs as a systemd service or in a Rex tab, as each installation chooses, and opens its runs in a fixed per-repo Rex session
