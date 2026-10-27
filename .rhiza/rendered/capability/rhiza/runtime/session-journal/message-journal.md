# message-journal

capability · `01M5WRWX0ZGMQPHT667PFCEY26` · parent [session-journal](../session-journal.md)

## acceptance

A testscript launches a session in a run with a fake Rex, then pipes a SendMessage PostToolUse payload naming that session into rhiza session hook PostToolUse: the run's journal gains one message entry with the step, the session ID, the session's actor, the recipient and the text, and the session's recorded status is the same as for any PostToolUse. The same payload for a session outside any run, and for a run whose journal is unwritable, exits 0 and adds no entry; the journal's existing entries still fold. A structured message is journaled marked structured.

## statement

When a session that belongs to a run sends a message, rhiza session hook appends a message entry to that run's journal: the step, the sending session, its actor as sender, the recipient, the text, whether it is structured, and the summary, at the time of the hook, even after the session's exit is recorded. A session that belongs to no run journals nothing, and a failure to journal is reported on stderr and never fails the agent's tool call.

## Links

### Out

- depends-on [message-hook-payload](message-hook-payload.md)
- decided-by [sessions-talk-freely](../../../../decision/sessions-talk-freely.md)
- decided-by [structured-messages-journaled](../../../../decision/structured-messages-journaled.md)
- decided-by [message-summary-journaled](../../../../decision/message-summary-journaled.md)

### In

- depends-on [run-status-messages](run-status-messages.md)
- relates-to [message-recipient-session](../../../../question/message-recipient-session.md) — message-journal records the recipient exactly as SendMessage gives it; resolving it to a session is open.
- relates-to [session-message-journaled](../../../../question/session-message-journaled.md) — The intent journals messages through a SendMessage hook; whether rhiza session message is journaled too is open.
- verifies [message-journal-tests](../../../../evidence/message-journal-tests.md)
- verifies [pr-38-session-messages-and-token-use](../../../../evidence/pr-38-session-messages-and-token-use.md)

## History

- `01M5WPTWM0MPSBEHRYDMV9Q7PM` Session messages and token use in the journal: every message between a run's sessions is recorded with its sender and recipient through a SendMessage hook, and each session's token use is recorded.
