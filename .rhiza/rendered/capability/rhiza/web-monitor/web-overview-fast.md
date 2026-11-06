# web-overview-fast

capability · `01M6QSR4PW1N772QAPCE2E9N9W` · parent [web-monitor](../web-monitor.md)

## acceptance

A test serves / and /fragments/daemon with a health source that blocks for seconds and asserts each answers well under 0.5 s; that the first answer shows the pills pending, not failed; that many concurrent requests start one refresh; that once the report lands the pills show it; and that a report older than the interval starts one new refresh while requests keep getting the old one

## statement

The Overview page at / and its daemon fragment answer at once from the last health report the web UI holds in memory, refreshed in the background no more than one at a time and when older than about 30 s, with the health pills shown as pending until the first report lands, while running, paused and slots stay a live read

## Links

### Out

- decided-by [overview-health-cached-background](../../../decision/overview-health-cached-background.md)
- derived-from [overview-fast-plan](../../../source/overview-fast-plan.md)

### In

- verifies [web-overview-fast-verified](../../../evidence/web-overview-fast-verified.md)
- verifies [web-overview-fast-qa-1](../../../evidence/web-overview-fast-qa-1.md)
- verifies [web-overview-fast-verified-2](../../../evidence/web-overview-fast-verified-2.md)
- verifies [pr-327-fix-the-web-ui-s-overview-page](../../../evidence/pr-327-fix-the-web-ui-s-overview-page.md)

## History

- `01M6QSJHF5WH4SZKCTTYD0W0SS` Fix The web UI's Overview page takes about 5 s to load (owner: "it's incredibly slow"). Measured on the VM: / 5.2 s and /fragments/daemon 4.7 s, while the inbox, live-runs and milestones sections and /runs and /roadmap answer in 0.01-0.11 s. The daemon panel's source (WebDaemonSource, internal/api/web_daemon.go) calls runHealth on every render, and health runs the full model check (rhiza health alone takes 5.0 s here). The page and the daemon section's live refreshes must answer fast (well under 0.5 s) without running the check on each request, e.g. by reusing the daemon's last health result instead of recomputing it.
