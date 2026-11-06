# root-daemon-self-update-qa-2

evidence · `01M6SKMFHFHY6G4QJMYEGTQVQ1` · actor [claude-code-32ca8314](../actor/claude-code-32ca8314.md)

## result

PASS

## statement

QA cycle 2 at b8630e96 with a scratch build (go build -trimpath ./cmd/rhiza) on scratch roots (rhiza root init) holding projects alpha and folder registered with project add --name other (each with three accepted change sets serving a milestone, one leaf each) and a bare uninitialised checkout, scratch HOME, fake rex/systemctl/make on PATH: with rhiza.daemon.selfUpdate true and folder's origin one commit ahead, a foreground daemon with no self_update_project ran no make and left folder at init; with self_update_project = "other" it logged 'self-update: updated to 8b71fbf78', fast-forwarded folder to the newer commit, ran make update with cwd folder, and left alpha untouched.

## Links

### Out

- verifies [root-daemon-self-update](../capability/rhiza/orchestration/daemon/root-daemon-self-update.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
