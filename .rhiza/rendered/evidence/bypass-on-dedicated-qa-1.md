# bypass-on-dedicated-qa-1

evidence · `01M6GXPXH0546SAW6RTT9GYS6S` · actor [claude-code-fd06db9c](../actor/claude-code-fd06db9c.md)

## result

PASS

## statement

At 9560eed9, built with go build -o $TMPDIR/rhiza-x ./cmd/rhiza and run by hand in a scratch repo with fake rex and claude on PATH and HOME/XDG_CONFIG_HOME pointed at scratch dirs: with no file and with trust = "daily", session start --permission-mode bypassPermissions and run start --workflow passes --permission-mode bypassPermissions exit 1 with 'needs a machine with trust = "dedicated" in <path>/machine.toml' and nothing reaches rex or claude. With trust = "dedicated", session start succeeds, session run launches claude with --permission-mode bypassPermissions, and run start runs its agent step to completion with claude launched with bypassPermissions. After switching the file to daily between session start and session run, session run is refused naming the file and claude is not launched. trust = "root" refuses bypass naming the file; acceptEdits still starts with a bad file. Separate finding: with HOME and XDG_CONFIG_HOME both unset, a file at the relative path '$XDG_CONFIG_HOME/rhiza/machine.toml' in the cwd grants bypass.

## Links

### Out

- verifies [bypass-on-dedicated](../capability/rhiza/runtime/bypass-on-dedicated.md)

## History

- `01M6FY9EXET6R8QN6GPEJTCF7B` Machine trust: a machine config can declare trust = dedicated (this VM does), and on such a machine the daemon and sessions may run with permission mode bypassPermissions; today bypassRefused (internal/api/session.go:437, session_run.go:153) refuses bypass everywhere because no machine config exists
