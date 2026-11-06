# web-api-read-only

capability · `01M6PFJ9736AS6TE510QA5VVBS` · parent [web-monitor](../web-monitor.md)

## acceptance

A web package test compares /api/run.ls and /api/roadmap bodies against the api handlers' JSON for the same fixture repo and shows land, accept, daemon.set, render and ship.merge answer 404 or 405; an internal/api test asserts the ReadOnly marker is set on the operations the web UI uses and is not set on any operation that writes the log, settings, files, Rex or GitHub (daemon.set, daemon.start, render, ship.merge and every log write among them).

## statement

GET /api/<op> runs an internal/api operation that the registry marks ReadOnly, with its request built from the query string, and answers the same JSON bytes rhiza <cmd> --json prints; every operation without the ReadOnly marker, and every non-GET method, is refused, so no write is reachable.

## Links

### Out

- decided-by [web-ui-in-daemon-read-only](../../../decision/web-ui-in-daemon-read-only.md)
- decided-by [web-api-read-only-marker](../../../decision/web-api-read-only-marker.md)

### In

- blocks [review-01m6pf30-r3-f1](../../../question/review-01m6pf30-r3-f1.md)
- verifies [web-api-read-only-verified](../../../evidence/web-api-read-only-verified.md)
- verifies [web-api-read-only-qa-1](../../../evidence/web-api-read-only-qa-1.md)
- verifies [pr-314-web-ui-server](../../../evidence/pr-314-web-ui-server.md)

## History

- `01M6PFEXSHWC6TR4F3GFKHV09R` Web server: the daemon serves a read-only web UI on 127.0.0.1:7787 (rhiza.web.port, or off), with read-only internal/api operations as JSON, server-sent events and embedded templates and assets
