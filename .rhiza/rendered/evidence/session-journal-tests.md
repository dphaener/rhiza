# session-journal-tests

evidence · `01M5XB701M3P4WTVG3ZNXFJ6GY` · actor [claude-code-d3710d3d](../actor/claude-code-d3710d3d.md)

## result

`go test ./cmd/rhiza -run 'TestScripts/session-journal$' -count=1` passed at 1466643

## statement

The session journal records what a run's sessions did.

## status

verified

## tests

TestScripts/session-journal

## Links

### Out

- verifies [session-journal](../capability/rhiza/runtime/session-journal.md)

## History

- `01M5WPTWM0MPSBEHRYDMV9Q7PM` Session messages and token use in the journal: every message between a run's sessions is recorded with its sender and recipient through a SendMessage hook, and each session's token use is recorded.
