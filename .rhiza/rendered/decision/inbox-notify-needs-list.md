# inbox-notify-needs-list

decision · `01M620KNHFWVNGE2J4CY4JWJD4`

## statement

inbox-notify depends on inbox-list, which defines what an inbox item is and when it stays, so its check runs against real items. The check stays as written.

## wording

the orchestrator's option, picked by the owner: "Add the needs (Recommended)" (daemon-status needs daemon-picks and daemon-run-tabs; inbox-notify needs inbox-list; the checks stay as strong as written)

## Links

### Out

- answers [review-01m61h5z-r3-f2](../question/review-01m61h5z-r3-f2.md)

### In

- decided-by [inbox-notify](../capability/rhiza/orchestration/inbox/inbox-notify.md)

## History

- `01M620GA5ZK6Z854D07NH0R8WX` The inbox: rhiza inbox lists every run's questions, sign-offs, failed runs and runs the daemon gave up on, with a recommendation for each question, opens each in its tab, takes sign-offs, and the daemon sends one Rex notification per new item.
