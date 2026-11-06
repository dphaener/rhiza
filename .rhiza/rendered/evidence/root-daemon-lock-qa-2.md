# root-daemon-lock-qa-2

evidence · `01M6SKMFK9N9M7H3QVPRH8HZBR` · actor [claude-code-32ca8314](../actor/claude-code-32ca8314.md)

## result

PASS

## statement

QA cycle 2 at b8630e96 with a scratch build (go build -trimpath ./cmd/rhiza) on scratch roots (rhiza root init) holding projects alpha and folder registered with project add --name other (each with three accepted change sets serving a milestone, one leaf each) and a bare uninitialised checkout, scratch HOME, fake rex/systemctl/make on PATH: with a daemon running in the background from the root, daemon --permission-mode manual --once from folder exited 2 with 'another daemon holds the lock: process rhiza:<pid>' and next: rhiza daemon status; daemon --once in a dir outside any root exited 2 with 'not inside a Rhiza root' and next: rhiza root init.

## Links

### Out

- verifies [root-daemon-lock](../capability/rhiza/orchestration/daemon/root-daemon-lock.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
