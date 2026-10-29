# inbox-signoff

capability · `01M61R7R5H6HYW9WR669EGA860` · parent [inbox](../inbox.md)

## acceptance

A testscript with a review run waiting at sign-off approves it through its inbox item and finds the run done with an approved entry, sends another back with one finding and finds a fix cycle started, and gets exit 2 approving a question item; describe-golden is regenerated.

## statement

rhiza inbox approve <item> and rhiza inbox request-changes <item> <file> act on a sign-off item through rhiza run approve and rhiza run request-changes on that item's run, and exit 2 on any other kind of item.

## Links

### Out

- depends-on [inbox-list](inbox-list.md)
- decided-by [answer-8cd68bc7-16](../../../../decision/answer-8cd68bc7-16.md)

### In

- verifies [inbox-signoff-verified](../../../../evidence/inbox-signoff-verified.md)
- verifies [inbox-signoff-qa-1](../../../../evidence/inbox-signoff-qa-1.md)
- verifies [pr-53-the-inbox](../../../../evidence/pr-53-the-inbox.md)

## History

- `01M620GA5ZK6Z854D07NH0R8WX` The inbox: rhiza inbox lists every run's questions, sign-offs, failed runs and runs the daemon gave up on, with a recommendation for each question, opens each in its tab, takes sign-offs, and the daemon sends one Rex notification per new item.
