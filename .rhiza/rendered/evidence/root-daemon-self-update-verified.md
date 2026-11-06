# root-daemon-self-update-verified

evidence · `01M6SEM241QYW1GCHR2PE9P1Z1` · actor [claude-code-1b977378](../actor/claude-code-1b977378.md)

## result

PASS

## statement

Verifier claude-code-1b977378 at 6c4b2376: go test -trimpath -count=1 ./internal/daemon/... ./internal/rootcfg/... pass; TestSelfUpdateProjectChoosesCheckout passes: with self_update_project=rhiza (path code/rhiza) every fake git call and the final fake make run in /r/code/rhiza, and with it unset ForCheckout reports false and neither git nor make runs. internal/api TestDaemonSelfUpdateSetting passes: selfUpdate=true with no project named runs nothing, with the project named git and make run in that project's folder; -run 'Daemon|SelfUpdate|Root' ./internal/api ok. gofmt and go vet clean. Diff read: runDaemon now resolves the checkout via daemon.SelfUpdateCheckout(root, rootcfg.Load(root)) instead of the home project's first worktree. The full internal/api package timed out at 9m under load average 10 on 4 cores, inside TestLinkEndpointKinds, which this change does not touch.

## Links

### Out

- verifies [root-daemon-self-update](../capability/rhiza/orchestration/daemon/root-daemon-self-update.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
