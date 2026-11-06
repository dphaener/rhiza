# web-overview-fast-verified

evidence · `01M6QVXA913J3W0MVRJ916PWTT` · actor [claude-code-ddef990a](../actor/claude-code-ddef990a.md)

## result

PASS

## statement

Verified at ddf1ad2d. Read the diffs of 268eb67f and ddf1ad2d: healthCache in internal/api/web_daemon.go serves the held report and starts a background refresh only when none is running and the report is missing or older than 30 s; running, paused and slots still come from runDaemonStatus on every call; the template shows a pending pill until the first report. TestWebDaemonHealthCachedInBackground (internal/api/web_daemon_test.go) blocks the health read until released and asserts / and /fragments/daemon each answer under 400 ms, pending (not failed) pills first, 20 concurrent requests start exactly one refresh (now waiting for it to start, fixing attempt 1's flake), the report shows once it lands, and a report a minute old is served to 5 requests while exactly one new refresh starts and its result then shows. Ran the attempt-1 repro `go test -trimpath -c -race -o api.test ./internal/api && api.test -test.run TestWebDaemon -test.count=300 -test.cpu 1,2,4,8` twice: 1200 PASS, 0 FAIL each. `go test -trimpath -race -count=1 ./internal/web/...` passed; `go vet ./internal/api ./internal/web` clean.

## Links

### Out

- verifies [web-overview-fast](../capability/rhiza/web-monitor/web-overview-fast.md)

## History

- `01M6QSJHF5WH4SZKCTTYD0W0SS` Fix The web UI's Overview page takes about 5 s to load (owner: "it's incredibly slow"). Measured on the VM: / 5.2 s and /fragments/daemon 4.7 s, while the inbox, live-runs and milestones sections and /runs and /roadmap answer in 0.01-0.11 s. The daemon panel's source (WebDaemonSource, internal/api/web_daemon.go) calls runHealth on every render, and health runs the full model check (rhiza health alone takes 5.0 s here). The page and the daemon section's live refreshes must answer fast (well under 0.5 s) without running the check on each request, e.g. by reusing the daemon's last health result instead of recomputing it.
