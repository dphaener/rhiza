# machine-show-verified

evidence · `01M6GV5M1R409CTFCX5MA2WCAQ` · actor [claude-code-b907b60e](../actor/claude-code-b907b60e.md)

## result

PASS

## statement

At 4a8bb2ef ran go test ./cmd/rhiza -run 'TestScripts/(machine-show|describe-golden)$' -count=1: ok. machine-show.txtar runs rhiza machine show with no file (trust: daily, no machine config, path shown, file false), with trust = "dedicated" (trust: dedicated from <xdg>/rhiza/machine.toml, file true), and with trust = "root" (non-zero, stderr and --json stdout name the file); --json is cmp-equal to rhiza api call machine.show in both good cases. describe-golden lists machine.show with write false and no other machine op. Read the diff: internal/api/machine.go only reads via machine.Read and stats the path; grep found no code writing machine.toml. go vet and fmt-check pass. Full go test ./... failed only internal/tui/modelapp's read-timing budget (untouched package, runner contention); make check self-check reports only stale .rhiza/rendered files, which landing renders.

## Links

### Out

- verifies [machine-show](../capability/rhiza/runtime/machine-show.md)

## History

- `01M6FY9EXET6R8QN6GPEJTCF7B` Machine trust: a machine config can declare trust = dedicated (this VM does), and on such a machine the daemon and sessions may run with permission mode bypassPermissions; today bypassRefused (internal/api/session.go:437, session_run.go:153) refuses bypass everywhere because no machine config exists
