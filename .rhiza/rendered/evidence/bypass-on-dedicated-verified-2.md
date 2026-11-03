# bypass-on-dedicated-verified-2

evidence · `01M6H228TDREWZKSZT9YSW1S9M` · actor [claude-code-c00ed026](../actor/claude-code-c00ed026.md)

## result

PASS

## statement

Verifier, attempt 2, at be71e7b8. Ran go test ./cmd/rhiza -run 'TestScripts/(machine-trust-bypass|autonomous-permissions)$' -count=1: both PASS. machine-trust-bypass.txtar, with XDG_CONFIG_HOME/rhiza/machine.toml set to trust = "dedicated", runs session start --permission-mode bypassPermissions, then session run on the fake claude, and claude.log holds bypassPermissions. With no file, session start refuses with the error naming <xdg>/rhiza/machine.toml and trust = "dedicated", and there is no rex.log. After the file is changed to daily between start and run, session run refuses with the same naming and there is no claude.log. autonomous-permissions passes with no machine config. In the fix diff, bypassCheck now calls machine.Read on the real path; machine.Read("") returns daily, and the $XDG_CONFIG_HOME placeholder appears only in the message. This closes the confirmed quality/qa f1 (a file planted under the cwd turned bypass on). The new TestBypassCheckNoConfigPathIsDaily passes. gofmt, go vet and go test ./... pass, except the timing test internal/api TestReadBudget, which failed under load average 6.7 on 4 cores. The diff does not touch any read path.

## Links

### Out

- verifies [bypass-on-dedicated](../capability/rhiza/runtime/bypass-on-dedicated.md)

## History

- `01M6FY9EXET6R8QN6GPEJTCF7B` Machine trust: a machine config can declare trust = dedicated (this VM does), and on such a machine the daemon and sessions may run with permission mode bypassPermissions; today bypassRefused (internal/api/session.go:437, session_run.go:153) refuses bypass everywhere because no machine config exists
