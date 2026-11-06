# root-daemon-service-qa-2

evidence · `01M6SKMFEYYVGYDB1QYHCKMWRM` · actor [claude-code-32ca8314](../actor/claude-code-32ca8314.md)

## result

PASS

## statement

QA cycle 2 at b8630e96 with a scratch build (go build -trimpath ./cmd/rhiza) on scratch roots (rhiza root init) holding projects alpha and folder registered with project add --name other (each with three accepted change sets serving a milestone, one leaf each) and a bare uninitialised checkout, scratch HOME, fake rex/systemctl/make on PATH: daemon install from alpha/sub wrote one rhiza-root.service with WorkingDirectory=<root> and ExecStart .../rhiza daemon --autonomous, running systemctl --user daemon-reload and enable --now rhiza-root.service; daemon status at the root and from folder both printed 'Daemon not running / Mode: service' and exited 0 (cycle 1 failed with 'sessions need a git repository'); stop and start at the root, stop from folder and uninstall called systemctl --user stop/start/stop/disable --now rhiza-root.service plus daemon-reload and removed the unit. With a tab-mode daemon running from the root, status at the root and from folder both said 'Daemon running, process rhiza:<pid>'.

## Links

### Out

- verifies [root-daemon-service](../capability/rhiza/orchestration/daemon/root-daemon-service.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
