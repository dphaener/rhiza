# web-status-shows-host-qa-1

evidence · `01M6QF4HRHB3JND14V9PBB004J` · actor [claude-code-6065d5f4](../actor/claude-code-6065d5f4.md)

## result

PASS

## statement

QA at 5ad7b858 with the built binary and a live daemon in a scratch repo on port 17991: rhiza daemon status showed Web UI: http://127.0.0.1:17991 at the default, http://0.0.0.0:17991 after rhiza daemon set rhiza.web.host 0.0.0.0 and a poll, http://[::1]:17991 (bracketed) after ::1, and Web UI: off after rhiza.web.port off. go test ./cmd/rhiza -run TestScripts/(daemon-settings|daemon-status|describe-golden) passed.

## Links

### Out

- verifies [web-status-shows-host](../capability/rhiza/web-monitor/web-url-in-daemon-status/web-status-shows-host.md)

## History

- `01M6QBSDK2D7SYV177YBJCAM2A` Let the web UI listen on all interfaces so a Caddy reverse proxy on another host can reach it. Owner, verbatim: "in order to proxy the web ui via caddy we need to bind to all hosts and not just local host." This supersedes the localhost part of decision web-ui-in-daemon-read-only. Add a rhiza.web.host setting (default 127.0.0.1; 0.0.0.0 binds all interfaces) beside rhiza.web.port, settable with rhiza daemon set and followed live like the port; daemon status shows the bound address. The UI stays read-only.
