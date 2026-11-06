# question-b6e0da8b-2

question · `01M6QC10Y9XWAETNDVS3HXEHTX`

## session

b6e0da8b-7ea6-46cb-ad32-01a1aebcef62

## statement

Change set 01M6QBS: per Decision web-host-setting-shape I amended the landed leaves web-serve-loopback and web-url-in-daemon-status in place. But build runs only build capabilities a change set creates (changeLeaves in internal/api/change.go:404, used by run_build.go:113 and daemon_fold.go:147), so the change set now has zero leaves and nothing would be built. How should the build work be carried?

## Links

### In

- answers [web-host-build-leaves](../decision/web-host-build-leaves.md)

## History

- `01M6QBSDK2D7SYV177YBJCAM2A` Let the web UI listen on all interfaces so a Caddy reverse proxy on another host can reach it. Owner, verbatim: "in order to proxy the web ui via caddy we need to bind to all hosts and not just local host." This supersedes the localhost part of decision web-ui-in-daemon-read-only. Add a rhiza.web.host setting (default 127.0.0.1; 0.0.0.0 binds all interfaces) beside rhiza.web.port, settable with rhiza daemon set and followed live like the port; daemon status shows the bound address. The UI stays read-only.
