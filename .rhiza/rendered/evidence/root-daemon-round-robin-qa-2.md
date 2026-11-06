# root-daemon-round-robin-qa-2

evidence · `01M6SKMFA3D4ZWFR96T2EYM1S7` · actor [claude-code-32ca8314](../actor/claude-code-32ca8314.md)

## result

PASS

## statement

QA cycle 2 at b8630e96 with a scratch build (go build -trimpath ./cmd/rhiza) on scratch roots (rhiza root init) holding projects alpha and folder registered with project add --name other (each with three accepted change sets serving a milestone, one leaf each) and a bare uninitialised checkout, scratch HOME, fake rex/systemctl/make on PATH: maxRuns 3 over alpha(3) and other(3) started alpha cs1, other cs1, alpha cs2 (two in alpha, one in other, alternating), and alpha cs3 waited for a slot; on a fresh root with maxRuns 4 and rhiza.project.maxRuns 1 set in alpha, one poll started four runs: alpha cs1 only, and other cs1, cs2, cs3.

## Links

### Out

- verifies [root-daemon-round-robin](../capability/rhiza/orchestration/daemon/root-daemon-round-robin.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
