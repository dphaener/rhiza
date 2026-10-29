# review-01m61h5z-r3-f2

question · `01M61S7K22QYKMWENXMP1VY8EX`

## statement

Plan review round 3 left finding f2 (missing-acceptance on inbox-notify): needs is empty, yet its acceptance covers 'a third [run] the daemon gives up on' and counts notifications per inbox item. A given-up run comes from daemon-resume-once (another change set), and what makes an inbox item, and when it stays or leaves, is defined by inbox-list. Built first, as needs [] allows, the check can't run. It needs inbox-list (which already needs daemon-resume-once). The critic confirmed it: inbox-notify has no needs, yet its statement ('each new inbox item') and its acceptance (four items: a question, a sign-off, a failed run and a run the daemon gave up on, with no repeat while each stays) rely on inbox-list to say what an item is and when it stays, and on daemon-resume-once for the given-up run. Built first, as an empty needs allows, the check has nothing to run against. It needs inbox-list, which already needs daemon-resume-once. How should it be settled?

## Links

### Out

- blocks [inbox-notify](../capability/rhiza/orchestration/inbox/inbox-notify.md)

### In

- answers [inbox-notify-needs-list](../decision/inbox-notify-needs-list.md)

## History

- `01M620GA5ZK6Z854D07NH0R8WX` The inbox: rhiza inbox lists every run's questions, sign-offs, failed runs and runs the daemon gave up on, with a recommendation for each question, opens each in its tab, takes sign-offs, and the daemon sends one Rex notification per new item.
