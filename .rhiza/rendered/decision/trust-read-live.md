# trust-read-live

decision · `01M6G0DRSM17ATXZV5AQV959V4`

## statement

Trust is read fresh from the machine config at every check (session start, session run, run start and resume, daemon start), and no session record or journal carries it

## wording

(assumed for the owner by the orchestrator acting as their proxy, the agent's option picked by the orchestrator)

## Links

### Out

- derived-from [machine-trust-plan-interview](../source/machine-trust-plan-interview.md)

### In

- decided-by [bypass-on-dedicated](../capability/rhiza/runtime/bypass-on-dedicated.md)

## History

- `01M6FY9EXET6R8QN6GPEJTCF7B` Machine trust: a machine config can declare trust = dedicated (this VM does), and on such a machine the daemon and sessions may run with permission mode bypassPermissions; today bypassRefused (internal/api/session.go:437, session_run.go:153) refuses bypass everywhere because no machine config exists
