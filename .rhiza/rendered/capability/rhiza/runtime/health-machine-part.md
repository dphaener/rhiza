# health-machine-part

capability · `01M6G0G7HP57PFFNW4WHZQJAA5` · parent [runtime](../runtime.md)

## acceptance

A testscript runs rhiza health --json with no machine config, with trust = "dedicated" and with a bad file, and finds a machine part that is ok with daily, ok with dedicated, and warn naming the file

## statement

rhiza health shows a machine part: ok with the trust and its file, or daily with no machine config, and warn naming the file when the machine config cannot be read

## Links

### Out

- decided-by [trust-shown](../../../decision/trust-shown.md)
- decided-by [trust-missing-or-bad](../../../decision/trust-missing-or-bad.md)
- derived-from [machine-trust-plan-interview](../../../source/machine-trust-plan-interview.md)

### In

- verifies [health-machine-part-verified](../../../evidence/health-machine-part-verified.md)
- verifies [health-machine-part-qa-1](../../../evidence/health-machine-part-qa-1.md)
- verifies [health-machine-part-qa-2](../../../evidence/health-machine-part-qa-2.md)
- verifies [pr-212-machine-trust](../../../evidence/pr-212-machine-trust.md)
- verifies [pr-301-delete-the-old-rex-apps](../../../evidence/pr-301-delete-the-old-rex-apps.md)

## History

- `01M6FY9EXET6R8QN6GPEJTCF7B` Machine trust: a machine config can declare trust = dedicated (this VM does), and on such a machine the daemon and sessions may run with permission mode bypassPermissions; today bypassRefused (internal/api/session.go:437, session_run.go:153) refuses bypass everywhere because no machine config exists
- `01M6MYNZJSMEXAE5Y3P3S9FB6F` Delete the old apps: remove changesapp, healthapp, modelapp, readapps, roadmapapp and writeapps with their Rex actions, keys and testscripts, and retire their capabilities; rhiza rex install registers exactly the five apps Inbox, Request, Runs, Daemon and Autonomy; the CLI keeps every command
