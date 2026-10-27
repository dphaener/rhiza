# token-use-read-at-done-and-step-end

decision · `01M5X95435CPCXPF0W59M0E3A8`

## statement

Token use is read when a session says done or its step ends, and again on exit; the later reading wins.

## wording

the orchestrator's ruling as arbiter for the owner

## Links

### In

- decided-by [session-token-use](../capability/rhiza/runtime/session-journal/session-token-use.md)
- decided-by [run-status-tokens](../capability/rhiza/runtime/session-journal/run-status-tokens.md)

## History

- `01M5WPTWM0MPSBEHRYDMV9Q7PM` Session messages and token use in the journal: every message between a run's sessions is recorded with its sender and recipient through a SendMessage hook, and each session's token use is recorded.
