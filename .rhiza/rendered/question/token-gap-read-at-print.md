# token-gap-read-at-print

question · `01M6505MMR01DJT7JVD7XBZA79`

## statement

When a run's session has no journaled token reading yet (still running, or its transcript unreadable), should the final printout and run status read its transcript at print time? Recommended: yes, falling back to the transcript, and name the reason for any residual gap. Alternative: journal only, and name the reason instead of 'no token use recorded'.

## Links

### Out

- derived-from [run-token-printout-autonomous-plan](../source/run-token-printout-autonomous-plan.md)

### In

- answers [printout-reads-transcript-when-unjournaled](../decision/printout-reads-transcript-when-unjournaled.md)

## History

- `01M61GW8RKSBM4BT1B8B3FR5JH` A run's final printout and run status show each session's token use, never 'no token use recorded' for sessions that used tokens.
