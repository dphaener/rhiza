# bypass-on-dedicated-verified

evidence · `01M6GVCVZ9JNXA28HD0FJGRHEZ` · actor [claude-code-c7f28746](../actor/claude-code-c7f28746.md)

## result

PASS

## statement

At 5c5bf355 ran go test ./cmd/rhiza -run 'TestScripts/(machine-trust-bypass|autonomous-permissions|session-start|daemon-picks)$' -count=1: all PASS. machine-trust-bypass.txtar, using fake rex and claude binaries: with no $XDG_CONFIG_HOME/rhiza/machine.toml, session start --permission-mode bypassPermissions is refused, stderr names trust = "dedicated" in <path>/machine.toml, and no rex.log is written. With trust = "dedicated", session start succeeds and session run launches claude, whose log shows bypassPermissions. After the file changes to trust = "daily" between start and run, session run is refused naming the file and no claude.log is written. An invalid trust is refused naming the file. autonomous-permissions still passes its bypass refusal with no machine config. In the diff, startAdapter and session run call bypassCheck, which calls machine.Read at each check, and checkAgentSteps passes ctx through, so run start, resume, signoff and the daemon get the same check. make check ran go test with no failures and failed only on stale .rhiza/rendered files from this run's own model writes. A second full run under load average ~24 hit timeouts in plan-workflow and plan-workflow-pr and budget-test failures; rerun alone, both plan-workflow scripts PASS.

## Links

### Out

- verifies [bypass-on-dedicated](../capability/rhiza/runtime/bypass-on-dedicated.md)

## History

- `01M6FY9EXET6R8QN6GPEJTCF7B` Machine trust: a machine config can declare trust = dedicated (this VM does), and on such a machine the daemon and sessions may run with permission mode bypassPermissions; today bypassRefused (internal/api/session.go:437, session_run.go:153) refuses bypass everywhere because no machine config exists
