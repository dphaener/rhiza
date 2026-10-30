# printout-gap-names-reason-qa-1

evidence · `01M65C38B5MBEHPPDQRRJ550JV` · actor [claude-code-a8fce9bb](../actor/claude-code-a8fce9bb.md)

## result

PASS

## statement

QA at 6b63cac with the built binary in a scratch repo whose journal had no tokens entries. Deleting the review session's transcript made run status print '- review, session …: no token use shown, transcript unreadable: no transcript of session …: looked in <home>/.claude/projects/*/<sid>.jsonl', and Totals dropped to the build session's 7/11/13/17 alone. Removing the build session record's exit and its transcript then printed 'no token use shown, still running'. Totals went to all zeros. run status --json gave no_tokens 'still running' and 'transcript unreadable: …' with no tokens objects, and 'no token use recorded' appeared nowhere.

## Links

### Out

- verifies [printout-gap-names-reason](../capability/rhiza/runtime/session-journal/run-status-tokens/printout-gap-names-reason.md)

## History

- `01M61GW8RKSBM4BT1B8B3FR5JH` A run's final printout and run status show each session's token use, never 'no token use recorded' for sessions that used tokens.
