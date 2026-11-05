# web-live-sections-verified

evidence · `01M6PTWJZDDCC064TCZN8A83TG` · actor [claude-code-1840efdb](../actor/claude-code-1840efdb.md)

## result

PASS

## statement

At 74600e5d ran go test -trimpath -count=1 ./internal/web/... (ok) and go vet ./internal/web (clean). TestLiveSectionAndFragment registers a sample section with events run and log, renders it via section() and sees data-rz-fragment="/fragments/sample", data-rz-events="run log" and its content; GET /fragments/sample returns exactly <p>runs: 3</p>, an unknown fragment 404s and POST is 405; GET /static/app.js is served and contains data-rz-fragment, /events and setStale(true). Read the diff: app.js listens for the named SSE kinds (change, run, log, inbox, daemon, matching events.go's event: lines), re-fetches and swaps innerHTML only of sections whose data-rz-events include the event type, re-fetches all after a reconnect, and on stream error toggles is-stale and the stale banner. No browser test was added.

## Links

### Out

- verifies [web-live-sections](../capability/rhiza/web-monitor/web-live-sections.md)

## History

- `01M6PFEXSHWC6TR4F3GFKHV09R` Web server: the daemon serves a read-only web UI on 127.0.0.1:7787 (rhiza.web.port, or off), with read-only internal/api operations as JSON, server-sent events and embedded templates and assets
