# usage-learned-limit-qa-1

evidence · `01M6JDS5MGHVWW03HGQQVDCT2Y` · actor [claude-code-590230ff](../actor/claude-code-590230ff.md)

## result

PASS

## statement

QA: seeded .git/rhiza/usage.limits.jsonl with a hit resetting 2h later, one resetting 4 days later, one with no reset and a garbage line, then ran rhiza daemon --once with no limit settings over real transcripts. The state's learned list had session 2772142 and weekly 16403273 tokens at those hits; the no-reset hit and the bad line taught nothing. The learned session limit held (reset moved to the hit's 06:39:21Z). Setting rhiza.daemon.usage.session and .weekly to 1e9 won over the learned limits and cleared the hold.

## Links

### Out

- verifies [usage-learned-limit](../capability/rhiza/orchestration/usage-governor/usage-learned-limit.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
