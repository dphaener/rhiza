# usage-hold-qa-1

evidence · `01M6JDSHE8AW1GAQ7X6TSDPGG5` · actor [claude-code-590230ff](../actor/claude-code-590230ff.md)

## result

FAIL

## statement

QA: rhiza daemon --once with a tiny session limit held (started/resumed/planned empty, usage_hold saved), and with a large limit it did not. Over real transcripts, though, the current session window depends on the 2x-window lookback: at the same instant usage.Current gave start 23:38:27Z, 4564260 used, reset 04:38:27Z with a 10h lookback but start 03:48:10Z, 1280338 used, reset 08:48:10Z with a 30h one. Consecutive health calls showed session use drop with no reset, and the reported reset moved from 04:38:27Z to about 04:45Z ten seconds later. Live-run projection was not exercised by hand.

## Links

### Out

- verifies [usage-hold](../capability/rhiza/orchestration/usage-governor/usage-hold.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
