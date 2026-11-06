# web-ui-in-daemon-read-only

decision · `01M6PFEKAMCWXEK0XHHRQ4C088`

## quote

This will run as part of the daemon process and be served on localhost. The user will be responsible for any port forwarding (if it's running on this VM for instance). First pass will be read only. For monitoring purposes only and we will build from there.

## statement

The first web UI is served by the daemon process on localhost, read-only and for monitoring only; the owner handles any port forwarding.

## Links

### Out

- derived-from [web-ui-milestone-interview](../source/web-ui-milestone-interview.md)

### In

- relates-to [web-host-setting-shape](web-host-setting-shape.md) — Amends only the 'served on localhost' clause: the UI binds to the address rhiza.web.host names (127.0.0.1 by default). The daemon-served, read-only, monitoring-only parts of web-ui-in-daemon-read-only still hold. (derived-from was ruled in web-host-amends-not-supersedes, but derived-from only targets Sources; relates-to carries the same amendment note.)
- decided-by [milestone-28](../intent/milestone-28.md)
- decided-by [web-monitor](../capability/rhiza/web-monitor.md)
- decided-by [web-serve-loopback](../capability/rhiza/web-monitor/web-serve-loopback.md)
- decided-by [web-api-read-only](../capability/rhiza/web-monitor/web-api-read-only.md)
- decided-by [web-events-stream](../capability/rhiza/web-monitor/web-events-stream.md)
- decided-by [web-templates-embedded](../capability/rhiza/web-monitor/web-templates-embedded.md)
- decided-by [web-live-sections](../capability/rhiza/web-monitor/web-live-sections.md)
- decided-by [web-url-in-daemon-status](../capability/rhiza/web-monitor/web-url-in-daemon-status.md)
- decided-by [web-overview-daemon](../capability/rhiza/web-monitor/web-overview-daemon.md)
- decided-by [web-overview-runs](../capability/rhiza/web-monitor/web-overview-runs.md)
- decided-by [web-overview-inbox](../capability/rhiza/web-monitor/web-overview-inbox.md)
- decided-by [web-overview-milestone](../capability/rhiza/web-monitor/web-overview-milestone.md)
- decided-by [web-overview-live](../capability/rhiza/web-monitor/web-overview-live.md)
- decided-by [web-runs-list](../capability/rhiza/web-monitor/web-runs-list.md)
- decided-by [web-run-detail](../capability/rhiza/web-monitor/web-run-detail.md)
- decided-by [web-run-log-tail](../capability/rhiza/web-monitor/web-run-log-tail.md)
- decided-by [web-roadmap-page](../capability/rhiza/web-monitor/web-roadmap-page.md)
- decided-by [web-change-set-page](../capability/rhiza/web-monitor/web-change-set-page.md)

## History

- `01M6PF31YCJ041HB94MDWQR331` I want to build out the first pass of a web UI that can be used to monitor the progress of Rhiza. This will run as part of the daemon process and be served on localhost. The user will be responsible for any port forwarding (if it's running on this VM for instance). First pass will be read only. For monitoring purposes only and we will build from there. The first requirement will be to obtain a design. You will generate a prompt that I will then feed to Claude design to generate a design system.
