# web-serve-loopback

capability · `01M6PFJ960GFE6G98GN9GVTABG` · parent [web-monitor](../web-monitor.md)

## acceptance

go test ./internal/web/... and ./internal/daemon/... cover the default 127.0.0.1:7787, a changed port, a changed host (0.0.0.0), an IPv6 host (::1, giving [::1]:7787), off, and rejected hosts (localhost, abc, empty, 127.0.0.1:80), through the handler and the computed address with no socket bound; a test drives the listener manager through successive polls with a fake listener factory: port 7787, then 8080, then host 0.0.0.0 with the port unchanged, then off, then 7787 again, and checks it starts, moves (including when only the host changes), stops and restarts the listener at each poll, and that --once and a set RHIZA_TEST_NOW open no listener. The daemon-settings testscript lists rhiza.web.port and rhiza.web.host with their defaults, accepts rhiza daemon set rhiza.web.port 8080 and off and rhiza.web.host 0.0.0.0 and ::1, and rejects port 0, 70000 and abc and host localhost, abc and an empty value with an error naming the key.

## statement

The long-running daemon serves the web UI over HTTP on the address rhiza.web.host names, 127.0.0.1 by default, at the port rhiza.web.port names, 7787 by default, listening on net.JoinHostPort(host, port), and serves nothing when rhiza.web.port is off; rhiza daemon --once never serves, and the daemon binds no listener when RHIZA_TEST_NOW is set, so the default test run binds no socket and parallel testscripts never race for the port. rhiza daemon set accepts rhiza.web.port only as a port number from 1 to 65535 or off, and rhiza.web.host only as a literal IPv4 or IPv6 address (such as 127.0.0.1, 0.0.0.0, ::1 or ::), rejecting an empty value, hostnames including localhost and anything with a port, with an error naming the key; a running daemon picks up a change to either at its next poll: it starts, moves or stops the listener to match. The UI stays read-only.

## Links

### Out

- decided-by [web-ui-in-daemon-read-only](../../../decision/web-ui-in-daemon-read-only.md)
- decided-by [web-host-setting-shape](../../../decision/web-host-setting-shape.md)

### In

- verifies [web-serve-loopback-verified](../../../evidence/web-serve-loopback-verified.md)
- verifies [web-serve-loopback-qa-1](../../../evidence/web-serve-loopback-qa-1.md)
- verifies [pr-314-web-ui-server](../../../evidence/pr-314-web-ui-server.md)
- verifies [pr-325-let-the-web-ui-listen-on-all](../../../evidence/pr-325-let-the-web-ui-listen-on-all.md)

## History

- `01M6PFEXSHWC6TR4F3GFKHV09R` Web server: the daemon serves a read-only web UI on 127.0.0.1:7787 (rhiza.web.port, or off), with read-only internal/api operations as JSON, server-sent events and embedded templates and assets
- `01M6QBSDK2D7SYV177YBJCAM2A` Let the web UI listen on all interfaces so a Caddy reverse proxy on another host can reach it. Owner, verbatim: "in order to proxy the web ui via caddy we need to bind to all hosts and not just local host." This supersedes the localhost part of decision web-ui-in-daemon-read-only. Add a rhiza.web.host setting (default 127.0.0.1; 0.0.0.0 binds all interfaces) beside rhiza.web.port, settable with rhiza daemon set and followed live like the port; daemon status shows the bound address. The UI stays read-only.
