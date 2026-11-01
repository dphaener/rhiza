# inbox-stalled-items-verified

evidence · `01M6ATJ7N8G0PF3H67PNAJ49FV` · actor [claude-code-c79d5731](../actor/claude-code-c79d5731.md)

## result

PASS

## statement

At ffd9ee2d, go test ./cmd/rhiza -run 'TestScripts/(inbox-stalled|describe-golden)$' -count=1 passes: inbox-stalled.txtar starts a live run whose agent reports working, sets its hook report to 2030-01-01T00:00:00Z, and rhiza inbox --json has no stalled item at 00:19 and one of kind stalled (step build, at 00:00:00Z, next rhiza session message <session> <text>) at 00:21, with the human output naming run, change set and step; turning the report blocked, then idle, leaves no stalled item. go test ./internal/tui/inboxapp -run TestViewGroupsStalled passes: the app shows a Stalled group with the item. go test ./... passes. Diff read: inboxItems takes a stalled func built from runStalledNow over the session store at ctx.Now, applied only to live runs. make check's self-check fails only on rendered actor files of this run's own sessions (unlanded run change set), the same artifact the earlier leaf carries, not on this leaf's code.

## Links

### Out

- verifies [inbox-stalled-items](../capability/rhiza/orchestration/inbox/inbox-list/inbox-stalled-items.md)

## History

- `01M6A0K4K58YMW04FA38373Y8P` A session with no activity for 20 minutes shows as stalled in run status and the inbox
