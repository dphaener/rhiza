# printout-reads-transcript-when-unjournaled

decision · `01M650632VFV6TVF96ACAVR0D9`

## reasoning

The change set intent is that the printout and run status never say 'no token use recorded' for a session that used tokens. Readings are journaled only at done, step end or exit, so a live session, or one whose transcript could not be read at that time, currently gets that line (internal/api/run.go:2297-2298). Reading the transcript at print time closes the live-session gap. Writing the reading back to the log from a print or status call would turn a read path into a write, and the log is append-only and written through change sets, so the read result stays unjournaled. An unreadable transcript still cannot yield counts, so that case gets an honest 'unreadable' line rather than a claim of no use.

## statement

Yes. When a session has no journaled token reading, the final printout and run status read its transcript at print time and show the counts they find, without journaling the result. Only when the transcript still cannot be read do they say so, instead of 'no token use recorded'.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-4ba1eee9-1](../question/question-4ba1eee9-1.md)
- answers [token-gap-read-at-print](../question/token-gap-read-at-print.md)
- derived-from [run-token-printout-autonomous-plan](../source/run-token-printout-autonomous-plan.md)

### In

- decided-by [printout-counts-used-tokens](../capability/rhiza/runtime/session-journal/run-status-tokens/printout-counts-used-tokens.md)
- decided-by [printout-gap-names-reason](../capability/rhiza/runtime/session-journal/run-status-tokens/printout-gap-names-reason.md)

## History

- `01M61GW8RKSBM4BT1B8B3FR5JH` A run's final printout and run status show each session's token use, never 'no token use recorded' for sessions that used tokens.
