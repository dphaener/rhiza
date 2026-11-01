# inbox-overlap-items-verified

evidence · `01M6A9NY2QQSRP9B86KPRKNAA4` · actor [claude-code-6bc8960c](../actor/claude-code-6bc8960c.md)

## result

PASS

## statement

At commit 77622d2b I re-ran the acceptance. go test ./cmd/rhiza -run 'TestScripts/(inbox-list|describe-golden)$' passes. inbox-list.txtar has a new case: two proposed change sets set thing's statement. rhiza inbox lists '<id> overlap change sets A and B', 'target: thing' and 'both write field:statement'. --json has kind overlap, target and other, and no run. inbox open on the item exits 2 with an overlap message. describe-golden.txtar is regenerated for the new kind and the other/target fields. TestInboxNoticeOverlap passes: the daemon notice reads 'an overlap with change set X waits on you', and daemonNotify appends overlap items. TestViewGroupsOverlaps passes: the inbox app adds an Overlaps group and shows 'Seven · with C8 · thing' on the title line. That fixes attempt 1's finding f1. go test ./... -count=1, go vet and gofmt are clean. The diff in inbox.go, daemon.go and inboxapp.go matches the statement. make self-check fails only on stale .rhiza/rendered from earlier verifiers' actor records, which is run bookkeeping and not this leaf's code.

## Links

### Out

- verifies [inbox-overlap-items](../capability/rhiza/orchestration/inbox/inbox-overlap-items.md)

## History

- `01M5SWTNWB9GEJ5AA8R7K867KB` Genesis detects conflicting decisions across change sets in flight and surfaces them to the owner.
