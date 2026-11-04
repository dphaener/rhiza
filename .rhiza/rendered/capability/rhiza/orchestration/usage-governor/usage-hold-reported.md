# usage-hold-reported

capability · `01M6HJ14J13C9N67W124QEQME1` · parent [usage-governor](../usage-governor.md)

## acceptance

A new usage-hold.txtar with a seeded daemon state that holds shows the usage-reserve skip and its detail in daemon status --json and text, a warn usage part in health --json, and one usage-reserve item in inbox --json; with the hold cleared none of them shows and health's usage part is ok; describe-golden.txtar is regenerated

## statement

While the daemon holds for usage, each held change set is a skip with reason usage-reserve whose detail gives the window, use, limit and reset time in rhiza daemon status; rhiza health has a usage part (warn while held, ok with each window's use otherwise) and rhiza inbox lists one usage-reserve item until the hold clears

## Links

### Out

- depends-on [usage-hold](usage-hold.md)
- decided-by [usage-governor-hold-rules](../../../../decision/usage-governor-hold-rules.md)

### In

- verifies [usage-hold-reported-verified](../../../../evidence/usage-hold-reported-verified.md)
- verifies [usage-hold-reported-qa-1](../../../../evidence/usage-hold-reported-qa-1.md)
- verifies [usage-hold-reported-qa-2](../../../../evidence/usage-hold-reported-qa-2.md)
- verifies [usage-hold-reported-verified-5](../../../../evidence/usage-hold-reported-verified-5.md)
- verifies [usage-hold-reported-qa-3](../../../../evidence/usage-hold-reported-qa-3.md)
- verifies [pr-238-usage-governor](../../../../evidence/pr-238-usage-governor.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
