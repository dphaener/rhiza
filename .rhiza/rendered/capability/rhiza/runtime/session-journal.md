# session-journal

capability · `01M5WRWWARC3CG3GN0Y485B6NJ` · parent [runtime](../runtime.md)

## acceptance

A testscript runs a workflow on a fake adapter whose agent sends one SendMessage and leaves a transcript with token usage: the run's journal has a message entry and a tokens entry for the session, run status --json lists the message and the session's token use, and session show --json gives the same token use.

## statement

A run's journal records every message its sessions send each other through Claude Code's SendMessage, with sender, recipient and text, and each of its sessions' token use, and rhiza run status and rhiza session show show both.

## Links

### Out

- decided-by [sessions-talk-freely](../../../decision/sessions-talk-freely.md)
- decided-by [m3-scope](../../../decision/m3-scope.md)

### In

- verifies [session-journal-tests](../../../evidence/session-journal-tests.md)
- verifies [pr-38-session-messages-and-token-use](../../../evidence/pr-38-session-messages-and-token-use.md)

## History

- `01M5WPTWM0MPSBEHRYDMV9Q7PM` Session messages and token use in the journal: every message between a run's sessions is recorded with its sender and recipient through a SendMessage hook, and each session's token use is recorded.
