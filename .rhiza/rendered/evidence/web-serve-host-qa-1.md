# web-serve-host-qa-1

evidence · `01M6QF4HPXM36YPFZH7QHVPD3V` · actor [claude-code-6065d5f4](../actor/claude-code-6065d5f4.md)

## result

PASS

## statement

QA at 5ad7b858 with the built binary in a scratch repo, live daemon (not --once), rhiza.web.port 17991: rhiza daemon set rhiza.web.host rejected localhost, abc, '' and 127.0.0.1:80 (exit 2, naming rhiza.web.host) and wrote nothing. Default bound 127.0.0.1:17991 only (ss; curl loopback 200, LAN 192.168.1.210 refused). While the daemon ran, set rhiza.web.host 0.0.0.0: the listener moved on a later poll to *:17991 and the UI served 200 over LAN and Tailscale IPs; POST to /, /api/call/change.new and /api/daemon.set from the LAN answered 405 only GET is allowed. Set 192.0.2.7 (not on the host): nothing listened, daemon stderr logged the bind error each poll. Set ::1: bound [::1]:17991 (200 on ::1, LAN refused). Port off: listener closed.

## Links

### Out

- verifies [web-serve-host](../capability/rhiza/web-monitor/web-serve-loopback/web-serve-host.md)

## History

- `01M6QBSDK2D7SYV177YBJCAM2A` Let the web UI listen on all interfaces so a Caddy reverse proxy on another host can reach it. Owner, verbatim: "in order to proxy the web ui via caddy we need to bind to all hosts and not just local host." This supersedes the localhost part of decision web-ui-in-daemon-read-only. Add a rhiza.web.host setting (default 127.0.0.1; 0.0.0.0 binds all interfaces) beside rhiza.web.port, settable with rhiza daemon set and followed live like the port; daemon status shows the bound address. The UI stays read-only.
