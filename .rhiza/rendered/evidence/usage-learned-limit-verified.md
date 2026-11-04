# usage-learned-limit-verified

evidence · `01M6J2TSSXR960K9K6DH9EG6B1` · actor [claude-code-1f6cb4be](../actor/claude-code-1f6cb4be.md)

## result

PASS

## statement

At 2f765569 ran go test ./internal/api -run 'TestLimitHitsBecomeLearnedLimits|TestSessionRunRecordsLimitHit' -count=1 -v: both PASS. With a fake adapter (usageFake) and fake token records, a limit hit with no reset time records nothing and teaches nothing; a hit resetting 3h out is recorded and learnLimits (the function the daemon poll calls before saving state) puts the session window's use at the hit (350) into daemon.State.Learned as kind session; a hit resetting 4 days out becomes the weekly learned limit (1349); a set setting wins over a learned limit and a window with neither has none. TestSessionRunRecordsLimitHit shows session.run (the run driver) appends the hit to the git common dir's usage.limits.jsonl. Read the diff: learnLimits is called each poll in runDaemon and Learned survives a restart via daemonStartState; KindOf maps reset <=5h to session, later to weekly, zero reset to nothing. go test ./internal/api (full, re-run), ./internal/usage and ./internal/daemon pass; go vet and gofmt -l clean. One earlier full internal/api run failed on a timing budget test under load and passed on re-run.

## Links

### Out

- verifies [usage-learned-limit](../capability/rhiza/orchestration/usage-governor/usage-learned-limit.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
