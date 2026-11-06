# daemon-service-install-qa-1

evidence · `01M6FN774V4WDJET7SY66D7C6M` · actor [claude-code-e6cae1c1](../actor/claude-code-e6cae1c1.md)

## result

PASS

## statement

QA at 78a3b478 with a binary built from the branch on macOS, in a scratch repo qarepo with HOME pointed at a scratch dir. With no systemctl on PATH, daemon install and daemon uninstall each exited 2 with 'needs systemd: systemctl is not on PATH' and wrote nothing. With logging fake systemctl and loginctl on PATH, install wrote ~/.config/systemd/user/rhiza-qarepo.service (marker comment, After=rex.service, WorkingDirectory the main checkout, ExecStart=<binary> daemon, Restart=always, WantedBy=default.target) and called systemctl --user daemon-reload, systemctl --user enable --now rhiza-qarepo.service and loginctl enable-linger. A second install said nothing changed and left the file as it was, with no daemon-reload (it repeated enable --now and enable-linger, which are idempotent). From a linked worktree, daemon status said mode service. Uninstall called systemctl --user disable --now and daemon-reload, removed the file, and status went back to mode tab.

## Links

### Out

- verifies daemon-service-install (retired)

## History

- `01M6FEMDQ1D2HWXT3GR4F3Z5H6` The daemon runs as a systemd service or in a Rex tab, as each installation chooses, and opens its runs in a fixed per-repo Rex session
