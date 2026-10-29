# inbox-change-set-intent

question · `01M61RKV40P66NSENPEPYSPRC1`

## statement

This change set's intent still reads 'Genesis watches runs and passes their questions to the owner, with a recommendation drawn from slice and why', but answers 4, 7 and 16 made it the whole inbox: sign-offs approved or sent back, failed runs and runs the daemon gave up on, and a notification for each. No rhiza command restates a change set's intent. Should the intent be restated (which needs a way to do it), or should the sign-off and failed-run leaves move to a new change set whose intent covers them?

## Links

### Out

- relates-to [inbox](../capability/rhiza/orchestration/inbox.md) — the inbox capability whose leaves fall outside the change set's intent as written
- blocks change set `01M5SWTKV6SCGKDS1QJ9KN026T` (Genesis watches runs and passes their questions to the owner, with a recommendation drawn from slice and why.)

### In

- answers [inbox-own-change-set](../decision/inbox-own-change-set.md)

## History

- `01M620GA5ZK6Z854D07NH0R8WX` The inbox: rhiza inbox lists every run's questions, sign-offs, failed runs and runs the daemon gave up on, with a recommendation for each question, opens each in its tab, takes sign-offs, and the daemon sends one Rex notification per new item.
