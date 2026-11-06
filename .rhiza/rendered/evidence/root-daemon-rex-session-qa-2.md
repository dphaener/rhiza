# root-daemon-rex-session-qa-2

evidence · `01M6SKMFQHMV4VH84YSR1HRSZN` · actor [claude-code-32ca8314](../actor/claude-code-32ca8314.md)

## result

PASS

## statement

QA cycle 2 at b8630e96 with a scratch build (go build -trimpath ./cmd/rhiza) on scratch roots (rhiza root init) holding projects alpha and folder registered with project add --name other (each with three accepted change sets serving a milestone, one leaf each) and a bare uninitialised checkout, scratch HOME, fake rex/systemctl/make on PATH: with REX_SESSION unset and rex session ls empty, one poll ran rex session new rhiza-alpha (-c alpha) and rhiza-other (-c folder) once each and opened alpha's builder windows with -s rhiza-alpha and folder's with -s rhiza-other; rex.log never names rhiza-folder.

## Links

### Out

- verifies [root-daemon-rex-session](../capability/rhiza/orchestration/daemon/root-daemon-rex-session.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
