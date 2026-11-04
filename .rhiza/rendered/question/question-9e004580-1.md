# question-9e004580-1

question · `01M6HHSVH0RABRX630W06KPEZV`

## session

9e004580-c5cb-495c-8f99-01a1a31a84fb

## statement

Usage governor: when does the daemon hold a new pick? / Usage governor: which tokens count toward a window's use? / Usage governor: what does a usage hold stop?

## Links

### In

- answers [usage-governor-hold-rules](../decision/usage-governor-hold-rules.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
