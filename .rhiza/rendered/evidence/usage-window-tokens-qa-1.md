# usage-window-tokens-qa-1

evidence · `01M6JDS024BBSD2X37WWMA68Z7` · actor [claude-code-590230ff](../actor/claude-code-590230ff.md)

## result

PASS

## statement

QA against this machine's real ~/.claude/projects (335 transcripts, 215MB): UsageRecords read 5897 records in about 6 s; rhiza health's usage part showed non-zero session and weekly use, and with HOME pointed at a missing dir it showed 0 for both. Correct but slow: every call re-reads and parses every transcript (see the qa finding).

## Links

### Out

- verifies [usage-window-tokens](../capability/rhiza/orchestration/usage-governor/usage-window-tokens.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
