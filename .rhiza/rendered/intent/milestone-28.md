# milestone-28

intent · `01M6PFEK82MGSAY2M7055FVPCV`

## milestone

28

## statement

The owner monitors Rhiza from a browser: the daemon serves a read-only web UI on localhost, designed from the owner's Claude Design system, showing the daemon and its health, live and recent runs, the inbox and the roadmap of milestones and change sets, kept live without reloading.

## Links

### Out

- decided-by [web-ui-in-daemon-read-only](../decision/web-ui-in-daemon-read-only.md)
- decided-by [web-ui-design-from-claude-design](../decision/web-ui-design-from-claude-design.md)

### In

- serves change set `01M6PFEXSHWC6TR4F3GFKHV09R` (Web server: the daemon serves a read-only web UI on 127.0.0.1:7787 (rhiza.web.port, or off), with read-only internal/api operations as JSON, server-sent events and embedded templates and assets)
- serves change set `01M6PFFFZNADHH2WGMB5MES8QT` (Web design system: the Claude Design system built from the owner's basics as one embedded stylesheet, embedded fonts and the shared layout and partials every page renders with)
- serves change set `01M6PFGXV205YBYX2MZ9CC590B` (Web Overview page: the daemon and its health, live runs and their current step, the inbox and the active milestone's progress on one live page)
- serves change set `01M6PFH1X52KNK0C5ZCDYEEF64` (Web Runs pages: runs live first then recent, and a run's detail with its steps, sessions, review summary and a live log tail)
- serves change set `01M6PFH625P3BYQAB8P6J38BBT` (Web Roadmap pages: milestone groups and the backlog, and a change set's page with its intent, state, leaves, questions, runs and history)
- serves change set `01M6QBSDK2D7SYV177YBJCAM2A` (Let the web UI listen on all interfaces so a Caddy reverse proxy on another host can reach it. Owner, verbatim: "in order to proxy the web ui via caddy we need to bind to all hosts and not just local host." This supersedes the localhost part of decision web-ui-in-daemon-read-only. Add a rhiza.web.host setting (default 127.0.0.1; 0.0.0.0 binds all interfaces) beside rhiza.web.port, settable with rhiza daemon set and followed live like the port; daemon status shows the bound address. The UI stays read-only.)

## History

- `01M6PF31YCJ041HB94MDWQR331` I want to build out the first pass of a web UI that can be used to monitor the progress of Rhiza. This will run as part of the daemon process and be served on localhost. The user will be responsible for any port forwarding (if it's running on this VM for instance). First pass will be read only. For monitoring purposes only and we will build from there. The first requirement will be to obtain a design. You will generate a prompt that I will then feed to Claude design to generate a design system.
