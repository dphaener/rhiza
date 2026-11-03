# health-machine-part-qa-1

evidence · `01M6GXQV6XVXSGV6BW53VKXEWA` · actor [claude-code-fd06db9c](../actor/claude-code-fd06db9c.md)

## result

PASS

## statement

At 9560eed9, built with go build -o $TMPDIR/rhiza-x ./cmd/rhiza and run by hand in a scratch repo with fake rex and claude on PATH and HOME/XDG_CONFIG_HOME pointed at scratch dirs: rhiza health --json puts machine between autonomy and disk; with no file it is ok 'trust daily (no machine config)', with trust = "dedicated" ok 'trust dedicated (from <path>)', with a malformed file warn 'could not be read: machine config <path>: ...'. Human rhiza health shows 'ok   machine  trust dedicated (from <path>)'.

## Links

### Out

- verifies [health-machine-part](../capability/rhiza/runtime/health-machine-part.md)

## History

- `01M6FY9EXET6R8QN6GPEJTCF7B` Machine trust: a machine config can declare trust = dedicated (this VM does), and on such a machine the daemon and sessions may run with permission mode bypassPermissions; today bypassRefused (internal/api/session.go:437, session_run.go:153) refuses bypass everywhere because no machine config exists
