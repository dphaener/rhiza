# inbox-clears-stale-verified

evidence · `01M63492JC1VQ8YDJHXMC5YMJF` · actor [claude-code-355f3b69](../actor/claude-code-355f3b69.md)

## result

PASS

## statement

At c3aa860: go test ./cmd/rhiza -run TestScripts/inbox-stale passes. The script builds fixture change sets and journals: a failed run and a gave-up run each on landed, withdrawn, rejected and superseded change sets (a dumped rhiza change ls --json showed all eight in those states, superseded included), one of each on an accepted change set, a failed run on unknown change set CSUNKNOWN, and a question and a sign-off on a landed change set. rhiza inbox and inbox --json list exactly F-ACCEPTED, G-ACCEPTED, F-UNKNOWN, Q-LANDED and S-LANDED (5 kinds in --json), with none of the eight stale runs. Mutation probe in a scratch worktree: removing the gone() check on ended failed runs makes the script fail (F-WITHDRAWN listed), and removing it on the gave-up branch alone also makes it fail. The diff applies goneChanges (State.Terminal: landed, rejected, withdrawn, superseded; an unknown change or unreadable log counts as not gone) to inbox, inbox open and daemon notify. make fmt-check, vet and go test ./... pass; make self-check fails only on stale rendered actor files for this run's sessions, not on this leaf's code.

## Links

### Out

- verifies [inbox-clears-stale](../capability/rhiza/orchestration/inbox/inbox-clears-stale.md)

## History

- `01M62MKZFRPZ2C7QTYAF815XAG` The inbox leaves out failed runs whose change set has since landed, been withdrawn or rejected, and rhiza run abandon also takes a failed run, so every failed item can be cleared; no item's next command is one that can only fail again.
