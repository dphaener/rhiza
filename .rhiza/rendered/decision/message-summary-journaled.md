# message-summary-journaled

decision · `01M5X9547A3G5CSZA4PRSAXHGV`

## statement

SendMessage's summary is journaled with the message and shown in run status.

## wording

the agent's option, picked by the orchestrator as arbiter for the owner

## Links

### In

- decided-by [message-hook-payload](../capability/rhiza/runtime/session-journal/message-hook-payload.md)
- decided-by [message-journal](../capability/rhiza/runtime/session-journal/message-journal.md)
- decided-by [run-status-messages](../capability/rhiza/runtime/session-journal/run-status-messages.md)

## History

- `01M5WPTWM0MPSBEHRYDMV9Q7PM` Session messages and token use in the journal: every message between a run's sessions is recorded with its sender and recipient through a SendMessage hook, and each session's token use is recorded.
