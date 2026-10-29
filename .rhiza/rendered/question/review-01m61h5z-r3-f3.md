# review-01m61h5z-r3-f3

question · `01M61S7KEZ0PGK49670XG9GWKD`

## statement

Plan review round 3 left finding f3 (contradiction on inbox-notify): It sends each notification 'through the blocked-run notice path', but that path is in the run process: Engine.Notify -> runRex.notice in internal/api/run.go:1482, called only from block/unblock in internal/run/engine.go:862,879. A run the daemon gave up on is interrupted, so no run process is left to notify. That item, and any item the daemon or the inbox spots after the fact, needs a sender outside the run, which neither the leaf nor the daemon leaves name. The critic confirmed it: The blocked-run notice path is in the run process. reportToRex sets Engine.Notify to runRex.notice (internal/api/run.go:1461-1500), which writes an OSC 7501 record to the run's own terminal, and the engine calls it only on block and unblock (internal/run/engine.go:885-889). A run the daemon gave up on is interrupted, so its process is gone and nothing can send that item's notice. No leaf gives the daemon or the inbox its own sender, so inbox-notify's acceptance can't be met for that item as written. How should it be settled?

## Links

### Out

- blocks [inbox-notify](../capability/rhiza/orchestration/inbox/inbox-notify.md)

### In

- answers [daemon-sends-inbox-notices](../decision/daemon-sends-inbox-notices.md)

## History

- `01M620GA5ZK6Z854D07NH0R8WX` The inbox: rhiza inbox lists every run's questions, sign-offs, failed runs and runs the daemon gave up on, with a recommendation for each question, opens each in its tab, takes sign-offs, and the daemon sends one Rex notification per new item.
