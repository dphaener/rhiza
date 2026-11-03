# machine-config-file

decision · `01M6G0BMCXN7C1SCFGZWAWV5RW`

## statement

A machine's trust lives in a machine config file, $XDG_CONFIG_HOME/rhiza/machine.toml (default ~/.config/rhiza/machine.toml), as trust = "dedicated" or "daily"; it is outside every repo, so no clone carries it

## wording

(assumed for the owner by the orchestrator acting as their proxy, the agent's option picked by the orchestrator)

## Links

### Out

- derived-from [machine-trust-plan-interview](../source/machine-trust-plan-interview.md)

### In

- decided-by [machine-config](../capability/rhiza/runtime/machine-config.md)

## History

- `01M6FY9EXET6R8QN6GPEJTCF7B` Machine trust: a machine config can declare trust = dedicated (this VM does), and on such a machine the daemon and sessions may run with permission mode bypassPermissions; today bypassRefused (internal/api/session.go:437, session_run.go:153) refuses bypass everywhere because no machine config exists
