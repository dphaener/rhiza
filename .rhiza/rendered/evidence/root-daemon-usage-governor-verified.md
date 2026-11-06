# root-daemon-usage-governor-verified

evidence · `01M6SA0BQM8H9D1DXRNGMRT6GS` · actor [claude-code-1511520f](../actor/claude-code-1511520f.md)

## result

PASS

## statement

At f2157764: go test -trimpath -count=1 -run Governor ./internal/daemon/ passes TestGovernorHoldSpansProjectsAndLifts (a hit recorded with RecordLimitHitIn in RootStateDir holds via ActiveHold before reset, still holds at reset-1ns, lifts at the reset time). Read the diff: recordLimitHit, learnLimits and usageWindowsNow now use limitsDir (the root's .rhiza-root when rootcfg.Find(ctx.Wd) finds a rhiza.toml, else the repo's own folder), and runDaemon falls back to governorHold when no reserve hold applies, passing it to daemonPoll1, which stops launches. A throwaway api test in a scratch copy (not committed) recorded a hit through limitsDir from project a's checkout under a root and saw governorHold return the hit's reset for project b and for a, nil for both at the reset time, and nil for a repo outside the root; it passed together with the api package's Usage/Limit/Daemon tests. Testscripts session-run-usage-limit, usage, usage-hold, usage-settings, tui-launch-usage, daemon-status and daemon-picks pass.

## Links

### Out

- verifies [root-daemon-usage-governor](../capability/rhiza/orchestration/usage-governor/root-daemon-usage-governor.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
