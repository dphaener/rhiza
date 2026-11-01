# inbox-disk-low-verified-2

evidence · `01M6B4ZBGMFZ6210J553XZ70WD` · actor [claude-code-75141a6c](../actor/claude-code-75141a6c.md)

## result

PASS

## statement

At 20922097: go test ./cmd/rhiza -run 'TestScripts/(inbox-disk-low|daemon-notify-disk-low|daemon-picks)$' passes; inbox-disk-low.txtar, with a daemon.state.json fixture recording a disk-low poll, shows exactly one "kind": "disk-low" item in rhiza inbox --json ('1.0 GiB free, under the minimum of 10.0 GiB') and items [] once the state from a later poll has no disk_low. End to end with a scratch build in a fresh repo (RHIZA_TEST_NOW set): rhiza daemon --once with RHIZA_TEST_FREE_BYTES=1 GiB then rhiza inbox --json listed one disk-low item naming the free space; a second --once poll at 20 GiB left items []. Diff read: diskLowItems is shared by the inbox and daemonNotify, which now uses this poll's in-memory state. make fmt-check vet test pass; make self-check fails only on stale .rhiza/rendered actor files from the review run's chore commits, not this leaf.

## Links

### Out

- verifies [inbox-disk-low](../capability/rhiza/orchestration/inbox/inbox-list/inbox-disk-low.md)

## History

- `01M6ATAS02XDRJ1MSQ4THGCNPA` The daemon holds new starts and resumes while disk space is low, and daemon status and the inbox say so
