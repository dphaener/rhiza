# root-daemon-rex-session-qa-1

evidence · `01M6SFB3KN8QMR4WJE216HSRA4` · actor [claude-code-78f55b1d](../actor/claude-code-78f55b1d.md)

## result

PASS

## statement

QA at 912172a0 with a scratch build (go build ./cmd/rhiza) on a scratch root (rhiza root init) holding projects alpha, folder registered with project add --name other (each with three accepted milestone change sets with a leaf) and a bare uninitialised checkout, scratch HOME, fake rex/systemctl/loginctl on PATH, web port moved to 17787: with REX_SESSION unset and rex session ls empty, daemon --once ran rex session new rhiza-alpha and rex session new rhiza-other (cwd each project) once each, opened alpha's builder windows with -s rhiza-alpha and other's with -s rhiza-other, and rex.log never names rhiza-folder.

## Links

### Out

- verifies [root-daemon-rex-session](../capability/rhiza/orchestration/daemon/root-daemon-rex-session.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
