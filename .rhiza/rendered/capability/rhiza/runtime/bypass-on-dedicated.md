# bypass-on-dedicated

capability · `01M6G0F9H1ZR0W65DV8FZHZK47` · parent [runtime](../runtime.md)

## acceptance

A testscript with trust = "dedicated" in $XDG_CONFIG_HOME/rhiza/machine.toml starts a session and a run with --permission-mode bypassPermissions on a fake adapter and the launch logs bypassPermissions; with no file, and again after the file is changed to trust = "daily" between session start and session run, both refuse with the error naming the file and nothing is launched; the existing autonomous-permissions refusal still passes with no machine config

## statement

On a machine whose config has trust = "dedicated", session start, session run, run start and resume, and the daemon accept a permission mode that bypasses the agent's checks; on any other machine they refuse it, the error naming the machine config file and the line trust = "dedicated", and trust is read fresh at each check

## Links

### Out

- decided-by [trust-lifts-refusal-only](../../../decision/trust-lifts-refusal-only.md)
- decided-by [trust-read-live](../../../decision/trust-read-live.md)
- derived-from [machine-trust-plan-interview](../../../source/machine-trust-plan-interview.md)

### In

- verifies [bypass-on-dedicated-verified](../../../evidence/bypass-on-dedicated-verified.md)
- verifies [bypass-on-dedicated-qa-1](../../../evidence/bypass-on-dedicated-qa-1.md)
- verifies [bypass-on-dedicated-verified-2](../../../evidence/bypass-on-dedicated-verified-2.md)
- verifies [bypass-on-dedicated-qa-2](../../../evidence/bypass-on-dedicated-qa-2.md)
- verifies [pr-212-machine-trust](../../../evidence/pr-212-machine-trust.md)

## History

- `01M6FY9EXET6R8QN6GPEJTCF7B` Machine trust: a machine config can declare trust = dedicated (this VM does), and on such a machine the daemon and sessions may run with permission mode bypassPermissions; today bypassRefused (internal/api/session.go:437, session_run.go:153) refuses bypass everywhere because no machine config exists
