# health-machine-part-verified

evidence · `01M6GX9Y2JD9AYWDAZ0XPQ7C0H` · actor [claude-code-15906f2f](../actor/claude-code-15906f2f.md)

## result

PASS

## statement

At 0d70f9d2 ran go test ./cmd/rhiza -run 'TestScripts/(health|describe-golden)' -count=1: health-machine-part, health, tui-health-app, daemon-version-in-health and describe-golden pass. Read health-machine-part.txtar: it runs rhiza health --json with no machine config, with trust = "dedicated" and with a bad file (trust = "root"), and each --json check pins the machine part's own adjacent part/state/summary lines: ok with 'trust daily (no machine config)', ok with 'trust dedicated (from <file>)', and warn with 'could not be read: ...<file>'. This fixes both attempt-1 findings. Read the diff of internal/api/health.go: healthMachine sits between autonomy and disk, and the part enum and describe golden are updated. gofmt -l is clean and go vet passes. go test ./... passes in every package except internal/slice's read-budget timing medians, which come from runner contention in a package this leaf does not touch.

## Links

### Out

- verifies [health-machine-part](../capability/rhiza/runtime/health-machine-part.md)

## History

- `01M6FY9EXET6R8QN6GPEJTCF7B` Machine trust: a machine config can declare trust = dedicated (this VM does), and on such a machine the daemon and sessions may run with permission mode bypassPermissions; today bypassRefused (internal/api/session.go:437, session_run.go:153) refuses bypass everywhere because no machine config exists
