# machine-show

capability · `01M6G0FSSB90F7DHTX1VVZQ86N` · parent [runtime](../runtime.md)

## acceptance

A testscript runs rhiza machine show with no file (daily, no machine config), with trust = "dedicated" (dedicated and the path) and with a bad file (exit non-zero naming the file), and --json matches rhiza api call; describe --json lists no operation that writes the machine config

## statement

rhiza machine show (and its API op) prints the machine's trust and the config file it came from, or daily with no machine config, and writes nothing; rhiza has no operation that writes trust

## Links

### Out

- decided-by [trust-shown](../../../decision/trust-shown.md)
- decided-by [trust-hand-edited](../../../decision/trust-hand-edited.md)
- derived-from [machine-trust-plan-interview](../../../source/machine-trust-plan-interview.md)

### In

- verifies [machine-show-verified](../../../evidence/machine-show-verified.md)
- verifies [machine-show-qa-1](../../../evidence/machine-show-qa-1.md)
- verifies [machine-show-qa-2](../../../evidence/machine-show-qa-2.md)
- verifies [pr-212-machine-trust](../../../evidence/pr-212-machine-trust.md)

## History

- `01M6FY9EXET6R8QN6GPEJTCF7B` Machine trust: a machine config can declare trust = dedicated (this VM does), and on such a machine the daemon and sessions may run with permission mode bypassPermissions; today bypassRefused (internal/api/session.go:437, session_run.go:153) refuses bypass everywhere because no machine config exists
