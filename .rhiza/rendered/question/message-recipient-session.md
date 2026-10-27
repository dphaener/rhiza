# message-recipient-session

question · `01M5WRWYMM148C2SMJSNEDB33D`

## statement

SendMessage names its recipient as Claude Code knows it, not as a Rhiza session ID. Should the journal also resolve the recipient to the run's session it reached, and how? (For example: each session's agent name set at launch from its step, so a name maps to one session.)

## Links

### Out

- relates-to [message-journal](../capability/rhiza/runtime/session-journal/message-journal.md) — message-journal records the recipient exactly as SendMessage gives it; resolving it to a session is open.

## History

- `01M5WPTWM0MPSBEHRYDMV9Q7PM` Session messages and token use in the journal: every message between a run's sessions is recorded with its sender and recipient through a SendMessage hook, and each session's token use is recorded.
