# printout-gap-names-reason

capability · `01M6505PVYH9A049T3Q3H2T28K` · parent [run-status-tokens](../run-status-tokens.md)

## acceptance

A testscript runs a session whose transcript file is missing and another still running; the printout and run status --json name each reason, the totals leave both out, and describe-golden is regenerated only if a response field changed

## statement

Where a session has no readable token use, the final printout and run status name the reason (still running, or transcript unreadable with its warning) instead of 'no token use recorded', and leave the session out of the totals

## Links

### Out

- decided-by [printout-reads-transcript-when-unjournaled](../../../../../decision/printout-reads-transcript-when-unjournaled.md)

### In

- verifies [printout-gap-names-reason-verified](../../../../../evidence/printout-gap-names-reason-verified.md)
- verifies [printout-gap-names-reason-qa-1](../../../../../evidence/printout-gap-names-reason-qa-1.md)
- verifies [pr-102-run-token-printout](../../../../../evidence/pr-102-run-token-printout.md)

## History

- `01M61GW8RKSBM4BT1B8B3FR5JH` A run's final printout and run status show each session's token use, never 'no token use recorded' for sessions that used tokens.
