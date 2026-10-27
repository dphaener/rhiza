# run-status-tokens

capability · `01M5WRWY9E4NFEB9ZYE7K54AXK` · parent [session-journal](../session-journal.md)

## acceptance

A testscript runs a two-step workflow whose two fake agents leave transcripts with known usage: run status --json lists both sessions with their steps and counts and totals equal to their sums, the markdown shows the same, and a session whose token use was not recorded is listed with no counts and left out of the totals. describe-golden is regenerated for the new response field.

## statement

rhiza run status shows each of a run's sessions with its step, its journaled token use and its helpers', and the run's totals of each, the latest reading of each session counting, so the token cost of a run can be read when it is over.

## Links

### Out

- depends-on [session-token-use](session-token-use.md)
- decided-by [m3-scope](../../../../decision/m3-scope.md)
- decided-by [helper-token-total](../../../../decision/helper-token-total.md)
- decided-by [token-use-read-at-done-and-step-end](../../../../decision/token-use-read-at-done-and-step-end.md)

### In

- verifies [run-status-tokens-tests](../../../../evidence/run-status-tokens-tests.md)
- verifies [pr-38-session-messages-and-token-use](../../../../evidence/pr-38-session-messages-and-token-use.md)

## History

- `01M5WPTWM0MPSBEHRYDMV9Q7PM` Session messages and token use in the journal: every message between a run's sessions is recorded with its sender and recipient through a SendMessage hook, and each session's token use is recorded.
