# message-journal-tests

evidence · `01M5XB6YZZFS61ME2YRCZBNV85` · actor [claude-code-d3710d3d](../actor/claude-code-d3710d3d.md)

## result

`go test ./cmd/rhiza -run 'TestScripts/session-hook-message$' -count=1` and `go test ./internal/api -run TestSessionHook -count=1` passed at 1466643

## statement

rhiza session hook journals the messages a run's sessions send.

## status

verified

## tests

TestScripts/session-hook-message, TestSessionHookJournalsAfterExit, TestSessionHookNoRunJournalsNothing

## Links

### Out

- verifies [message-journal](../capability/rhiza/runtime/session-journal/message-journal.md)

## History

- `01M5WPTWM0MPSBEHRYDMV9Q7PM` Session messages and token use in the journal: every message between a run's sessions is recorded with its sender and recipient through a SendMessage hook, and each session's token use is recorded.
