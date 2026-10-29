# daemon-sends-inbox-notices

decision · `01M620KP4AZ2MJDFP7565KD9Z9`

## statement

The daemon sends inbox notifications: on each poll it finds inbox items it has not notified yet and sends one Rex notification per item, none again while the item stays. A run keeps its own blocked-run notice.

## wording

the orchestrator's option, picked by the owner: "The daemon (Recommended)" (on each poll the daemon notices new inbox items and sends one Rex notification per item, with no repeat while it stays; runs keep their own blocked notice; inbox-notify then needs the daemon)

## Links

### Out

- answers [review-01m61h5z-r3-f3](../question/review-01m61h5z-r3-f3.md)

### In

- decided-by [inbox-notify](../capability/rhiza/orchestration/inbox/inbox-notify.md)

## History

- `01M620GA5ZK6Z854D07NH0R8WX` The inbox: rhiza inbox lists every run's questions, sign-offs, failed runs and runs the daemon gave up on, with a recommendation for each question, opens each in its tab, takes sign-offs, and the daemon sends one Rex notification per new item.
