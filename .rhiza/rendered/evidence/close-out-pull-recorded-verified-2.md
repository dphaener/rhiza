# close-out-pull-recorded-verified-2

evidence · `01M6DGJZQ5V26AWAA5X2GM46KF` · actor [claude-code-69ead661](../actor/claude-code-69ead661.md)

## result

PASS

## statement

At aeed9da6 ran go test ./internal/api -run TestShipCloseOutRecordsStaleMain -count=1: PASS. The test runs ship.close-out with a dirty main checkout (status shows ' M README.md'): MainUpdated is false, main-stale.json records 'has uncommitted changes', inboxMainStale lists one main-stale item, and healthMain is warn with summary 'the last close-out could not pull' and detail naming the change set and failure; a repeat keeps the item ID; with main clean a later close-out pulls, clears the record, the inbox item and the health warning. Read the diff: internal/api/health.go:253 reads daemon.ReadMainStale, internal/api/inbox.go:533 lists it. make check: fmt, vet and all go tests pass; self-check fails only on this verifier session's own unrendered actor file, not on the commit.

## Links

### Out

- verifies [close-out-pull-recorded](../capability/rhiza/orchestration/daemon/daemon-status/health-command/close-out-pull-recorded.md)

## History

- `01M6CRZSRFQ94MAJ8CRTEE44P2` rhiza health and a health app show the whole system's state on one screen
