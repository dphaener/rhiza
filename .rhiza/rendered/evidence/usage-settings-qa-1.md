# usage-settings-qa-1

evidence · `01M6JDRT41DMZF4E8E6GTRZFF4` · actor [claude-code-590230ff](../actor/claude-code-590230ff.md)

## result

PASS

## statement

QA with a binary built from 0d62831c in a scratch repo: rhiza daemon settings lists rhiza.daemon.usage.session and .weekly unset and .reserve 20 by default; daemon set took reserve 25 and session/weekly token counts; reserve 101, session 1.5 and weekly abc each exited 2 with a clear message; settings --json shows the git config source.

## Links

### Out

- verifies [usage-settings](../capability/rhiza/orchestration/usage-governor/usage-settings.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
