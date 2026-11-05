# web-events-stream

capability · `01M6PFJ98BKT1ES8F6WJTVJ893` · parent [web-monitor](../web-monitor.md)

## acceptance

A web package test drives the events handler with a recorder, appends a log entry, a run journal move, driver log lines and an inbox item and changes the daemon's status in a fixture repo, and reads one matching named event for each from the stream.

## statement

GET /events streams server-sent events naming what changed (a change set moved or landed, a run's journal moved, a run's driver log grew, an inbox item arrived or left, the daemon's status changed: paused, slots or health), so an open page knows to re-fetch.

## Links

### Out

- decided-by [web-ui-in-daemon-read-only](../../../decision/web-ui-in-daemon-read-only.md)

### In

- verifies [web-events-stream-verified](../../../evidence/web-events-stream-verified.md)
- verifies [web-events-stream-qa-1](../../../evidence/web-events-stream-qa-1.md)
- verifies [pr-314-web-ui-server](../../../evidence/pr-314-web-ui-server.md)

## History

- `01M6PFEXSHWC6TR4F3GFKHV09R` Web server: the daemon serves a read-only web UI on 127.0.0.1:7787 (rhiza.web.port, or off), with read-only internal/api operations as JSON, server-sent events and embedded templates and assets
