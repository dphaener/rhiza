# inbox-notify-verified

evidence · `01M62K23SVESZWXBR39W4ZDMH1` · actor [claude-code-9228c9a7](../actor/claude-code-9228c9a7.md)

## result

PASS

## statement

At c0a0879 (builder commit 78720ab), ran go test ./cmd/rhiza -run 'TestScripts/daemon-notify$' -count=1 -v. With a fake rex on PATH and REX_BLOCK set, the first rhiza daemon --once over four fixture runs (one blocked on a question, one waiting at sign-off, one finished failed, one the daemon gave up on after a dead resumed process) wrote exactly four OSC 7501 state=blocked reports to stderr. Each was preceded by a working report, and their decoded msgs were 'CSQUESTION: a question waits on you: Which store?', 'CSSIGNOFF: a sign-off waits on you', 'CSFAILED: a run failed' and 'CSGAVEUP: the daemon gave up on a run'. A second --once with the same items wrote none. After the sign-off was approved and a new failed run was added, a third --once wrote exactly one ('CSFAILEDTWO: a run failed'). Read the diff: daemonNotify in internal/api/daemon.go runs on every poll, and the notified IDs persist in daemon.State.Notified, pruned to the items that still wait. Run notices are untouched. make fmt-check vet test passed with every package ok. In a clean clone at c0a0879, rhiza check printed ok. make self-check in the run checkout failed only on stale renders of actor records that the run added to the shared .git/rhiza log after the commit, including this verifier's own actor record.

## Links

### Out

- verifies [inbox-notify](../capability/rhiza/orchestration/inbox/inbox-notify.md)

## History

- `01M620GA5ZK6Z854D07NH0R8WX` The inbox: rhiza inbox lists every run's questions, sign-offs, failed runs and runs the daemon gave up on, with a recommendation for each question, opens each in its tab, takes sign-offs, and the daemon sends one Rex notification per new item.
