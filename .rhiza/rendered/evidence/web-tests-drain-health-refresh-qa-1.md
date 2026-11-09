# web-tests-drain-health-refresh-qa-1

evidence · `01M6Z4GGBAYG5D9Z4Y12HEDYGG` · actor [claude-code-8be04309](../actor/claude-code-8be04309.md)

## result

PASS

## statement

QA at branch head 0e7b9a5f. From a scratch worktree outside ~/code: go test -trimpath -count=50 -run TestWebRootHandlerServesEveryProject ./internal/api/ passed 50/50; with the t.Cleanup(healthRefreshing.Wait) line removed, the same run failed 5 of 50 with 'TempDir RemoveAll cleanup: directory not empty', so the drain is what fixes it. go test -trimpath -race -count=10 -run 'TestWebRoot|TestWebDaemon' ./internal/api/ passed. go test -trimpath ./internal/api -run 'TestWebRoot|TestWebDaemon|TestWebHealth' -count=1 and go test -trimpath ./internal/web -count=1 passed. (Run inside ~/code, the -count=50 run passed too, but the enclosing-root guard flagged .rhiza-root/usage.live.json, which a live agent's status line writes; that is unrelated.) Hands-on: built rhiza from the branch and ran rhiza daemon on a scratch root with two projects (alpha, beta), web port 18779, daemon paused. GET / returned 200 in 0.16s listing /p/alpha/ and /p/beta/; /p/alpha/runs and /p/beta/runs returned 200; /p/nope/ returned 404. The first /p/beta/fragments/daemon answered in 0.01s with Health: pending, and 3s later the same fragment showed the landed report (daemon ok, autonomy warn paused, machine ok, disk ok). So the health refresh still runs in the background, off the request path; production behaviour is unchanged.

## Links

### Out

- verifies [web-tests-drain-health-refresh](../capability/rhiza/integrity/self-check/web-tests-drain-health-refresh.md)

## History

- `01M6YKX97C202RDFNV5DGVS033` Fix Flaky test on main: TestWebRootHandlerServesEveryProject (internal/api) fails about 1 run in 10 with 'TempDir RemoveAll cleanup: directory not empty'. Reproduced locally on an unmodified git archive of origin/main with go test -trimpath -count=10 -run TestWebRootHandlerServesEveryProject ./internal/api/ (found while verifying the #398 merge, 2026-11-08). Likely a goroutine or background process (web handler, daemon state write, or git) still writing into t.TempDir() after the test returns. Fix the root cause so -count=50 passes; don't add retries.
