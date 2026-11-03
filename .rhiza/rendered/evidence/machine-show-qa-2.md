# machine-show-qa-2

evidence · `01M6H2A4CCSAVH2S59848XS9ZD` · actor [claude-code-f2af322a](../actor/claude-code-f2af322a.md)

## result

PASS

## statement

At b03a5a16, built with go build -o $TMPDIR/rhiza-x ./cmd/rhiza and run by hand in scratch repos with fake rex and claude on PATH and HOME/XDG_CONFIG_HOME pointed at scratch dirs: rhiza machine show with no file prints 'trust: daily (no machine config; it would be <xdg>/rhiza/machine.toml)', and machine show --json is byte-identical to rhiza api call machine.show; with trust = "dedicated" it prints 'trust: dedicated (from <path>)'; with trust = "root" it exits 1 naming the file. With XDG_CONFIG_HOME empty it falls back to <home>/.config/rhiza/machine.toml. With HOME and XDG_CONFIG_HOME both unset it prints daily but an empty path ('it would be )').

## Links

### Out

- verifies [machine-show](../capability/rhiza/runtime/machine-show.md)

## History

- `01M6FY9EXET6R8QN6GPEJTCF7B` Machine trust: a machine config can declare trust = dedicated (this VM does), and on such a machine the daemon and sessions may run with permission mode bypassPermissions; today bypassRefused (internal/api/session.go:437, session_run.go:153) refuses bypass everywhere because no machine config exists
