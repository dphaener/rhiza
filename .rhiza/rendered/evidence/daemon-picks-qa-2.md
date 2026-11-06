# daemon-picks-qa-2

evidence · `01M6SKMF5NBWBQJ113WGY3HF22` · actor [claude-code-32ca8314](../actor/claude-code-32ca8314.md)

## result

PASS

## statement

QA cycle 2 at b8630e96 with a scratch build (go build -trimpath ./cmd/rhiza) on scratch roots (rhiza root init) holding projects alpha and folder registered with project add --name other (each with three accepted change sets serving a milestone, one leaf each) and a bare uninitialised checkout, scratch HOME, fake rex/systemctl/make on PATH: with rhiza.daemon.maxRuns = 3 in rhiza.toml one daemon --once started exactly three runs (cycle 1 started six); a second --once started none, skipping alpha cs1/cs2 'has a live run' and alpha cs3 'waiting for a slot'; bare listed 'not initialised'.

## Links

### Out

- verifies [daemon-picks](../capability/rhiza/orchestration/daemon/daemon-picks.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
