# usage-hold-reported-verified-5

evidence · `01M6JR04JV00FF9D39JBPCS37Y` · actor [claude-code-564b0ad4](../actor/claude-code-564b0ad4.md)

## result

PASS

## statement

Verifier at 28586e95 (attempt 5). Ran go test ./cmd/rhiza -run 'TestScripts/(usage-hold|describe-golden)$' -count=1: ok. usage-hold.txtar seeds a holding daemon state and checks: the usage-reserve skip and its detail (window, 900 used, 950 projected, limit 1000, 20% reserve, reset time) in daemon status --json and text; a warn usage part in health --json and text; exactly one usage-reserve item in inbox --json and text. With the cleared fixture it checks: no usage-reserve in daemon status or inbox (items []), and health usage ok with each window's use ('session 300, no limit known; weekly 0, no limit known', the weekly window being past its reset). describe-golden passes, and it changed against the merge-base 81f92766. Read the fix diff: healthUsage now reads the UsageWindow.Used the daemon poll saves (usageWindowsNow stores w.Used) and no longer calls usageWindowsNow, which only internal/api/daemon.go:163 calls now. go test ./internal/api -run 'Usage|Health|Inbox' -count=1: ok (TestUsageWindowsAnchored checks the saved session Used=400). internal/daemon: ok. Hands-on fix for the earlier major finding: strace -e openat of the built binary's rhiza health on the real repo opened no ~/.claude/projects transcript (its only .jsonl opens were .git/rhiza/runs logs). Back to back under the same load (loadavg ~21 on 4 cores), user CPU was 19.0-19.6 s against 19.0-20.1 s for a merge-base build, so the regression is gone; health's cost is the model fold. The whole internal/api package hit the default 10m go test timeout at that load, in the unrelated TestLinkEndpointKinds.

## Links

### Out

- verifies [usage-hold-reported](../capability/rhiza/orchestration/usage-governor/usage-hold-reported.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
