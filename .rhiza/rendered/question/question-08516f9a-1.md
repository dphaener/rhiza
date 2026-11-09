# question-08516f9a-1

question · `01M6YM62WARWNNJ9BB48AY58M3`

## session

08516f9a-2c06-4f2d-96f1-01a1bd3f13ae

## statement

TestWebRootHandlerServesEveryProject flakes because GET / starts a background health refresh (healthCache.get in internal/api/web_daemon.go) that keeps running and touching files under the test's TempDir after the test returns. Where should the fix live?

## Links

### In

- answers [health-refresh-drained-in-tests](../decision/health-refresh-drained-in-tests.md)

## History

- `01M6YKX97C202RDFNV5DGVS033` Fix Flaky test on main: TestWebRootHandlerServesEveryProject (internal/api) fails about 1 run in 10 with 'TempDir RemoveAll cleanup: directory not empty'. Reproduced locally on an unmodified git archive of origin/main with go test -trimpath -count=10 -run TestWebRootHandlerServesEveryProject ./internal/api/ (found while verifying the #398 merge, 2026-11-08). Likely a goroutine or background process (web handler, daemon state write, or git) still writing into t.TempDir() after the test returns. Fix the root cause so -count=50 passes; don't add retries.
