# usage-window-tokens-qa-2

evidence · `01M6JM1QBM5AP5Y9JH4E1C0F7Q` · actor [claude-code-3977a5fc](../actor/claude-code-3977a5fc.md)

## result

PASS

## statement

QA cycle 2, built binary at 46cc1d58 in a scratch repo on the VM with the real ~/.claude transcripts (Rex faked, starts captured by RHIZA_TEST_DETACH_LOG with RHIZA_TEST_NOW set to the wall clock): the daemon's session window use, 4798557 tokens for [2026-11-04T01:29:50.689Z, 06:25:58Z), matched exactly an independent Python count over every .jsonl under ~/.claude/projects (assistant lines, each message id once at its earliest timestamp, input+output+cache_creation, no cache reads).

## Links

### Out

- verifies [usage-window-tokens](../capability/rhiza/orchestration/usage-governor/usage-window-tokens.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
