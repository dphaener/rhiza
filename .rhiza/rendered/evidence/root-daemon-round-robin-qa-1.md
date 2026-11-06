# root-daemon-round-robin-qa-1

evidence · `01M6SFB3S9PTJT9BH0STQ81FPN` · actor [claude-code-78f55b1d](../actor/claude-code-78f55b1d.md)

## result

FAIL

## statement

QA at 912172a0 with a scratch build (go build ./cmd/rhiza) on a scratch root (rhiza root init) holding projects alpha, folder registered with project add --name other (each with three accepted milestone change sets with a leaf) and a bare uninitialised checkout, scratch HOME, fake rex/systemctl/loginctl on PATH, web port moved to 17787: with rhiza.daemon.maxRuns = 3 in rhiza.toml (daemon status shows Max runs: 3), one daemon --once started six runs: all three of alpha's, then all three of other's. The machine-wide cap is applied per project and the picks are not round-robin; daemon.RoundRobin has no caller outside tests.

## Links

### Out

- verifies [root-daemon-round-robin](../capability/rhiza/orchestration/daemon/root-daemon-round-robin.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
