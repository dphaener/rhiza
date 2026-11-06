# web-run-log-tail-verified-2

evidence · `01M6Q9VVY2GYFXZCK67KTVVQ8Z` · actor [claude-code-db656a66](../actor/claude-code-db656a66.md)

## result

PASS

## statement

Verified at d5ceb505 (attempt 2). Read the diff: run.status is now ReadOnly (internal/api/run.go:58) and listed in TestReadOnlyMarker, so the run page no longer 404s through the daemon's read-only API; TestRunLogTailFollowsDriverLog now also GETs /runs/run1 through api.ReadOnlyHandler and asserts the page shows the tail (rz-term, 'first line') and the run-log fragment. Ran go test -trimpath -count=1 ./internal/web/ : ok, with TestRunPageLogTail (fixture tail rendered in order with cursor, live on log events), TestRunPageNoLog (No log state on page and fragment), TestRunPageThroughReadOnlyAPI and TestRunLogTailFollowsDriverLog (lines appended to the driver log raise the stream's log event and show in the tail fragment; no-log before) all PASS. go test ./internal/api -run 'TestReadOnlyMarker|Web|Tail|RunLog' PASS; go test ./cmd/rhiza -run 'TestScripts/(run-tail|describe-golden|web)' PASS (run tail --limit 2 keeps only the last lines; no log exits 2). The full internal/api package timed out in TestReadBudget/TestInboxReadsChangeSetsOnce under load avg 12-20, unrelated to this leaf.

## Links

### Out

- verifies [web-run-log-tail](../capability/rhiza/web-monitor/web-run-log-tail.md)

## History

- `01M6PFH1X52KNK0C5ZCDYEEF64` Web Runs pages: runs live first then recent, and a run's detail with its steps, sessions, review summary and a live log tail
