# session-token-use-tests

evidence · `01M5XB6Z8R44FKVBD35Z0Z36EP` · actor [claude-code-d3710d3d](../actor/claude-code-d3710d3d.md)

## result

`go test ./cmd/rhiza -run 'TestScripts/session-run-tokens$' -count=1` and `go test ./internal/run -run AgentEnded -count=1` passed at 1466643

## statement

A session's token use is recorded at done, step end and exit.

## status

verified

## tests

TestScripts/session-run-tokens, TestAgentEnded*

## Links

### Out

- verifies [session-token-use](../capability/rhiza/runtime/session-journal/session-token-use.md)

## History

- `01M5WPTWM0MPSBEHRYDMV9Q7PM` Session messages and token use in the journal: every message between a run's sessions is recorded with its sender and recipient through a SendMessage hook, and each session's token use is recorded.
