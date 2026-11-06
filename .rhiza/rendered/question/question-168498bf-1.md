# question-168498bf-1

question · `01M6QSNV3XH5D44ES3Y10TG239`

## session

168498bf-491f-475d-a975-01a1af9967ce

## statement

Fix for the slow web UI Overview (WebDaemonSource runs the full health report, rhiza check included, about 5 s, on every render of / and /fragments/daemon). The daemon keeps no health result today, so there is nothing to reuse yet. How should the Overview get health without running it on each request?

## Links

### In

- answers [overview-health-cached-background](../decision/overview-health-cached-background.md)

## History

- `01M6QSJHF5WH4SZKCTTYD0W0SS` Fix The web UI's Overview page takes about 5 s to load (owner: "it's incredibly slow"). Measured on the VM: / 5.2 s and /fragments/daemon 4.7 s, while the inbox, live-runs and milestones sections and /runs and /roadmap answer in 0.01-0.11 s. The daemon panel's source (WebDaemonSource, internal/api/web_daemon.go) calls runHealth on every render, and health runs the full model check (rhiza health alone takes 5.0 s here). The page and the daemon section's live refreshes must answer fast (well under 0.5 s) without running the check on each request, e.g. by reusing the daemon's last health result instead of recomputing it.
