# daemon-picks-root-qa-1

evidence · `01M6SFB3V942SCFZ80EVQBEGQ7` · actor [claude-code-78f55b1d](../actor/claude-code-78f55b1d.md)

## result

FAIL

## statement

QA at 912172a0 with a scratch build (go build ./cmd/rhiza) on a scratch root (rhiza root init) holding projects alpha, folder registered with project add --name other (each with three accepted milestone change sets with a leaf) and a bare uninitialised checkout, scratch HOME, fake rex/systemctl/loginctl on PATH, web port moved to 17787: with machine maxRuns 3, one poll started six runs (three per project): the daemon does not stop when its running runs reach the machine-wide maxRuns. The milestone filter held: change sets serving no milestone were skipped with 'serves no milestone'.

## Links

### Out

- verifies [daemon-picks](../capability/rhiza/orchestration/daemon/daemon-picks.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
