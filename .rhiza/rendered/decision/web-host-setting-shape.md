# web-host-setting-shape

decision · `01M6QBXQ8SJ1HGA3RBETE3S7X4`

## reasoning

The change set's intent quotes the owner ('we need to bind to all hosts and not just local host') and says it supersedes the localhost part of web-ui-in-daemon-read-only, so the work belongs to the same milestone (milestone-28, whose web server change set introduced rhiza.web.port) and the conflicting clause of m28-read-only should be corrected, not left contradicting the model; the read-only and off parts of that done point are untouched because the intent says the UI stays read-only. For values, the intent names only 127.0.0.1 and 0.0.0.0; literal IPs cover those plus IPv6 with no DNS lookup at bind time and mirror the strict, error-naming validation ParsePort already applies to rhiza.web.port in internal/web/web.go, while hostnames would add resolution ambiguity the Caddy use case does not need. For the leaves, web-serve-loopback and web-url-in-daemon-status already own binding and the status URL and are childless leaves, so amending them keeps one capability per behaviour and makes a build run rebuild exactly those; adding new leaves would leave the old ones asserting 'always 127.0.0.1', contradicting the new setting. Showing the bound address as recorded matches the intent's 'daemon status shows the bound address' and keeps the default output unchanged.

## statement

Change set 01M6QBS serves milestone-28 (the web UI milestone), and amends its done point m28-read-only in place: the UI binds to the address rhiza.web.host names, 127.0.0.1 by default, instead of 'only to 127.0.0.1'; the allowlisted read-only operations and 'rhiza.web.port off leaves the daemon as it was' stay unchanged. rhiza daemon set rhiza.web.host accepts only a literal IP address (IPv4 or IPv6, e.g. 127.0.0.1, 0.0.0.0, ::1, ::), checked with net.ParseIP; it rejects an empty value, hostnames (including localhost) and anything with a port, with an error naming the key; the listen address is net.JoinHostPort(host, port). The landed leaves are amended in this change set rather than shadowed by new ones: web-serve-loopback's statement and acceptance become 'serves on the address rhiza.web.host names (default 127.0.0.1) at rhiza.web.port', its tests cover the default 127.0.0.1, a changed host, IPv6, and the listener manager moving the listener when only the host changes, and the daemon-settings testscript lists rhiza.web.host with its default, accepts 0.0.0.0 and ::1 and rejects localhost, abc and an empty value; web-url-in-daemon-status shows the recorded bound address as http://<host>:<port> (IPv6 bracketed), so the default still reads http://127.0.0.1:7787 and 0.0.0.0 reads http://0.0.0.0:7787.

## wording

(assumed, would have escalated)

## Links

### Out

- relates-to [web-ui-in-daemon-read-only](web-ui-in-daemon-read-only.md) — Amends only the 'served on localhost' clause: the UI binds to the address rhiza.web.host names (127.0.0.1 by default). The daemon-served, read-only, monitoring-only parts of web-ui-in-daemon-read-only still hold. (derived-from was ruled in web-host-amends-not-supersedes, but derived-from only targets Sources; relates-to carries the same amendment note.)
- answers [question-b6e0da8b-1](../question/question-b6e0da8b-1.md)
- derived-from [web-host-plan-interview](../source/web-host-plan-interview.md)

### In

- decided-by [m28-read-only](../intent/milestone-28/m28-read-only.md)
- decided-by [web-serve-loopback](../capability/rhiza/web-monitor/web-serve-loopback.md)
- decided-by [web-url-in-daemon-status](../capability/rhiza/web-monitor/web-url-in-daemon-status.md)
- decided-by [web-serve-host](../capability/rhiza/web-monitor/web-serve-loopback/web-serve-host.md)
- decided-by [web-status-shows-host](../capability/rhiza/web-monitor/web-url-in-daemon-status/web-status-shows-host.md)

## History

- `01M6QBSDK2D7SYV177YBJCAM2A` Let the web UI listen on all interfaces so a Caddy reverse proxy on another host can reach it. Owner, verbatim: "in order to proxy the web ui via caddy we need to bind to all hosts and not just local host." This supersedes the localhost part of decision web-ui-in-daemon-read-only. Add a rhiza.web.host setting (default 127.0.0.1; 0.0.0.0 binds all interfaces) beside rhiza.web.port, settable with rhiza daemon set and followed live like the port; daemon status shows the bound address. The UI stays read-only.
