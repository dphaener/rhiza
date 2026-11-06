# m28-read-only

intent · `01M6PFEK8PKRGYC4C6WZS9BY5Y` · parent [milestone-28](../milestone-28.md)

## statement

Nothing the web UI serves writes: only allowlisted read operations of internal/api are reachable, it binds to the address rhiza.web.host names (127.0.0.1 by default), and turning it off (rhiza.web.port off) leaves the daemon as it was.

## Links

### Out

- decided-by [web-host-setting-shape](../../decision/web-host-setting-shape.md)

## History

- `01M6PF31YCJ041HB94MDWQR331` I want to build out the first pass of a web UI that can be used to monitor the progress of Rhiza. This will run as part of the daemon process and be served on localhost. The user will be responsible for any port forwarding (if it's running on this VM for instance). First pass will be read only. For monitoring purposes only and we will build from there. The first requirement will be to obtain a design. You will generate a prompt that I will then feed to Claude design to generate a design system.
- `01M6QBSDK2D7SYV177YBJCAM2A` Let the web UI listen on all interfaces so a Caddy reverse proxy on another host can reach it. Owner, verbatim: "in order to proxy the web ui via caddy we need to bind to all hosts and not just local host." This supersedes the localhost part of decision web-ui-in-daemon-read-only. Add a rhiza.web.host setting (default 127.0.0.1; 0.0.0.0 binds all interfaces) beside rhiza.web.port, settable with rhiza daemon set and followed live like the port; daemon status shows the bound address. The UI stays read-only.
