# trust-hand-edited

decision · `01M6G0BXZ346G3WKP358VWH6AB`

## statement

Trust is set only by hand-editing the machine config: rhiza has no command or API operation that writes it, so no agent raises trust through rhiza

## wording

(assumed for the owner by the orchestrator acting as their proxy, the agent's option picked by the orchestrator)

## Links

### Out

- derived-from [machine-trust-plan-interview](../source/machine-trust-plan-interview.md)

### In

- decided-by [machine-show](../capability/rhiza/runtime/machine-show.md)

## History

- `01M6FY9EXET6R8QN6GPEJTCF7B` Machine trust: a machine config can declare trust = dedicated (this VM does), and on such a machine the daemon and sessions may run with permission mode bypassPermissions; today bypassRefused (internal/api/session.go:437, session_run.go:153) refuses bypass everywhere because no machine config exists
