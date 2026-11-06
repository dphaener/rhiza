# overview-health-cached-background

decision · `01M6QSQ4GJQ8QRYX9JND6PMSRQ`

## reasoning

The change set's intent asks for / and /fragments/daemon well under 0.5 s without running the check per request, and suggests reusing a last health result; the question notes none is kept today. Reading internal/api/web_daemon.go and health.go: runHealth's cost is the full rhiza check (plus roadmap), while runDaemonStatus is a cheap file read. The web handler is built inside runDaemon (internal/api/daemon.go) from a ctx snapshot, so an in-process cache needs no new storage. Saving health into daemon state on each poll (the usage-window pattern) would add about 5 s to every poll and couple poll cadence to UI freshness, so it is rejected. A bounded background refresh with single-flight keeps requests fast, keeps CPU bounded no matter how often the page or its live refreshes poll, and keeps the fix inside the web source, the smallest change that meets the intent. Health a few tens of seconds stale is acceptable for an at-a-glance panel; rhiza health stays the fresh, full report.

## statement

The Overview's daemon panel never runs health in the request path. WebDaemonSource keeps the last health report in memory (the web UI already lives in the daemon process) and refreshes it in the background: a request that finds the report missing or older than a short interval (about 30 s) starts one refresh (never more than one in flight) and answers at once with what it has. Before the first report lands the health pills show as pending (no result yet), not as failed. The daemon status part of the panel (running, slots, limit, paused) stays a live read on every request, since it is cheap. The daemon's poll loop and its state file are not changed: no health result is written to state and the poll does not run the check.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-168498bf-1](../question/question-168498bf-1.md)
- derived-from [overview-fast-plan](../source/overview-fast-plan.md)

### In

- decided-by [web-overview-fast](../capability/rhiza/web-monitor/web-overview-fast.md)

## History

- `01M6QSJHF5WH4SZKCTTYD0W0SS` Fix The web UI's Overview page takes about 5 s to load (owner: "it's incredibly slow"). Measured on the VM: / 5.2 s and /fragments/daemon 4.7 s, while the inbox, live-runs and milestones sections and /runs and /roadmap answer in 0.01-0.11 s. The daemon panel's source (WebDaemonSource, internal/api/web_daemon.go) calls runHealth on every render, and health runs the full model check (rhiza health alone takes 5.0 s here). The page and the daemon section's live refreshes must answer fast (well under 0.5 s) without running the check on each request, e.g. by reusing the daemon's last health result instead of recomputing it.
