# web-tests-drain-health-refresh

capability · `01M6YMA101E06VZYTZ470T5WER` · parent [self-check](../self-check.md)

## acceptance

healthCache tracks each refresh healthCache.get starts with an unexported WaitGroup or in-flight counter, and TestWebRootHandlerServesEveryProject (plus any other internal/api test serving a page that reaches WebDaemonSource) waits for in-flight refreshes in a t.Cleanup registered after its t.TempDir calls, with no sleep, retry, exported API, synchronous refresh or change to WebRootHandler's signature; go test -trimpath -count=50 -run TestWebRootHandlerServesEveryProject ./internal/api/ passes, and go test -trimpath ./internal/api -run 'TestWebRoot|TestWebDaemon|TestWebHealth' -count=1 and go test -trimpath ./internal/web -count=1 pass

## statement

internal/api's web tests that serve pages wait for every background health refresh the web UI started before their temp folders are removed, so TestWebRootHandlerServesEveryProject no longer fails with 'TempDir RemoveAll cleanup: directory not empty'; the web UI still answers at once from the held report and refreshes in the background

## Links

### Out

- decided-by [health-refresh-drained-in-tests](../../../../decision/health-refresh-drained-in-tests.md)
- derived-from [flaky-web-root-test-plan](../../../../source/flaky-web-root-test-plan.md)

### In

- verifies [web-tests-drain-health-refresh-verified](../../../../evidence/web-tests-drain-health-refresh-verified.md)
- verifies [web-tests-drain-health-refresh-qa-1](../../../../evidence/web-tests-drain-health-refresh-qa-1.md)
- verifies [pr-405-fix-flaky-test-on-main](../../../../evidence/pr-405-fix-flaky-test-on-main.md)

## History

- `01M6YKX97C202RDFNV5DGVS033` Fix Flaky test on main: TestWebRootHandlerServesEveryProject (internal/api) fails about 1 run in 10 with 'TempDir RemoveAll cleanup: directory not empty'. Reproduced locally on an unmodified git archive of origin/main with go test -trimpath -count=10 -run TestWebRootHandlerServesEveryProject ./internal/api/ (found while verifying the #398 merge, 2026-11-08). Likely a goroutine or background process (web handler, daemon state write, or git) still writing into t.TempDir() after the test returns. Fix the root cause so -count=50 passes; don't add retries.
