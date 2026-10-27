# message-hook-payload

capability · `01M5WRWWNXTXNFMZMGZFM46D6F` · parent [session-journal](../session-journal.md)

## acceptance

Go adapter tests feed PostToolUse payloads shaped like a real SendMessage call, with made-up content (one plain, one with newlines and quotes in its text, one whose message is a JSON object, one with an empty-string message), a PostToolUse payload for another tool, SendMessage inputs with no recipient and with no message, and each other hook event: SendMessage gives the recipient, the verbatim text and the summary, the object gives its raw JSON marked structured, the rest give no message, the malformed inputs are errors, and HookStatus returns the same report as before for all of them. The settings JSON still installs one PostToolUse hook with no matcher.

## statement

The Claude Code adapter reads a PostToolUse hook payload whose tool is SendMessage as a message: the recipient, the text exactly as the tool's input gives it (a message that is not a string kept as its raw JSON and marked structured) and its summary, if any. Input with no recipient or no message is malformed. Every other tool and every other event gives no message, and the status the hook reports is unchanged.

## Links

### Out

- decided-by [sessions-talk-freely](../../../../decision/sessions-talk-freely.md)
- decided-by [m3-scope](../../../../decision/m3-scope.md)
- decided-by [structured-messages-journaled](../../../../decision/structured-messages-journaled.md)
- decided-by [message-summary-journaled](../../../../decision/message-summary-journaled.md)

### In

- depends-on [message-journal](message-journal.md)
- verifies [message-hook-payload-tests](../../../../evidence/message-hook-payload-tests.md)
- verifies [pr-38-session-messages-and-token-use](../../../../evidence/pr-38-session-messages-and-token-use.md)

## History

- `01M5WPTWM0MPSBEHRYDMV9Q7PM` Session messages and token use in the journal: every message between a run's sessions is recorded with its sender and recipient through a SendMessage hook, and each session's token use is recorded.
