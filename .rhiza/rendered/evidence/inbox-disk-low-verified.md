# inbox-disk-low-verified

evidence · `01M6B345X6WXTS2W8FEH5DTRX8` · actor [claude-code-a951bbd8](../actor/claude-code-a951bbd8.md)

## result

PASS

## statement

At e8e4de61 ran go test ./cmd/rhiza -run 'TestScripts/(inbox-disk-low|daemon-picks)$' -count=1: both PASS. inbox-disk-low.txtar puts a daemon state with disk_low (1 GiB free, 10 GiB min) and rhiza inbox --json shows exactly one "kind": "disk-low" item naming '1.0 GiB free, under the minimum of 10.0 GiB'; after a later state without disk_low it shows none (items []). daemon-picks.txtar shows the same through real daemon polls (RHIZA_TEST_FREE_BYTES low, then 20 GiB clears it). go test ./internal/tui/inboxapp -run DiskLow PASS (app groups it under 'Disk space'). make test, fmt-check, vet pass; make self-check passes with RHIZA_ROOT unset. Read the diff: daemon sets/clears State.DiskLow per successful poll and keeps Since across polls; inbox open refuses the item.

## Links

### Out

- verifies [inbox-disk-low](../capability/rhiza/orchestration/inbox/inbox-list/inbox-disk-low.md)

## History

- `01M6ATAS02XDRJ1MSQ4THGCNPA` The daemon holds new starts and resumes while disk space is low, and daemon status and the inbox say so
