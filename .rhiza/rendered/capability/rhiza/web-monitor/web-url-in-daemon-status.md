# web-url-in-daemon-status

capability · `01M6PFJ9CD8XKQZK2PNA3B6973` · parent [web-monitor](../web-monitor.md)

## acceptance

The daemon-status testscript runs rhiza daemon --once, then rhiza daemon status shows http://127.0.0.1:7787 in text and --json; after rhiza daemon set rhiza.web.port off and another rhiza daemon --once, status shows the web UI off.

## statement

rhiza daemon status shows the web UI as the daemon's state file recorded it from the last poll's effective rhiza.web.port: its URL (http://127.0.0.1:<port>) for a port, or that the web UI is off; it reads the recorded setting and does not probe the listener.

## Links

### Out

- decided-by [web-ui-in-daemon-read-only](../../../decision/web-ui-in-daemon-read-only.md)

### In

- verifies [web-url-in-daemon-status-verified](../../../evidence/web-url-in-daemon-status-verified.md)
- verifies [web-url-in-daemon-status-qa-1](../../../evidence/web-url-in-daemon-status-qa-1.md)
- verifies [pr-314-web-ui-server](../../../evidence/pr-314-web-ui-server.md)

## History

- `01M6PFEXSHWC6TR4F3GFKHV09R` Web server: the daemon serves a read-only web UI on 127.0.0.1:7787 (rhiza.web.port, or off), with read-only internal/api operations as JSON, server-sent events and embedded templates and assets
