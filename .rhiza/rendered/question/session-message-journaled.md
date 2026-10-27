# session-message-journaled

question · `01M5WRWYW1NWMRQ5ZCP0THPAYM`

## statement

Should messages sent into a run's step with rhiza session message (step-messages), not through SendMessage, also be journaled as message entries? Today only an amendment is journaled.

## Links

### Out

- relates-to [message-journal](../capability/rhiza/runtime/session-journal/message-journal.md) — The intent journals messages through a SendMessage hook; whether rhiza session message is journaled too is open.

## History

- `01M5WPTWM0MPSBEHRYDMV9Q7PM` Session messages and token use in the journal: every message between a run's sessions is recorded with its sender and recipient through a SendMessage hook, and each session's token use is recorded.
