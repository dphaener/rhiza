# web-overview-fast-qa-1

evidence · `01M6QW6MEN1FKVSVMKYTY9DDV9` · actor [claude-code-e0e47890](../actor/claude-code-e0e47890.md)

## result

PASS

## statement

QA at e8d75459 with the built binary (go build -trimpath ./cmd/rhiza), in a scratch clone of this repo whose own `rhiza health` takes 11.9 s. Ran `rhiza daemon --max-runs 1 --permission-mode plan` with rhiza.daemon.paused=true and rhiza.web.port=17891, then curl'd the real listener. Right after start: / answered 200 in 0.27 s and 0.14 s, /fragments/daemon in 0.019 s and 0.008 s, each with one 'Checking…' rz-pill--pending pill and no failed pills; 30 concurrent /fragments/daemon requests all answered under 0.1 s. About 10 s after start the report landed: / 0.16-0.17 s and /fragments/daemon 0.008-0.011 s showing 10 health pills (7 ok, 3 warn, e.g. 'daemon: ok · running…', 'autonomy: warn · daemon paused…'), live daemon status in them. Polling /fragments/daemon once a second for 45 s across the 30 s staleness mark: max 0.014 s, never pending again. Daemon CPU time: 40 s with no requests rose 1.0 s (poll loop only, no health run); one request on the stale report answered in 0.17 s and the next 15 s added 8.0 s CPU (one background health run), then 10 s idle added 0.3 s. `go test -trimpath -count=1 -run TestWebDaemon ./internal/api/` and `go test -trimpath ./internal/web/...` pass.

## Links

### Out

- verifies [web-overview-fast](../capability/rhiza/web-monitor/web-overview-fast.md)

## History

- `01M6QSJHF5WH4SZKCTTYD0W0SS` Fix The web UI's Overview page takes about 5 s to load (owner: "it's incredibly slow"). Measured on the VM: / 5.2 s and /fragments/daemon 4.7 s, while the inbox, live-runs and milestones sections and /runs and /roadmap answer in 0.01-0.11 s. The daemon panel's source (WebDaemonSource, internal/api/web_daemon.go) calls runHealth on every render, and health runs the full model check (rhiza health alone takes 5.0 s here). The page and the daemon section's live refreshes must answer fast (well under 0.5 s) without running the check on each request, e.g. by reusing the daemon's last health result instead of recomputing it.
