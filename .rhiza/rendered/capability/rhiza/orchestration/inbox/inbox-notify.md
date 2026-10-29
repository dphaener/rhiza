# inbox-notify

capability · `01M61R7SANZRP2X52437ENZGTH` · parent [inbox](../inbox.md)

## acceptance

With a fake Rex, rhiza daemon --once over a run that blocks on a question, one waiting at sign-off, one that failed and one the daemon gives up on sends exactly four notifications, one per item, and a second rhiza daemon --once with the same items sends none.

## statement

On each poll, rhiza daemon sends one Rex notification for each inbox item it has not notified yet, naming the change set and what waits (a question, a sign-off, a failed run, a run the daemon gave up on or an interrupted run that was not resumed), and none again while the item stays; a run keeps its own blocked-run notice.

## Links

### Out

- depends-on [inbox-list](inbox-list.md)
- depends-on [daemon-picks](../daemon/daemon-picks.md)
- decided-by [answer-8cd68bc7-18](../../../../decision/answer-8cd68bc7-18.md)
- decided-by [daemon-sends-inbox-notices](../../../../decision/daemon-sends-inbox-notices.md)
- decided-by [inbox-notify-needs-list](../../../../decision/inbox-notify-needs-list.md)

### In

- blocks [review-01m61h5z-r3-f2](../../../../question/review-01m61h5z-r3-f2.md)
- blocks [review-01m61h5z-r3-f3](../../../../question/review-01m61h5z-r3-f3.md)
- verifies [inbox-notify-verified](../../../../evidence/inbox-notify-verified.md)
- verifies [inbox-notify-qa-1](../../../../evidence/inbox-notify-qa-1.md)
- verifies [pr-53-the-inbox](../../../../evidence/pr-53-the-inbox.md)
- verifies [pr-64-resume-only-in-flight-runs](../../../../evidence/pr-64-resume-only-in-flight-runs.md)

## History

- `01M620GA5ZK6Z854D07NH0R8WX` The inbox: rhiza inbox lists every run's questions, sign-offs, failed runs and runs the daemon gave up on, with a recommendation for each question, opens each in its tab, takes sign-offs, and the daemon sends one Rex notification per new item.
- `01M62MKZ6SA46DYC1K2R0W254S` The daemon resumes only interrupted runs whose change set is still accepted: a run left interrupted on a change set that has since landed, been withdrawn or rejected is never resumed.
