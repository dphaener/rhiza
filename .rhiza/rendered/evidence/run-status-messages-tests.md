# run-status-messages-tests

evidence · `01M5XB6ZH4WJFJ1ZH1KSXQZ5XF` · actor [claude-code-d3710d3d](../actor/claude-code-d3710d3d.md)

## result

`go test ./cmd/rhiza -run 'TestScripts/run-status-messages$' -count=1` passed at 1466643

## statement

run status lists the messages a run's sessions sent.

## status

verified

## tests

TestScripts/run-status-messages

## Links

### Out

- verifies [run-status-messages](../capability/rhiza/runtime/session-journal/run-status-messages.md)

## History

- `01M5WPTWM0MPSBEHRYDMV9Q7PM` Session messages and token use in the journal: every message between a run's sessions is recorded with its sender and recipient through a SendMessage hook, and each session's token use is recorded.
