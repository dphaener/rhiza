# printout-counts-used-tokens

capability · `01M6505NQCFYY8Q7F22NEN8DCE` · parent [run-status-tokens](../run-status-tokens.md)

## acceptance

A testscript runs a workflow whose fake agent's transcript shows known usage but whose journal has no tokens entry for it; the final printout and run status --json list that session with its counts, the totals include them, and no 'no token use recorded' line appears

## statement

A run's final printout and run status list every session's token use, including a session whose transcript shows tokens but whose journal has no reading yet, and never print 'no token use recorded' for such a session

## Links

### Out

- decided-by [printout-reads-transcript-when-unjournaled](../../../../../decision/printout-reads-transcript-when-unjournaled.md)

### In

- verifies [printout-counts-used-tokens-verified](../../../../../evidence/printout-counts-used-tokens-verified.md)
- verifies [printout-counts-used-tokens-qa-1](../../../../../evidence/printout-counts-used-tokens-qa-1.md)
- verifies [pr-102-run-token-printout](../../../../../evidence/pr-102-run-token-printout.md)

## History

- `01M61GW8RKSBM4BT1B8B3FR5JH` A run's final printout and run status show each session's token use, never 'no token use recorded' for sessions that used tokens.
