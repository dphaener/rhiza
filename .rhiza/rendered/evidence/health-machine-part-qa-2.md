# health-machine-part-qa-2

evidence · `01M6H2AMW9NF3C16V6G7QPEVGG` · actor [claude-code-f2af322a](../actor/claude-code-f2af322a.md)

## result

PASS

## statement

At b03a5a16, built with go build -o $TMPDIR/rhiza-x ./cmd/rhiza and run by hand in scratch repos with fake rex and claude on PATH and HOME/XDG_CONFIG_HOME pointed at scratch dirs: rhiza health --json shows the machine part ok 'trust daily (no machine config)' with no file and with HOME and XDG_CONFIG_HOME unset, human health shows 'ok machine trust dedicated (from <path>)' with trust = "dedicated", and with trust = "root" the part is warn 'could not be read: machine config <path>: trust = "root": want "dedicated" or "daily"'.

## Links

### Out

- verifies [health-machine-part](../capability/rhiza/runtime/health-machine-part.md)

## History

- `01M6FY9EXET6R8QN6GPEJTCF7B` Machine trust: a machine config can declare trust = dedicated (this VM does), and on such a machine the daemon and sessions may run with permission mode bypassPermissions; today bypassRefused (internal/api/session.go:437, session_run.go:153) refuses bypass everywhere because no machine config exists
