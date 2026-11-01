# inbox-stalled-items-verified-2

evidence · `01M6AVGM78DX3R860E599NDRE5` · actor [claude-code-aa98a951](../actor/claude-code-aa98a951.md)

## result

PASS

## statement

At eb20db85 (attempt 2): go test ./cmd/rhiza -run TestScripts/inbox-stalled passes (live run, hand-set working hook report at 00:00Z: no stalled item at 00:19Z, one kind stalled at 00:21Z with step build, at, and next rhiza session message; blocked and idle records give none). go test ./internal/tui/inboxapp -run TestViewGroupsStalled passes (Stalled group shown). The fix diff replaces runStalledNow with runStalledAll in inboxStalledOf so the inbox adds one item per stalled session of the current step; TestRunStalledNow now checks parallel checks plus fixer all listed and none for a stopped run; this addresses QA finding f2. go vet and gofmt clean; make check fails only on stale .rhiza/rendered files from the run's own model writes.

## Links

### Out

- verifies [inbox-stalled-items](../capability/rhiza/orchestration/inbox/inbox-list/inbox-stalled-items.md)

## History

- `01M6A0K4K58YMW04FA38373Y8P` A session with no activity for 20 minutes shows as stalled in run status and the inbox
