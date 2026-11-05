# web-api-read-only-qa-1

evidence · `01M6PWYQWWFZENP0VXCNKF4ZP9` · actor [claude-code-e4071629](../actor/claude-code-e4071629.md)

## result

PASS

## statement

QA at dfd42a2b: a scratch build (go build ./cmd/rhiza) ran a real long-running rhiza daemon (paused, --max-runs 0) in a fresh scratch repo with the run env unset, and curl hit it. GET /api/run.ls, /api/roadmap, /api/ls and /api/status were byte-identical (cmp) to rhiza run ls/roadmap/ls/status --json; ?limit=1 was honoured, ?bogus=1 gave 400 with a usage error, HEAD /api/ls 200. GET /api/land, accept, daemon.set, daemon.start, render, ship.merge, add, apply and change.new all gave 404; POST, PUT, DELETE and PATCH on /api/run.ls gave 405.

## Links

### Out

- verifies [web-api-read-only](../capability/rhiza/web-monitor/web-api-read-only.md)

## History

- `01M6PFEXSHWC6TR4F3GFKHV09R` Web server: the daemon serves a read-only web UI on 127.0.0.1:7787 (rhiza.web.port, or off), with read-only internal/api operations as JSON, server-sent events and embedded templates and assets
