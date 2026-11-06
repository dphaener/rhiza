# web-status-shows-host-verified

evidence · `01M6QESHD0C6CESJTCT67076ZT` · actor [claude-code-9dc87e15](../actor/claude-code-9dc87e15.md)

## result

PASS

## statement

At edc98e1b ran go test -trimpath ./cmd/rhiza -run 'TestScripts/^daemon-status$' -count=1: PASS. The script, after rhiza daemon --once, asserts 'Web UI: http://127.0.0.1:7787' in text and "web": "http://127.0.0.1:7787" in --json; after daemon set rhiza.web.host 0.0.0.0 and another --once, http://0.0.0.0:7787 in both; after daemon set rhiza.web.host ::1 and another --once, http://[::1]:7787 in both. Read the code: webStatus in internal/api/daemon.go builds the URL from the recorded host and port settings via web.Addr, which uses net.JoinHostPort, so IPv6 hosts are bracketed and an empty host defaults to 127.0.0.1.

## Links

### Out

- verifies [web-status-shows-host](../capability/rhiza/web-monitor/web-url-in-daemon-status/web-status-shows-host.md)

## History

- `01M6QBSDK2D7SYV177YBJCAM2A` Let the web UI listen on all interfaces so a Caddy reverse proxy on another host can reach it. Owner, verbatim: "in order to proxy the web ui via caddy we need to bind to all hosts and not just local host." This supersedes the localhost part of decision web-ui-in-daemon-read-only. Add a rhiza.web.host setting (default 127.0.0.1; 0.0.0.0 binds all interfaces) beside rhiza.web.port, settable with rhiza daemon set and followed live like the port; daemon status shows the bound address. The UI stays read-only.
