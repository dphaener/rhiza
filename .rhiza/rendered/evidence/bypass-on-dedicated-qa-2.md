# bypass-on-dedicated-qa-2

evidence · `01M6H29QZ7K5C5MS4TT67EHVPE` · actor [claude-code-f2af322a](../actor/claude-code-f2af322a.md)

## result

PASS

## statement

At b03a5a16, built with go build -o $TMPDIR/rhiza-x ./cmd/rhiza and run by hand in scratch repos with fake rex and claude on PATH and HOME/XDG_CONFIG_HOME pointed at scratch dirs: with no file and with trust = "daily", session start --permission-mode bypassPermissions and run start --workflow (one agent step) --permission-mode bypassPermissions (with and without --autonomous) exit 1 with 'needs a machine with trust = "dedicated" in <xdg>/rhiza/machine.toml', and neither rex nor claude is invoked. With trust = "dedicated", session start succeeds, session run launches claude with --permission-mode bypassPermissions, and run start runs its agent step to succeeded with claude launched with bypassPermissions. After the file is switched to daily between session start and session run, session run is refused naming the file and claude is not launched. trust = "root" refuses bypass naming the file; acceptEdits still starts. The cycle-1 finding is closed: with HOME and XDG_CONFIG_HOME both unset and trust = "dedicated" planted at both ./$XDG_CONFIG_HOME/rhiza/machine.toml and ./.config/rhiza/machine.toml in the cwd, session start and session run both refuse bypass and nothing is launched. go test for machine-trust-bypass, autonomous-permissions and TestBypassCheckNoConfigPathIsDaily passes.

## Links

### Out

- verifies [bypass-on-dedicated](../capability/rhiza/runtime/bypass-on-dedicated.md)

## History

- `01M6FY9EXET6R8QN6GPEJTCF7B` Machine trust: a machine config can declare trust = dedicated (this VM does), and on such a machine the daemon and sessions may run with permission mode bypassPermissions; today bypassRefused (internal/api/session.go:437, session_run.go:153) refuses bypass everywhere because no machine config exists
