# usage-governor

capability · `01M5SWTC9RDZ6PDFHA1M25BGPK` · parent [orchestration](../orchestration.md)

## statement

The usage governor paces agents to keep a reserve of the owner's usage on every window.

## Links

### Out

- decided-by [usage-governor-hold-rules](../../../decision/usage-governor-hold-rules.md)

### In

- verifies [pr-238-usage-governor](../../../evidence/pr-238-usage-governor.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
