# machine-config-verified

evidence · `01M6GQ7B8GCWEQY57G29N71172` · actor [claude-code-876bd2f5](../actor/claude-code-876bd2f5.md)

## result

PASS

## statement

At 95a41720 ran go test -count=1 -v ./internal/machine: TestAbsentIsDaily, TestTrustValues, TestHomeDefault, TestBadFilesNameThePath, TestPathUnknown all PASS. Read the diff: Path uses $XDG_CONFIG_HOME/rhiza/machine.toml when set, else ~/.config/rhiza/machine.toml; Read returns daily for an absent file, accepts trust = "dedicated" and "daily", and errors naming the path for non-TOML input and for trust = "root". The tests cover no file -> daily, both trust values via XDG_CONFIG_HOME, dedicated via the ~/.config default, and both bad files with the path in the error. make check passed fmt-check, vet and go test ./...; self-check failed only on stale .rhiza/rendered actor/index files, which is expected before render on landing.

## Links

### Out

- verifies [machine-config](../capability/rhiza/runtime/machine-config.md)

## History

- `01M6FY9EXET6R8QN6GPEJTCF7B` Machine trust: a machine config can declare trust = dedicated (this VM does), and on such a machine the daemon and sessions may run with permission mode bypassPermissions; today bypassRefused (internal/api/session.go:437, session_run.go:153) refuses bypass everywhere because no machine config exists
