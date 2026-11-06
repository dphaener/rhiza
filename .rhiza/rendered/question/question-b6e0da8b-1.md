# question-b6e0da8b-1

question · `01M6QBW5TMT8B4DASNECB5AGGW`

## session

b6e0da8b-7ea6-46cb-ad32-01a1aebcef62

## statement

Change set 01M6QBS (web UI host setting) serves no intent yet, and milestone 28's done point m28-read-only says the UI 'binds only to 127.0.0.1', which this change contradicts. What should it serve, and what happens to that done point? / Which values should rhiza daemon set rhiza.web.host accept? / How should the existing landed leaves web-serve-loopback (asserts always 127.0.0.1) and web-url-in-daemon-status (shows http://127.0.0.1:<port>) be handled?

## Links

### In

- answers [web-host-setting-shape](../decision/web-host-setting-shape.md)

## History

- `01M6QBSDK2D7SYV177YBJCAM2A` Let the web UI listen on all interfaces so a Caddy reverse proxy on another host can reach it. Owner, verbatim: "in order to proxy the web ui via caddy we need to bind to all hosts and not just local host." This supersedes the localhost part of decision web-ui-in-daemon-read-only. Add a rhiza.web.host setting (default 127.0.0.1; 0.0.0.0 binds all interfaces) beside rhiza.web.port, settable with rhiza daemon set and followed live like the port; daemon status shows the bound address. The UI stays read-only.
