# usage-window-tokens

capability · `01M6HJ0DA7EGP43KHWDZZA0ZP8` · parent [usage-governor](../usage-governor.md)

## acceptance

A test in internal/adapter/claudecode builds a fake config dir with transcripts in two project dirs, a message repeated across content-block lines, a helper's lines, cache-read tokens and records before and after the span, and gets the exact expected records, each with its timestamp and token count, in time order; no file outside the fake config dir is read

## statement

An adapter reports the token records between two times across every one of its transcripts on the machine, each with its timestamp and its counted tokens; the Claude Code adapter reads sessions and their helpers under <config>/projects, each message once by id, counting input, output and cache-creation tokens and never cache reads, and only records timestamped inside the span

## Links

### Out

- decided-by [usage-governor-hold-rules](../../../../decision/usage-governor-hold-rules.md)

### In

- depends-on [usage-windows](usage-windows.md)
- verifies [usage-window-tokens-verified](../../../../evidence/usage-window-tokens-verified.md)
- verifies [usage-window-tokens-qa-1](../../../../evidence/usage-window-tokens-qa-1.md)
- verifies [usage-window-tokens-verified-2](../../../../evidence/usage-window-tokens-verified-2.md)
- verifies [usage-window-tokens-qa-2](../../../../evidence/usage-window-tokens-qa-2.md)
- verifies [pr-238-usage-governor](../../../../evidence/pr-238-usage-governor.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
