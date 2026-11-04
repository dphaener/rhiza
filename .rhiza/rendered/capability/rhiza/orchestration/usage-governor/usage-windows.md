# usage-windows

capability · `01M6HJ0KXZCA4KV63M0R3FN1BW` · parent [usage-governor](../usage-governor.md)

## acceptance

Unit tests in a pure internal package, over timestamped token records as usage-window-tokens returns them and limit resets with fixed times, cover a first window, a record just after a window ended starting the next one, a gap with no records, and a usage-limit reset moving the start, for both window kinds

## statement

The governor tracks a 5-hour session window and a 7-day weekly window, each starting at the first token record after the previous window of its kind ended, or at the reset time a usage-limit message gave, and knows each window's start, reset and use so far

## Links

### Out

- depends-on [usage-window-tokens](usage-window-tokens.md)
- decided-by [usage-governor-hold-rules](../../../../decision/usage-governor-hold-rules.md)

### In

- depends-on [usage-hold](usage-hold.md)
- depends-on [usage-learned-limit](usage-learned-limit.md)
- verifies [usage-windows-verified](../../../../evidence/usage-windows-verified.md)
- verifies [usage-windows-qa-2](../../../../evidence/usage-windows-qa-2.md)
- verifies [pr-238-usage-governor](../../../../evidence/pr-238-usage-governor.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
