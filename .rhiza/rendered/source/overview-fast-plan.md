# overview-fast-plan

source · `01M6QSQXRDHBMB5YX5DPN0AZ26`

## quotes

```json
[
  {
    "by": "agent (recommended; the arbiter decided it as overview-health-cached-background)",
    "text": "Web-side background cache: the daemon's web server computes health in the background, single-flight, pages render the last result and never wait for it; before the first result the panel shows health as checking.",
    "where": "interview question 1"
  }
]
```

## raw_sha256

490d16aba5b367866820c2252a2be72d8308658fd7bc6dcfb715a58cb4289ec4

## session

168498bf-491f-475d-a975-01a1af9967ce

## summary

Autonomous plan run for the slow Overview page (/ and /fragments/daemon about 5 s, because WebDaemonSource runs the full health report with rhiza check on every render). The daemon keeps no health result to reuse today. One question was raised: how the Overview gets health without running it per request. The arbiter decided on an in-memory health report in WebDaemonSource, refreshed in the background (single flight, about 30 s staleness), pills shown as pending before the first report, daemon status read live, no change to the daemon poll or its state file. Planned one leaf under web-monitor, serving m28-browser-monitor. Nothing left open.

## title

Fast web UI Overview, autonomous plan

## Links

### In

- derived-from [overview-health-cached-background](../decision/overview-health-cached-background.md)
- derived-from [web-overview-fast](../capability/rhiza/web-monitor/web-overview-fast.md)

## History

- `01M6QSJHF5WH4SZKCTTYD0W0SS` Fix The web UI's Overview page takes about 5 s to load (owner: "it's incredibly slow"). Measured on the VM: / 5.2 s and /fragments/daemon 4.7 s, while the inbox, live-runs and milestones sections and /runs and /roadmap answer in 0.01-0.11 s. The daemon panel's source (WebDaemonSource, internal/api/web_daemon.go) calls runHealth on every render, and health runs the full model check (rhiza health alone takes 5.0 s here). The page and the daemon section's live refreshes must answer fast (well under 0.5 s) without running the check on each request, e.g. by reusing the daemon's last health result instead of recomputing it.
