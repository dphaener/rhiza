# root-daemon-polls-projects-qa-1

evidence · `01M6SFB3HJEH92KQP217WDN9JC` · actor [claude-code-78f55b1d](../actor/claude-code-78f55b1d.md)

## result

PASS

## statement

QA at 912172a0 with a scratch build (go build ./cmd/rhiza) on a scratch root (rhiza root init) holding projects alpha, folder registered with project add --name other (each with three accepted milestone change sets with a leaf) and a bare uninitialised checkout, scratch HOME, fake rex/systemctl/loginctl on PATH, web port moved to 17787: one daemon --once from the root started runs in both alpha and other (run logs in each project's own .git/rhiza/logs), daemon status in alpha lists 'Skipped project: bare: not initialised', and bare/.git/rhiza was never created.

## Links

### Out

- verifies [root-daemon-polls-projects](../capability/rhiza/orchestration/daemon/root-daemon-polls-projects.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
