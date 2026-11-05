# usage-settings

capability · `01M6HJ0R8X5Q00XXMB8CRVN89Q` · parent [usage-governor](../usage-governor.md)

## acceptance

A new usage-settings.txtar shows the three keys with their defaults in rhiza daemon settings --json, sets each with rhiza daemon set, rejects a reserve outside 0 to 100 and a non-number limit with exit 2, and describe-golden.txtar is regenerated

## statement

Daemon settings rhiza.daemon.usage.session and rhiza.daemon.usage.weekly set each window's token limit (unset by default) and rhiza.daemon.usage.reserve sets the reserve as a percent of the limit (default 20), shown by rhiza daemon settings and written by rhiza daemon set

## Links

### Out

- decided-by [usage-governor-hold-rules](../../../../decision/usage-governor-hold-rules.md)

### In

- depends-on [usage-hold](usage-hold.md)
- verifies [usage-settings-verified](../../../../evidence/usage-settings-verified.md)
- verifies [usage-settings-qa-1](../../../../evidence/usage-settings-qa-1.md)
- verifies [usage-settings-qa-2](../../../../evidence/usage-settings-qa-2.md)
- verifies [pr-238-usage-governor](../../../../evidence/pr-238-usage-governor.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
