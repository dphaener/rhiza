# review-01m61h5z-r3-f4

question · `01M61S7KVVQ9F9Z9CF3RN3S568`

## statement

Plan review round 3 left finding f4 (contradiction on inbox-open): It focuses 'the Rex tab ... of the run for a sign-off, a failed run or a run the daemon gave up on', but no record holds a run's tab. The journal keeps tabs only for blocked sessions (Block.TabSession/TabWindow/TabBlock in internal/run/journal.go:200), run-process-record journals only host and pid, and daemon-run-tabs doesn't journal the tab it opens. A leaf has to record the run's tab (or inbox-open has to name another target), or this part of the acceptance has nothing to focus. The critic confirmed it: The journal keeps a tab only in Block (TabSession/TabWindow/TabBlock, internal/run/journal.go:198-206), for a session blocked on the owner. Started has no tab (journal.go:175-190), run-process-record journals only host and pid, and daemon-run-tabs doesn't say it journals the tab it opens. inbox-open's acceptance focuses 'the run's tab' for a sign-off, and its statement does the same for failed and given-up runs, but nothing records that tab. A leaf has to journal the run's tab, or inbox-open has to name another target. How should it be settled?

## Links

### Out

- blocks [inbox-open](../capability/rhiza/orchestration/inbox/inbox-open.md)

### In

- answers [inbox-open-run-tab](../decision/inbox-open-run-tab.md)

## History

- `01M620GA5ZK6Z854D07NH0R8WX` The inbox: rhiza inbox lists every run's questions, sign-offs, failed runs and runs the daemon gave up on, with a recommendation for each question, opens each in its tab, takes sign-offs, and the daemon sends one Rex notification per new item.
