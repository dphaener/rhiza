# web-status-shows-host

capability · `01M6QC3F6QYZ63H13RNGX78808` · parent [web-url-in-daemon-status](../web-url-in-daemon-status.md)

## acceptance

The daemon-status testscript, after rhiza daemon --once, shows http://127.0.0.1:7787 in text and --json; after rhiza daemon set rhiza.web.host 0.0.0.0 and another --once, it shows http://0.0.0.0:7787; after rhiza daemon set rhiza.web.host ::1 and another --once, it shows http://[::1]:7787.

## statement

rhiza daemon status shows the web UI's URL with the effective rhiza.web.host the daemon's last poll recorded, as http://<host>:<port> with an IPv6 host bracketed, so the default still reads http://127.0.0.1:7787.

## Links

### Out

- depends-on [web-serve-host](../web-serve-loopback/web-serve-host.md)
- decided-by [web-host-setting-shape](../../../../decision/web-host-setting-shape.md)
- decided-by [web-host-build-leaves](../../../../decision/web-host-build-leaves.md)

### In

- verifies [web-status-shows-host-verified](../../../../evidence/web-status-shows-host-verified.md)
- verifies [web-status-shows-host-qa-1](../../../../evidence/web-status-shows-host-qa-1.md)
- verifies [pr-325-let-the-web-ui-listen-on-all](../../../../evidence/pr-325-let-the-web-ui-listen-on-all.md)

## History

- `01M6QBSDK2D7SYV177YBJCAM2A` Let the web UI listen on all interfaces so a Caddy reverse proxy on another host can reach it. Owner, verbatim: "in order to proxy the web ui via caddy we need to bind to all hosts and not just local host." This supersedes the localhost part of decision web-ui-in-daemon-read-only. Add a rhiza.web.host setting (default 127.0.0.1; 0.0.0.0 binds all interfaces) beside rhiza.web.port, settable with rhiza daemon set and followed live like the port; daemon status shows the bound address. The UI stays read-only.
