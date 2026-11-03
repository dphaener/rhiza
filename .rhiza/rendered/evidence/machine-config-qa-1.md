# machine-config-qa-1

evidence · `01M6GXP0N4HXAZAKSCPX3PFHWV` · actor [claude-code-fd06db9c](../actor/claude-code-fd06db9c.md)

## result

PASS

## statement

At 9560eed9, built with go build -o $TMPDIR/rhiza-x ./cmd/rhiza and run by hand in a scratch repo with fake rex and claude on PATH and HOME/XDG_CONFIG_HOME pointed at scratch dirs: with no machine.toml, rhiza machine show says trust daily (no machine config) naming $XDG_CONFIG_HOME/rhiza/machine.toml; trust = "dedicated" reads as dedicated; with XDG_CONFIG_HOME unset it reads ~/.config/rhiza/machine.toml (daily when absent, dedicated once written); trust = "root" and 'not = [toml' both exit 1 with an error naming the file's path.

## Links

### Out

- verifies [machine-config](../capability/rhiza/runtime/machine-config.md)

## History

- `01M6FY9EXET6R8QN6GPEJTCF7B` Machine trust: a machine config can declare trust = dedicated (this VM does), and on such a machine the daemon and sessions may run with permission mode bypassPermissions; today bypassRefused (internal/api/session.go:437, session_run.go:153) refuses bypass everywhere because no machine config exists
