# daemon-service-install-verified

evidence · `01M6FMG9QFZZYFX3YXAJQETATV` · actor [claude-code-f90d97d7](../actor/claude-code-f90d97d7.md)

## result

PASS

## statement

Verified at 0afdeea5. cmd/rhiza/testdata/script/daemon-install.txtar puts fake systemctl and loginctl on PATH. It checks that install writes $HOME/.config/systemd/user/rhiza-repo.service with ExecStart=<rhiza> daemon, WorkingDirectory set to the main checkout (also when run from a worktree), Restart=always and After=rex.service. It checks the calls in order: systemctl --user daemon-reload, systemctl --user enable --now rhiza-repo.service, loginctl enable-linger. A second install leaves the file byte-identical and does no daemon-reload. Uninstall calls systemctl --user disable --now and daemon-reload and removes the file. With PATH holding only the test's rhiza dir (RHIZA_TEST_BIN), both commands exit 2 with 'needs systemd: systemctl is not on PATH' and write nothing. Read internal/api/daemon_service.go: it matches the test, and a unit not written by install is never overwritten or removed. Attempt-1 finding f1 (the test inherited the host PATH) is fixed. With logging fake systemctl/loginctl prepended to the host PATH, go test ./cmd/rhiza -run TestScripts/daemon-install -count=1 passed and the host fakes were never run. In a clean clone at the commit, go test ./... passed (an earlier make test run failed only on the read-timing budget test, which does not touch this leaf), and make fmt-check vet self-check passed.

## Links

### Out

- verifies [daemon-service-install](../capability/rhiza/orchestration/daemon/daemon-service-install.md)

## History

- `01M6FEMDQ1D2HWXT3GR4F3Z5H6` The daemon runs as a systemd service or in a Rex tab, as each installation chooses, and opens its runs in a fixed per-repo Rex session
