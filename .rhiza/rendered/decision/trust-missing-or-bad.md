# trust-missing-or-bad

decision · `01M6G0CADH4TRB0T2ZN4DB8J7F`

## statement

No machine config means trust = daily; a malformed file or an unknown trust value is an error naming the file, and bypass stays refused

## wording

(assumed for the owner by the orchestrator acting as their proxy, the agent's option picked by the orchestrator)

## Links

### Out

- derived-from [machine-trust-plan-interview](../source/machine-trust-plan-interview.md)

### In

- decided-by [machine-config](../capability/rhiza/runtime/machine-config.md)
- decided-by [health-machine-part](../capability/rhiza/runtime/health-machine-part.md)

## History

- `01M6FY9EXET6R8QN6GPEJTCF7B` Machine trust: a machine config can declare trust = dedicated (this VM does), and on such a machine the daemon and sessions may run with permission mode bypassPermissions; today bypassRefused (internal/api/session.go:437, session_run.go:153) refuses bypass everywhere because no machine config exists
