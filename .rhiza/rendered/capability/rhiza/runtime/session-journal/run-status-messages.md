# run-status-messages

capability · `01M5WRWXC0MEDRAPVSV5P23ZVK` · parent [session-journal](../session-journal.md)

## acceptance

A testscript journals two messages in a run through rhiza session hook: run status --json has a messages list with both in order and every field, the markdown shows each with its sender, recipient and text, and a run with no messages has an empty list. describe-golden is regenerated for the new response field.

## statement

rhiza run status lists a run's journaled messages in order, each with its time, step, sending session, sender, recipient, text, summary and whether it is structured, so the conversation between its sessions can be read back after the run.

## Links

### Out

- depends-on [message-journal](message-journal.md)
- decided-by [sessions-talk-freely](../../../../decision/sessions-talk-freely.md)
- decided-by [message-summary-journaled](../../../../decision/message-summary-journaled.md)

### In

- verifies [run-status-messages-tests](../../../../evidence/run-status-messages-tests.md)
- verifies [pr-38-session-messages-and-token-use](../../../../evidence/pr-38-session-messages-and-token-use.md)

## History

- `01M5WPTWM0MPSBEHRYDMV9Q7PM` Session messages and token use in the journal: every message between a run's sessions is recorded with its sender and recipient through a SendMessage hook, and each session's token use is recorded.
