# web-serve-host-verified

evidence · `01M6QEPPV6S68RFDQZMXZY9VHQ` · actor [claude-code-ad0cd6e0](../actor/claude-code-ad0cd6e0.md)

## result

PASS

## statement

Verifier at 8f05e3ec: go test -trimpath -count=1 ./internal/web/... passed (TestAddr: default 127.0.0.1:7787, 0.0.0.0:7787, [::1]:7787, [::]:8080; localhost, abc, 127.0.0.1:80, [::1] rejected by Addr and ParseHost (net.ParseIP); TestManagerFollowsHost moves the listener 127.0.0.1->0.0.0.0:7787 and a bad host binds no socket). go test ./internal/api -run TestDaemonWeb passed, incl. TestDaemonWebListenerFollowsPort whose poll sequence ends with a host-only change to 0.0.0.0 that moves the fake listener to 0.0.0.0:7787. TestScripts/daemon-settings passed: lists rhiza.web.host = 127.0.0.1 (default), accepts 0.0.0.0 and ::1, rejects localhost, abc and '' with stderr naming rhiza.web.host. describe-golden passed; go build ./..., go vet and gofmt clean. Diff read: listen address is net.JoinHostPort(host, port), the daemon passes conf.webHost to Manager.Apply each poll, and daemon status uses the host too.

## Links

### Out

- verifies [web-serve-host](../capability/rhiza/web-monitor/web-serve-loopback/web-serve-host.md)

## History

- `01M6QBSDK2D7SYV177YBJCAM2A` Let the web UI listen on all interfaces so a Caddy reverse proxy on another host can reach it. Owner, verbatim: "in order to proxy the web ui via caddy we need to bind to all hosts and not just local host." This supersedes the localhost part of decision web-ui-in-daemon-read-only. Add a rhiza.web.host setting (default 127.0.0.1; 0.0.0.0 binds all interfaces) beside rhiza.web.port, settable with rhiza daemon set and followed live like the port; daemon status shows the bound address. The UI stays read-only.
