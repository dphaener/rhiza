# trust-lifts-refusal-only

decision · `01M6G0CREKTXG4ZVZHTP1VVRW8`

## statement

Trust changes no default permission mode: on a dedicated machine it only stops refusing a mode that bypasses the agent's checks, and the mode stays explicit

## wording

(assumed for the owner by the orchestrator acting as their proxy, the agent's option picked by the orchestrator)

## Links

### Out

- derived-from [machine-trust-plan-interview](../source/machine-trust-plan-interview.md)

### In

- decided-by [bypass-on-dedicated](../capability/rhiza/runtime/bypass-on-dedicated.md)

## History

- `01M6FY9EXET6R8QN6GPEJTCF7B` Machine trust: a machine config can declare trust = dedicated (this VM does), and on such a machine the daemon and sessions may run with permission mode bypassPermissions; today bypassRefused (internal/api/session.go:437, session_run.go:153) refuses bypass everywhere because no machine config exists
