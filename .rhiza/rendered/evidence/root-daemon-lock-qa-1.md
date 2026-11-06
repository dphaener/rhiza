# root-daemon-lock-qa-1

evidence · `01M6SFB3FYJ5RY1A2FRP629V2D` · actor [claude-code-78f55b1d](../actor/claude-code-78f55b1d.md)

## result

PASS

## statement

QA at 912172a0 with a scratch build (go build ./cmd/rhiza) on a scratch root (rhiza root init) holding projects alpha, folder registered with project add --name other (each with three accepted milestone change sets with a leaf) and a bare uninitialised checkout, scratch HOME, fake rex/systemctl/loginctl on PATH, web port moved to 17787: daemon --once in an empty dir outside the root exits 2 with 'not inside a Rhiza root' and next: rhiza root init. With one daemon running in the background from the root, a second daemon --once from alpha/sub and from the root each exit 2 with 'another daemon holds the lock: process rhiza:<pid>'; the root's .rhiza-root/ holds daemon.lock.

## Links

### Out

- verifies [root-daemon-lock](../capability/rhiza/orchestration/daemon/root-daemon-lock.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
