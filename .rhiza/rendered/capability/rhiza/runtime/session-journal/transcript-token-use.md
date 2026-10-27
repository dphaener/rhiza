# transcript-token-use

capability · `01M5WRWXQ0MVD3SNCJ9GFXZ6KK` · parent [session-journal](../session-journal.md)

## acceptance

Go adapter tests read a transcript fixture with three assistant messages, one of them split over several lines that repeat its usage, two sidechain assistant lines of one helper message, plus a user line and a broken line: the session's totals equal the sum over the three messages counted once each, the helper totals equal the helper message's counts once; an empty transcript gives zeros.

## statement

An adapter reads a session's token use from its own transcript: for Claude Code, the input, output, cache-creation and cache-read tokens of its assistant messages, each message counted once however many transcript lines repeat it, with lines its helper agents left in the same transcript counted as a separate helper total. Helper agents' own transcripts are not counted, and unreadable lines are skipped, not fatal.

## Links

### Out

- decided-by [m3-scope](../../../../decision/m3-scope.md)
- decided-by [helper-token-total](../../../../decision/helper-token-total.md)

### In

- depends-on [session-token-use](session-token-use.md)
- verifies [transcript-token-use-tests](../../../../evidence/transcript-token-use-tests.md)
- verifies [pr-38-session-messages-and-token-use](../../../../evidence/pr-38-session-messages-and-token-use.md)

## History

- `01M5WPTWM0MPSBEHRYDMV9Q7PM` Session messages and token use in the journal: every message between a run's sessions is recorded with its sender and recipient through a SendMessage hook, and each session's token use is recorded.
