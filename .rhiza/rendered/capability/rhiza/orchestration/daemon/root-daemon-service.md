# root-daemon-service

capability · `01M6RN85E14DDFJGJ96GKN7ZGE` · parent [daemon](../daemon.md)

## acceptance

A testscript with a fake systemctl on PATH runs daemon install from inside a project and finds one unit whose WorkingDirectory is the root, then start, stop and uninstall call systemctl on that unit.

## statement

rhiza daemon install writes and enables one systemd user unit for the root (rhiza-root.service, working directory the root), and rhiza daemon start, stop, status and uninstall act on the root's daemon from anywhere under the root.

## Links

### Out

- decided-by [m29-one-root-daemon](../../../../decision/m29-one-root-daemon.md)
- supersedes daemon-service-install (retired)

### In

- verifies [root-daemon-service-verified](../../../../evidence/root-daemon-service-verified.md)
- verifies [root-daemon-service-qa-1](../../../../evidence/root-daemon-service-qa-1.md)
- verifies [root-daemon-service-verified-2](../../../../evidence/root-daemon-service-verified-2.md)
- verifies [root-daemon-service-qa-2](../../../../evidence/root-daemon-service-qa-2.md)
- verifies [pr-338-root-daemon](../../../../evidence/pr-338-root-daemon.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
