# usage-hold-reported-qa-1

evidence · `01M6JDSBK07XJMSZGW7XMFR802` · actor [claude-code-590230ff](../actor/claude-code-590230ff.md)

## result

PASS

## statement

QA with real daemon --once polls (no seeded state): with session limit 1000 the state saved usage_hold; daemon status showed the settings, health showed warn usage with the window/use/limit/reserve/reset detail and inbox listed one usage-reserve item. After raising the limit, the next --once cleared usage_hold, health's usage part was ok (session N of limit; weekly N, no limit known) and the inbox was empty. The per-change-set skip was not exercised by hand (no accepted change set in the scratch repo).

## Links

### Out

- verifies [usage-hold-reported](../capability/rhiza/orchestration/usage-governor/usage-hold-reported.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
