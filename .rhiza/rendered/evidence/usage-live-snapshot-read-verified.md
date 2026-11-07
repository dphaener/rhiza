# usage-live-snapshot-read-verified

evidence · `01M6TCFMXNPRPVY23N05V8TCM5` · actor [claude-code-8dbeee5d](../actor/claude-code-8dbeee5d.md)

## result

PASS

## statement

At 82b33999 ran go test -trimpath -count=1 ./internal/adapter/... (ok) and TestLiveSnapshots -v (PASS); go vet and gofmt -l clean. The test writes a fake state folder's usage.live.json and reads it back through claudecode.Adapter.LiveSnapshots as one snapshot per kind (session, weekly, sorted) carrying used_percentage, resets_at and captured_at; an absent file and a corrupt file each return nil with no error (the method has no error return, so a poll cannot be stopped). Read the diff: adapter.LiveUsageReader is an optional interface with a compile-time assertion for Adapter, and the reader's file name and JSON tags match daemon.LiveFile and daemon.RecordLive's array output.

## Links

### Out

- verifies [usage-live-snapshot-read](../capability/rhiza/orchestration/usage-governor/usage-live-snapshot-read.md)

## History

- `01M6T2DCTWTJ0GW5AQMA2FZJF8` Fix the usage governor's limits are static (rhiza.daemon.usage.session/weekly or a limit learned from one hit), so they go stale whenever the account's plan or Anthropic's limits change (the owner moved to Max 20x and every limit was wrong); they must be dynamic: read the account's live usage from what Claude Code reports (its status line input carries rate_limits with used_percentage and resets_at per window, e.g. seven_day 49%) or the equivalent API, derive each window's limit and reset from it at every poll, and fall back to the settings only when no live figure is available
