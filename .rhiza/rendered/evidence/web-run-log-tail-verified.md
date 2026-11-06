# web-run-log-tail-verified

evidence · `01M6Q71NVZQ2A3H8AAHB87EWTZ` · actor [claude-code-b98edd62](../actor/claude-code-b98edd62.md)

## result

PASS

## statement

At b30eaa9b, as verifier claude-code-b98edd62: go test -trimpath -count=1 ./internal/web -run 'TestRunLogTailFollowsDriverLog|TestRunPageLogTail|TestRunPageNoLog' -v all PASS, and the whole internal/web package passes. TestRunPageLogTail renders a fixture run's page with its last lines (HTML-escaped) in the rz-term block with the live cursor, live on log events; TestRunLogTailFollowsDriverLog (real read-only API over fixtureCtx) asserts the no-log state, appends lines to the run's driver log, waits for the stream's 'event: log' with the run id, and reads them through /fragments/run-log, then again for a further appended line; TestRunPageNoLog asserts 'No log' on page and fragment. Also passed: internal/api TestReadOnlyMarker (run.tail is ReadOnly) and cmd/rhiza TestScripts/run-tail and describe-golden. Read the diff: run.tail reads a bounded window of the log, 404 maps to the no-log Empty state. Note: a full internal/api run hit the 600s go test timeout on a machine at load ~20; the scoped api test passes.

## Links

### Out

- verifies [web-run-log-tail](../capability/rhiza/web-monitor/web-run-log-tail.md)

## History

- `01M6PFH1X52KNK0C5ZCDYEEF64` Web Runs pages: runs live first then recent, and a run's detail with its steps, sessions, review summary and a live log tail
