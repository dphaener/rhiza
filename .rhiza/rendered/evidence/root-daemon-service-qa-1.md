# root-daemon-service-qa-1

evidence · `01M6SFB3Q04Z0T6GGHJ9PMY43Z` · actor [claude-code-78f55b1d](../actor/claude-code-78f55b1d.md)

## result

FAIL

## statement

QA at 912172a0 with a scratch build (go build ./cmd/rhiza) on a scratch root (rhiza root init) holding projects alpha, folder registered with project add --name other (each with three accepted milestone change sets with a leaf) and a bare uninitialised checkout, scratch HOME, fake rex/systemctl/loginctl on PATH, web port moved to 17787: daemon install from alpha/sub writes one rhiza-root.service with WorkingDirectory=<root> and ExecStart ... daemon --autonomous and runs systemctl --user enable --now rhiza-root.service; daemon stop, start and uninstall at the root call systemctl --user stop/start/disable --now rhiza-root.service and remove the unit. But rhiza daemon status at the root (service mode and tab mode alike) fails with 'sessions need a git repository: git rev-parse --git-common-dir: fatal: not a git repository'; it works only inside a project.

## Links

### Out

- verifies [root-daemon-service](../capability/rhiza/orchestration/daemon/root-daemon-service.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
