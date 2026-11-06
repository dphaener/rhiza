# web-serve-host

capability · `01M6QC3F43SZJGVMFP61832VJW` · parent [web-serve-loopback](../web-serve-loopback.md)

## acceptance

go test ./internal/web/... and go test ./internal/api -run TestDaemonWeb (the listener manager and its fake-listener poll test TestDaemonWebListenerFollowsPort, plus the rhiza.web.* validation in internal/api/daemon_settings.go) cover the computed address for the default host (127.0.0.1:7787), 0.0.0.0 (0.0.0.0:7787) and ::1 ([::1]:7787), and reject localhost, abc, an empty value and 127.0.0.1:80 with no socket bound; the fake-listener poll sequence includes a poll where only rhiza.web.host changes, to 0.0.0.0, and checks the listener moves to 0.0.0.0:<port>. The daemon-settings testscript lists rhiza.web.host with its default 127.0.0.1, accepts rhiza daemon set rhiza.web.host 0.0.0.0 and ::1, and rejects localhost, abc and an empty value with an error naming the key.

## statement

The rhiza.web.host machine setting, kept in the root's rhiza.toml, names the address the daemon serves the web UI on, 127.0.0.1 by default; rhiza daemon set accepts it only as a literal IPv4 or IPv6 address checked with net.ParseIP (such as 127.0.0.1, 0.0.0.0, ::1 or ::) and rejects an empty value, hostnames including localhost and anything with a port, with an error naming the key; the listen address is net.JoinHostPort(host, port), and a running daemon moves the listener at its next poll when only the host changes.

## Links

### Out

- decided-by [web-host-setting-shape](../../../../decision/web-host-setting-shape.md)
- decided-by [web-host-build-leaves](../../../../decision/web-host-build-leaves.md)

### In

- depends-on [web-status-shows-host](../web-url-in-daemon-status/web-status-shows-host.md)
- verifies [web-serve-host-verified](../../../../evidence/web-serve-host-verified.md)
- verifies [web-serve-host-qa-1](../../../../evidence/web-serve-host-qa-1.md)
- verifies [pr-325-let-the-web-ui-listen-on-all](../../../../evidence/pr-325-let-the-web-ui-listen-on-all.md)
- verifies [pr-338-root-daemon](../../../../evidence/pr-338-root-daemon.md)

## History

- `01M6QBSDK2D7SYV177YBJCAM2A` Let the web UI listen on all interfaces so a Caddy reverse proxy on another host can reach it. Owner, verbatim: "in order to proxy the web ui via caddy we need to bind to all hosts and not just local host." This supersedes the localhost part of decision web-ui-in-daemon-read-only. Add a rhiza.web.host setting (default 127.0.0.1; 0.0.0.0 binds all interfaces) beside rhiza.web.port, settable with rhiza daemon set and followed live like the port; daemon status shows the bound address. The UI stays read-only.
- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
