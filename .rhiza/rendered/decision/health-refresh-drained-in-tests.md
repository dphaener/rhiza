# health-refresh-drained-in-tests

decision · `01M6YM87Q7DQJ2JRS7DVB5DJ15`

## reasoning

The question names the cause: healthCache.get (internal/api/web_daemon.go) fires `go h.refresh()`, which runs runHealth over the project under t.TempDir and can still be writing when the test returns, so TempDir's RemoveAll finds a non-empty dir. The background refresh is deliberate product behaviour (the comments on healthCache, WebHealth and fingerprint: the page answers from memory and the /events stream moves when a report lands), so the fix must not change it. WebRootHandler and WebProjectHandler build their WebHealth internally, so a test cannot reach the cache directly; a package-internal tracker of in-flight refreshes, drained by the test in a LIFO cleanup that runs before TempDir's, removes the race at its root and satisfies the change set's intent ('fix the root cause so -count=50 passes; don't add retries'). It stays inside internal/api, touches no request/response types (describe-golden unaffected), and also covers TestOverviewReadsDaemonFromRepo-style tests that use NewWebHealth over a TempDir.

## statement

Fix the flake in internal/api, at the health cache: track each background health refresh healthCache.get starts (an unexported WaitGroup or in-flight counter, package-internal), and have the web tests that serve pages wait for in-flight refreshes in a t.Cleanup registered after their t.TempDir calls, so the refresh finishes before the TempDir is removed. Production behaviour stays as is: GET / and the Overview still return the held report at once and refresh in the background, with no new exported API, no sleeps and no retries. Don't make the refresh synchronous and don't change WebRootHandler's signature just for the test.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-08516f9a-1](../question/question-08516f9a-1.md)

### In

- decided-by [web-tests-drain-health-refresh](../capability/rhiza/integrity/self-check/web-tests-drain-health-refresh.md)

## History

- `01M6YKX97C202RDFNV5DGVS033` Fix Flaky test on main: TestWebRootHandlerServesEveryProject (internal/api) fails about 1 run in 10 with 'TempDir RemoveAll cleanup: directory not empty'. Reproduced locally on an unmodified git archive of origin/main with go test -trimpath -count=10 -run TestWebRootHandlerServesEveryProject ./internal/api/ (found while verifying the #398 merge, 2026-11-08). Likely a goroutine or background process (web handler, daemon state write, or git) still writing into t.TempDir() after the test returns. Fix the root cause so -count=50 passes; don't add retries.
