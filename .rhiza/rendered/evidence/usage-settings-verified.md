# usage-settings-verified

evidence · `01M6J27B7BHBFNKZC78J3PV7JB` · actor [claude-code-b1d1957d](../actor/claude-code-b1d1957d.md)

## result

PASS

## statement

At 86877c2f: go test ./cmd/rhiza -run 'TestScripts/(usage-settings|describe-golden|daemon)' passes; usage-settings.txtar covers the three keys' defaults, rhiza daemon set for each, and exit 2 with nothing written for reserve 101, -1 and lots and for session many and weekly 1.5. By hand with a scratch build in a fresh repo: daemon settings --json shows usage.session and usage.weekly at default "" and usage.reserve at default 20; reserve 101 and lots, and weekly many, 1.5 and 0, each exit 2 and leave git config unset. describe-golden.txtar carries the updated daemon.set desc. make fmt-check, vet and self-check pass, so attempt 1's finding (stale rendered actor files) is fixed. make check's test step failed only TestReadBudget (internal/api read-timing medians over 50ms) on a 4-core VM at load 11 to 22; the same test fails the same way at the leaf's base b67b44b7, and the leaf touches no read path. Every other package passes.

## Links

### Out

- verifies [usage-settings](../capability/rhiza/orchestration/usage-governor/usage-settings.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
