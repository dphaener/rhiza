# transcript-token-use-tests

evidence · `01M5XB6YQ0KDBBSMAAJANPRSGE` · actor [claude-code-d3710d3d](../actor/claude-code-d3710d3d.md)

## result

`go test ./internal/adapter/claudecode -run TokenUse -count=1` passed at 1466643

## statement

The Claude Code adapter sums token use from a session transcript, tolerating empty, broken and very long lines.

## status

verified

## tests

TestTokenUse, TestTokenUseEmpty, TestTokenUseBroken, TestTokenUseLongLine

## Links

### Out

- verifies [transcript-token-use](../capability/rhiza/runtime/session-journal/transcript-token-use.md)

## History

- `01M5WPTWM0MPSBEHRYDMV9Q7PM` Session messages and token use in the journal: every message between a run's sessions is recorded with its sender and recipient through a SendMessage hook, and each session's token use is recorded.
