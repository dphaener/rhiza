# structured-messages-journaled

decision · `01M5X953Z19VAN334QMR0BAZ8Z`

## statement

A SendMessage whose message is not a string is journaled as its raw JSON marked structured; only input with no recipient or no message is malformed.

## wording

the orchestrator's ruling as arbiter for the owner

## Links

### In

- decided-by [message-hook-payload](../capability/rhiza/runtime/session-journal/message-hook-payload.md)
- decided-by [message-journal](../capability/rhiza/runtime/session-journal/message-journal.md)

## History

- `01M5WPTWM0MPSBEHRYDMV9Q7PM` Session messages and token use in the journal: every message between a run's sessions is recorded with its sender and recipient through a SendMessage hook, and each session's token use is recorded.
