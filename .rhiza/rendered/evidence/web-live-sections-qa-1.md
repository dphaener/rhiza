# web-live-sections-qa-1

evidence · `01M6PWYR3ZAF6PYAR6HRH1D0GQ` · actor [claude-code-e4071629](../actor/claude-code-e4071629.md)

## result

PASS

## statement

QA at dfd42a2b: a scratch build (go build ./cmd/rhiza) ran a real long-running rhiza daemon (paused, --max-runs 0) in a fresh scratch repo with the run env unset, and curl hit it. /static/app.js is served, subscribes to /events, swaps sections by data-rz-events and data-rz-fragment, and toggles is-stale and the #rz-stale banner when the stream drops; the rendered page carries the #rz-stale banner. /fragments/nope gave 404. No page yet declares a live section, so no swap was seen in a browser (the leaf asks for no browser test).

## Links

### Out

- verifies [web-live-sections](../capability/rhiza/web-monitor/web-live-sections.md)

## History

- `01M6PFEXSHWC6TR4F3GFKHV09R` Web server: the daemon serves a read-only web UI on 127.0.0.1:7787 (rhiza.web.port, or off), with read-only internal/api operations as JSON, server-sent events and embedded templates and assets
