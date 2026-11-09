# flaky-web-root-test-plan

source · `01M6YM9KRGT5QZ3CWRHBWC5Y6T`

## quotes

```json
[]
```

## raw_sha256

7b7590161d01f055748addd3c4c70ba2e68ce38ad0b83350fef6b78e706a211a

## session

08516f9a-2c06-4f2d-96f1-01a1bd3f13ae

## summary

Autonomous plan run, no owner present. Reproduced the flake on this branch: go test -trimpath -count=20 -run TestWebRootHandlerServesEveryProject ./internal/api/ failed 2 of 20 with 'TempDir RemoveAll cleanup: unlinkat .../001: directory not empty'. Root cause: GET / and /p/<name>/runs reach WebDaemonSource, whose healthCache.get (internal/api/web_daemon.go) starts a background runHealth goroutine; it is still running and briefly creating files under the projects when the test returns and t.TempDir's RemoveAll runs. Files listed in a cleanup before RemoveAll and again 2 s later match, so the writer makes short-lived files only. The planner asked where the fix should live (product-side awaitable refresh or test only); the arbiter decided health-refresh-drained-in-tests: track each background refresh inside internal/api (unexported WaitGroup or counter) and have the page-serving web tests wait for in-flight refreshes in a t.Cleanup registered after their t.TempDir calls; production behaviour, exported API and WebRootHandler's signature unchanged; no sleeps, no retries, no synchronous refresh. Side note: the same run also tripped the enclosing-root guard on .rhiza-root/usage.live.json, which a live agent's status line on the machine writes; it is unrelated to this test.

## title

Plan of the flaky TestWebRootHandlerServesEveryProject fix (autonomous)

## Links

### In

- derived-from [web-tests-drain-health-refresh](../capability/rhiza/integrity/self-check/web-tests-drain-health-refresh.md)

## History

- `01M6YKX97C202RDFNV5DGVS033` Fix Flaky test on main: TestWebRootHandlerServesEveryProject (internal/api) fails about 1 run in 10 with 'TempDir RemoveAll cleanup: directory not empty'. Reproduced locally on an unmodified git archive of origin/main with go test -trimpath -count=10 -run TestWebRootHandlerServesEveryProject ./internal/api/ (found while verifying the #398 merge, 2026-11-08). Likely a goroutine or background process (web handler, daemon state write, or git) still writing into t.TempDir() after the test returns. Fix the root cause so -count=50 passes; don't add retries.
