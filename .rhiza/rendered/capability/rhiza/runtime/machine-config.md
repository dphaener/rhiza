# machine-config

capability · `01M6G0E9J86Z2AE53V9R82JXD9` · parent [runtime](../runtime.md)

## acceptance

Go tests over the machine config reader: no file gives daily; trust = "dedicated" and trust = "daily" read as such, from XDG_CONFIG_HOME when set and ~/.config otherwise; a file that is not TOML, or trust = "root", is an error naming the file's path

## statement

Rhiza reads a per-machine config at $XDG_CONFIG_HOME/rhiza/machine.toml (default ~/.config/rhiza/machine.toml) whose trust is "dedicated" or "daily", daily when the file is absent, and a malformed file or unknown trust is an error naming the file

## Links

### Out

- decided-by [machine-config-file](../../../decision/machine-config-file.md)
- decided-by [trust-missing-or-bad](../../../decision/trust-missing-or-bad.md)
- derived-from [machine-trust-plan-interview](../../../source/machine-trust-plan-interview.md)

### In

- verifies [machine-config-verified](../../../evidence/machine-config-verified.md)
- verifies [machine-config-qa-1](../../../evidence/machine-config-qa-1.md)
- verifies [pr-212-machine-trust](../../../evidence/pr-212-machine-trust.md)

## History

- `01M6FY9EXET6R8QN6GPEJTCF7B` Machine trust: a machine config can declare trust = dedicated (this VM does), and on such a machine the daemon and sessions may run with permission mode bypassPermissions; today bypassRefused (internal/api/session.go:437, session_run.go:153) refuses bypass everywhere because no machine config exists
