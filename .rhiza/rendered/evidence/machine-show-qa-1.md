# machine-show-qa-1

evidence · `01M6GXRSP0VGP6C55DJE2S6C5R` · actor [claude-code-fd06db9c](../actor/claude-code-fd06db9c.md)

## result

PASS

## statement

At 9560eed9, built with go build -o $TMPDIR/rhiza-x ./cmd/rhiza and run by hand in a scratch repo with fake rex and claude on PATH and HOME/XDG_CONFIG_HOME pointed at scratch dirs: rhiza machine show prints 'trust: daily (no machine config; it would be <path>)' with no file and 'trust: dedicated (from <path>)' with one; a bad file exits 1 naming the file; rhiza machine show --json and rhiza api call machine.show give identical bytes (cmp); describe --json lists machine.show as the only machine operation and none writes trust.

## Links

### Out

- verifies [machine-show](../capability/rhiza/runtime/machine-show.md)

## History

- `01M6FY9EXET6R8QN6GPEJTCF7B` Machine trust: a machine config can declare trust = dedicated (this VM does), and on such a machine the daemon and sessions may run with permission mode bypassPermissions; today bypassRefused (internal/api/session.go:437, session_run.go:153) refuses bypass everywhere because no machine config exists
