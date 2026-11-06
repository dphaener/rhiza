# m28-browser-monitor

intent · `01M6PFEK8B9YWN12YC15CAJZH0` · parent [milestone-28](../milestone-28.md)

## statement

With the daemon running, the owner opens http://127.0.0.1:7787 (through their own port forward) and sees the daemon and its health, live runs and their current step, what waits in the inbox, and each milestone's change sets, without a terminal.

## Links

### In

- serves change set `01M6QSJHF5WH4SZKCTTYD0W0SS` (Fix The web UI's Overview page takes about 5 s to load (owner: "it's incredibly slow"). Measured on the VM: / 5.2 s and /fragments/daemon 4.7 s, while the inbox, live-runs and milestones sections and /runs and /roadmap answer in 0.01-0.11 s. The daemon panel's source (WebDaemonSource, internal/api/web_daemon.go) calls runHealth on every render, and health runs the full model check (rhiza health alone takes 5.0 s here). The page and the daemon section's live refreshes must answer fast (well under 0.5 s) without running the check on each request, e.g. by reusing the daemon's last health result instead of recomputing it.)

## History

- `01M6PF31YCJ041HB94MDWQR331` I want to build out the first pass of a web UI that can be used to monitor the progress of Rhiza. This will run as part of the daemon process and be served on localhost. The user will be responsible for any port forwarding (if it's running on this VM for instance). First pass will be read only. For monitoring purposes only and we will build from there. The first requirement will be to obtain a design. You will generate a prompt that I will then feed to Claude design to generate a design system.
