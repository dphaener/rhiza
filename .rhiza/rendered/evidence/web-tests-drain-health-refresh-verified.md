# web-tests-drain-health-refresh-verified

evidence · `01M6Z44VXV1Q5G3S32VQ9D096J` · actor [claude-code-61be03b0](../actor/claude-code-61be03b0.md)

## result

PASS

## statement

Verifier claude-code-61be03b0 at 0e2f1c19. Diff: healthCache.get calls an unexported package-level sync.WaitGroup (healthRefreshing) Add(1) before go h.refresh(), and refresh defers Done(). TestWebRootHandlerServesEveryProject registers t.Cleanup(healthRefreshing.Wait) after all its t.TempDir calls. There is no sleep, retry, exported API or synchronous refresh, and WebRootHandler's signature is unchanged. The other internal/api tests reaching WebDaemonSource (web_daemon_test.go) use fake read funcs that never touch disk, so they need no drain. Run from a scratch worktree outside ~/code with RHIZA_ROOT unset, because under ~/code the live root daemon changes .rhiza-root/usage.live.json and trips the state-folder guard: go test -trimpath -count=50 -run TestWebRootHandlerServesEveryProject ./internal/api/ ok; go test -trimpath ./internal/api -run 'TestWebRoot|TestWebDaemon|TestWebHealth' -count=1 ok; go test -trimpath ./internal/web -count=1 ok; the same api subset under -race -count=5 ok.

## Links

### Out

- verifies [web-tests-drain-health-refresh](../capability/rhiza/integrity/self-check/web-tests-drain-health-refresh.md)

## History

- `01M6YKX97C202RDFNV5DGVS033` Fix Flaky test on main: TestWebRootHandlerServesEveryProject (internal/api) fails about 1 run in 10 with 'TempDir RemoveAll cleanup: directory not empty'. Reproduced locally on an unmodified git archive of origin/main with go test -trimpath -count=10 -run TestWebRootHandlerServesEveryProject ./internal/api/ (found while verifying the #398 merge, 2026-11-08). Likely a goroutine or background process (web handler, daemon state write, or git) still writing into t.TempDir() after the test returns. Fix the root cause so -count=50 passes; don't add retries.
