# web-url-in-daemon-status

capability · `01M6PFJ9CD8XKQZK2PNA3B6973` · parent [web-monitor](../web-monitor.md)

## acceptance

The daemon-status testscript runs rhiza daemon --once, then rhiza daemon status shows http://127.0.0.1:7787 in text and --json; after rhiza daemon set rhiza.web.host 0.0.0.0 and another rhiza daemon --once, status shows http://0.0.0.0:7787; after rhiza daemon set rhiza.web.host ::1 and another --once, it shows http://[::1]:7787; after rhiza daemon set rhiza.web.port off and another --once, status shows the web UI off.

## statement

rhiza daemon status shows the web UI as the daemon's state file recorded it from the last poll's effective rhiza.web.host and rhiza.web.port: its URL http://<host>:<port> (an IPv6 host bracketed) for a port, or that the web UI is off; it reads the recorded settings and does not probe the listener.

## Links

### Out

- decided-by [web-ui-in-daemon-read-only](../../../decision/web-ui-in-daemon-read-only.md)
- decided-by [web-host-setting-shape](../../../decision/web-host-setting-shape.md)

### In

- verifies [web-url-in-daemon-status-verified](../../../evidence/web-url-in-daemon-status-verified.md)
- verifies [web-url-in-daemon-status-qa-1](../../../evidence/web-url-in-daemon-status-qa-1.md)
- verifies [pr-314-web-ui-server](../../../evidence/pr-314-web-ui-server.md)
- verifies [pr-325-let-the-web-ui-listen-on-all](../../../evidence/pr-325-let-the-web-ui-listen-on-all.md)

## History

- `01M6PFEXSHWC6TR4F3GFKHV09R` Web server: the daemon serves a read-only web UI on 127.0.0.1:7787 (rhiza.web.port, or off), with read-only internal/api operations as JSON, server-sent events and embedded templates and assets
- `01M6QBSDK2D7SYV177YBJCAM2A` Let the web UI listen on all interfaces so a Caddy reverse proxy on another host can reach it. Owner, verbatim: "in order to proxy the web ui via caddy we need to bind to all hosts and not just local host." This supersedes the localhost part of decision web-ui-in-daemon-read-only. Add a rhiza.web.host setting (default 127.0.0.1; 0.0.0.0 binds all interfaces) beside rhiza.web.port, settable with rhiza daemon set and followed live like the port; daemon status shows the bound address. The UI stays read-only.
