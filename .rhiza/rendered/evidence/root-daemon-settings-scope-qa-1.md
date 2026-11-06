# root-daemon-settings-scope-qa-1

evidence · `01M6SFB3NBS0XEF9CSGVKDA6ED` · actor [claude-code-78f55b1d](../actor/claude-code-78f55b1d.md)

## result

PASS

## statement

QA at 912172a0 with a scratch build (go build ./cmd/rhiza) on a scratch root (rhiza root init) holding projects alpha, folder registered with project add --name other (each with three accepted milestone change sets with a leaf) and a bare uninitialised checkout, scratch HOME, fake rex/systemctl/loginctl on PATH, web port moved to 17787: daemon set rhiza.daemon.maxRuns 3 (also from inside alpha) lands in rhiza.toml [daemon]; daemon set rhiza.autonomy.ship owner inside alpha and --project other rhiza.autonomy.plan auto at the root land in that repo's git config only; a project key at the root without --project exits 2 naming --project; a bad maxRuns exits 2; daemon settings and settings --project other --json list each key with scope machine/project and its source, trimGoCache machine.

## Links

### Out

- verifies [root-daemon-settings-scope](../capability/rhiza/orchestration/daemon/root-daemon-settings-scope.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
