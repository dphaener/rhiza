# usage-window-tokens-verified-2

evidence · `01M6JFXYJMST1T13M2XW1MSNB9` · actor [claude-code-7c4785dd](../actor/claude-code-7c4785dd.md)

## result

PASS

## statement

At 7e926bfa (attempt 2), I ran go test -race -run Usage ./internal/adapter/claudecode and it passed. TestUsageRecords builds a fake CLAUDE_CONFIG_DIR with transcripts in two project dirs (-a, -b), message m1 repeated across lines, a helper transcript under s1/subagents, cache-read tokens, records before and after the span, a file outside the config dir and a symlink to it. It gets exactly 4 records in time order with the expected timestamps and tokens (7, 60, 300, 3): cache reads are not counted, the outside file is not read, and nothing outside the span is returned. TestUsageRecordsSkipsOldAndRereadsChanged also passes. I read the diff: files last written before from are skipped, and parsed records are cached by path, size and mtime. On this VM's real ~/.claude/projects (225 MB, 351 .jsonl files), I ran the reader through a go test overlay without editing the repo. The first call takes 2.4-4 s and later calls take 16-32 ms, so the QA finding f2 (a poll taking longer than the 5 s interval) is fixed. With an overlay of the pre-fix claudecode.go, the 5h and 7d spans ending 2026-11-04T04:00Z return identical records (same counts, same sha256). make fmt-check, vet and test all pass. make self-check fails only on stale rendered actor files from this run's records, not on code.

## Links

### Out

- verifies [usage-window-tokens](../capability/rhiza/orchestration/usage-governor/usage-window-tokens.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
