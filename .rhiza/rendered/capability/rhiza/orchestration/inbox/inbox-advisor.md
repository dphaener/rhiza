# inbox-advisor

capability · `01M61R7RR0R434KJZZJ75FE90Z` · parent [inbox](../inbox.md)

## acceptance

A testscript with a fake adapter whose agent blocks on a question finds an advisor launched with that question, its recommendation journalled and shown by rhiza inbox, no message sent to the asking session and the run still blocked; with a failing fake advisor the item shows no recommendation and the run is still blocked.

## statement

When a normal run blocks on a question, the engine launches a short-lived advisor session (the arbiter's guide in recommend-only mode) that reads slice and why for the question and journals a recommendation with its reasons, never answering it; the run does not wait for the advisor, the inbox shows the recommendation beside the question, and a failed advisor leaves the question without one.

## Links

### Out

- depends-on [inbox-list](inbox-list.md)
- decided-by [answer-8cd68bc7-17](../../../../decision/answer-8cd68bc7-17.md)

### In

- verifies [inbox-advisor-verified](../../../../evidence/inbox-advisor-verified.md)
- verifies [inbox-advisor-qa-1](../../../../evidence/inbox-advisor-qa-1.md)
- verifies [pr-53-the-inbox](../../../../evidence/pr-53-the-inbox.md)

## History

- `01M620GA5ZK6Z854D07NH0R8WX` The inbox: rhiza inbox lists every run's questions, sign-offs, failed runs and runs the daemon gave up on, with a recommendation for each question, opens each in its tab, takes sign-offs, and the daemon sends one Rex notification per new item.
