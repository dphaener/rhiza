# web-serve-loopback

capability · `01M6PFJ960GFE6G98GN9GVTABG` · parent [web-monitor](../web-monitor.md)

## acceptance

go test ./internal/web/... and ./internal/daemon/... cover the default 7787, a changed port, off, and that the listen address is always 127.0.0.1, through the handler and the computed address with no socket bound; a test drives the listener manager through successive polls with rhiza.web.port 7787, then 8080, then off, then 7787 again, using a fake listener factory, and checks it starts, moves, stops and restarts the listener at each poll, and that --once and a set RHIZA_TEST_NOW open no listener. The daemon-settings testscript lists rhiza.web.port with its default, accepts rhiza daemon set rhiza.web.port 8080 and off, and rejects 0, 70000 and abc with an error naming the key.

## statement

The long-running daemon serves the web UI over HTTP on 127.0.0.1 at the port rhiza.web.port names, 7787 by default, and serves nothing when it is set to off; rhiza daemon --once never serves, and the daemon binds no listener when RHIZA_TEST_NOW is set, so the default test run binds no socket and parallel testscripts never race for the port. rhiza daemon set accepts rhiza.web.port only as a port number from 1 to 65535 or off, and a running daemon picks up a change at its next poll: it starts, moves or stops the listener to match.

## Links

### Out

- decided-by [web-ui-in-daemon-read-only](../../../decision/web-ui-in-daemon-read-only.md)

### In

- verifies [web-serve-loopback-verified](../../../evidence/web-serve-loopback-verified.md)
- verifies [web-serve-loopback-qa-1](../../../evidence/web-serve-loopback-qa-1.md)
- verifies [pr-314-web-ui-server](../../../evidence/pr-314-web-ui-server.md)

## History

- `01M6PFEXSHWC6TR4F3GFKHV09R` Web server: the daemon serves a read-only web UI on 127.0.0.1:7787 (rhiza.web.port, or off), with read-only internal/api operations as JSON, server-sent events and embedded templates and assets
