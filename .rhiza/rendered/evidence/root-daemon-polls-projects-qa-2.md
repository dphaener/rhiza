# root-daemon-polls-projects-qa-2

evidence · `01M6SKMFNHGWYF1Q062F2JFRDH` · actor [claude-code-32ca8314](../actor/claude-code-32ca8314.md)

## result

PASS

## statement

QA cycle 2 at b8630e96 with a scratch build (go build -trimpath ./cmd/rhiza) on scratch roots (rhiza root init) holding projects alpha and folder registered with project add --name other (each with three accepted change sets serving a milestone, one leaf each) and a bare uninitialised checkout, scratch HOME, fake rex/systemctl/make on PATH: one daemon --once from the root started runs in both alpha and folder (run logs in each project's own .git/rhiza/logs) and daemon status listed 'Skipped project: bare: not initialised'.

## Links

### Out

- verifies [root-daemon-polls-projects](../capability/rhiza/orchestration/daemon/root-daemon-polls-projects.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
