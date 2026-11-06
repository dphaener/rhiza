# daemon-service-install-verified-2

evidence · `01M6FPMG00EF77TV1R04P9HK73` · actor [claude-code-b8008485](../actor/claude-code-b8008485.md)

## result

PASS

## statement

At f01bef3c (attempt 2), go test ./cmd/rhiza -run TestScripts/daemon -count=1 passed. daemon-install.txtar puts fake systemctl and loginctl on PATH. It checks that install, run from a worktree, writes $HOME/.config/systemd/user/rhiza-repo.service with ExecStart=<bin> daemon --autonomous, WorkingDirectory set to the main checkout, Restart=always and After=rex.service. It checks the calls are exactly systemctl --user daemon-reload, systemctl --user enable --now rhiza-repo.service and loginctl enable-linger. A second install prints 'nothing changed', leaves the unit byte-identical and does not reload. Uninstall calls systemctl --user disable --now and daemon-reload, and the file is gone. With only rhiza on PATH, both commands fail with 'needs systemd: systemctl is not on PATH' and write nothing. I also built the binary and ran it by hand with PATH=/usr/bin:/bin: install and uninstall each exited 2 with that message. Diff read: the fix shares daemonArgs={daemon,--autonomous} between the unit's ExecStart and the Rex tab, which resolves the quality finding f1. go vet and gofmt are clean.

## Links

### Out

- verifies daemon-service-install (retired)

## History

- `01M6FEMDQ1D2HWXT3GR4F3Z5H6` The daemon runs as a systemd service or in a Rex tab, as each installation chooses, and opens its runs in a fixed per-repo Rex session
