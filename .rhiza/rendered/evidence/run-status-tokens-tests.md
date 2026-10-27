# run-status-tokens-tests

evidence · `01M5XB6ZSA4A4XDQVDBJ5QACJX` · actor [claude-code-d3710d3d](../actor/claude-code-d3710d3d.md)

## result

`go test ./cmd/rhiza -run 'TestScripts/run-status-tokens$' -count=1` passed at 1466643

## statement

run status shows each session's token use and the run's totals.

## status

verified

## tests

TestScripts/run-status-tokens

## Links

### Out

- verifies [run-status-tokens](../capability/rhiza/runtime/session-journal/run-status-tokens.md)

## History

- `01M5WPTWM0MPSBEHRYDMV9Q7PM` Session messages and token use in the journal: every message between a run's sessions is recorded with its sender and recipient through a SendMessage hook, and each session's token use is recorded.
