# root-daemon-lock-verified

evidence · `01M6SBY83C3XNW2ZQ8CNAH43S7` · actor [claude-code-a098c954](../actor/claude-code-a098c954.md)

## result

PASS

## statement

At dac8ca29 (attempt 2) I read the diffs of 0497b08d and dac8ca29. runDaemon resolves the root with rootcfg.Find(ctx.Wd), exits as a usage error with next 'rhiza root init' outside a root, and takes daemon.TryLockRoot(root) in <root>/.rhiza-root before the per-repo lock. go test -trimpath -count=1 ./internal/daemon passes, including TestRootLockIsOnePerRoot and TestRootLockFromProjectAndRootAlike. The second test locks from the root and then from <root>/app/sub, resolved through rootcfg.Find, and also in the reverse order. Each time the second lock gets a *Held naming this process, and there is one lock file in <root>/.rhiza-root. Both attempt-1 findings are addressed. The testscripts daemon-no-root, daemon-settings, daemon-pause, daemon-event-hook, daemon-repo-session and stage-autonomy pass. daemon-no-root asserts exit=2, 'not inside a Rhiza root' and '^next: rhiza root init$' in $NOROOT. git diff 0497b08d~1..HEAD shows the five existing daemon scripts unedited. go test ./internal/api -run Daemon also passes.

## Links

### Out

- verifies [root-daemon-lock](../capability/rhiza/orchestration/daemon/root-daemon-lock.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
