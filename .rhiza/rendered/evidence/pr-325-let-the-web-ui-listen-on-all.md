# pr-325-let-the-web-ui-listen-on-all

evidence · `01M6QG10TCB9D8KG7MTJGV5EVX` · actor [dphaener](../actor/dphaener.md)

## head

fa72fba2226ce311695bb953ed3c2db7075c709e

## pr

https://github.com/dphaener/rhiza/pull/325

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #325 (https://github.com/dphaener/rhiza/pull/325) passed CI at fa72fba

## status

verified

## Links

### Out

- verifies [web-serve-loopback](../capability/rhiza/web-monitor/web-serve-loopback.md)
- verifies [web-url-in-daemon-status](../capability/rhiza/web-monitor/web-url-in-daemon-status.md)
- verifies [web-serve-host](../capability/rhiza/web-monitor/web-serve-loopback/web-serve-host.md)
- verifies [web-status-shows-host](../capability/rhiza/web-monitor/web-url-in-daemon-status/web-status-shows-host.md)

## History

- `01M6QBSDK2D7SYV177YBJCAM2A` Let the web UI listen on all interfaces so a Caddy reverse proxy on another host can reach it. Owner, verbatim: "in order to proxy the web ui via caddy we need to bind to all hosts and not just local host." This supersedes the localhost part of decision web-ui-in-daemon-read-only. Add a rhiza.web.host setting (default 127.0.0.1; 0.0.0.0 binds all interfaces) beside rhiza.web.port, settable with rhiza daemon set and followed live like the port; daemon status shows the bound address. The UI stays read-only.
