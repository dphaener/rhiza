# inbox-open

capability · `01M61R7QK0MX5D9HZ3V9J9FT24` · parent [inbox](../inbox.md)

## acceptance

A testscript with a fake Rex opens a question item and sees the asking session's tab focused, opens a sign-off item and sees the run's journalled tab focused, gets the run named with its tab gone once the fake Rex no longer has that tab, and gets exit 2 for an unknown item.

## statement

rhiza inbox open <item> focuses the Rex tab of the session that asked a question, or, for a sign-off, a failed run or a run the daemon gave up on, the run's tab as its journal records it, so the owner answers in the agent's own menu; a run whose tab is gone is named with its tab gone, and an unknown item exits 2.

## Links

### Out

- depends-on [inbox-list](inbox-list.md)
- depends-on [run-process-record](../../workflows/workflow-runs/run-process-record.md)
- decided-by [answer-8cd68bc7-16](../../../../decision/answer-8cd68bc7-16.md)
- decided-by [inbox-open-run-tab](../../../../decision/inbox-open-run-tab.md)

### In

- blocks [review-01m61h5z-r3-f4](../../../../question/review-01m61h5z-r3-f4.md)
- verifies [inbox-open-verified](../../../../evidence/inbox-open-verified.md)
- verifies [inbox-open-qa-1](../../../../evidence/inbox-open-qa-1.md)
- verifies [pr-53-the-inbox](../../../../evidence/pr-53-the-inbox.md)

## History

- `01M620GA5ZK6Z854D07NH0R8WX` The inbox: rhiza inbox lists every run's questions, sign-offs, failed runs and runs the daemon gave up on, with a recommendation for each question, opens each in its tab, takes sign-offs, and the daemon sends one Rex notification per new item.
