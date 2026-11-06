# web-overview-fast-verified-2

evidence · `01M6QWXBAVCTQXKKX4JV2VQRH5` · actor [claude-code-9192c011](../actor/claude-code-9192c011.md)

## result

PASS

## statement

At 22615c6c (attempt 2) I read the diff of internal/api/web_daemon.go, web_events.go and daemon.go and ran go test -race -count=5 -run TestWebDaemonHealth ./internal/api plus go test -count=1 ./internal/api/... ./internal/web/... and go vet: all pass. TestWebDaemonHealthCachedInBackground serves / and /fragments/daemon with a health read that blocks until released and asserts each answers under 400 ms, the first answer shows the Checking… pending pill and no report pills, 20 concurrent requests start exactly one refresh, the landed report replaces pending, and a report a minute old is served unchanged while exactly one new refresh starts and then shows. The attempt-1 confirmed quality finding is fixed: daemon.go builds one NewWebHealth shared by WebDaemonSource and WebEventSources, the daemon event fingerprint carries the held report, and TestWebDaemonHealthLandingSendsDaemonEvent shows no daemon event before the report lands, one when it lands (the re-fetched panel shows model: ok, not pending), and none for an identical refresh. Running, paused and slots are still read from the state file per request.

## Links

### Out

- verifies [web-overview-fast](../capability/rhiza/web-monitor/web-overview-fast.md)

## History

- `01M6QSJHF5WH4SZKCTTYD0W0SS` Fix The web UI's Overview page takes about 5 s to load (owner: "it's incredibly slow"). Measured on the VM: / 5.2 s and /fragments/daemon 4.7 s, while the inbox, live-runs and milestones sections and /runs and /roadmap answer in 0.01-0.11 s. The daemon panel's source (WebDaemonSource, internal/api/web_daemon.go) calls runHealth on every render, and health runs the full model check (rhiza health alone takes 5.0 s here). The page and the daemon section's live refreshes must answer fast (well under 0.5 s) without running the check on each request, e.g. by reusing the daemon's last health result instead of recomputing it.
