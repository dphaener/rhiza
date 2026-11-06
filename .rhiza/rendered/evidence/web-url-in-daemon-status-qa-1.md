# web-url-in-daemon-status-qa-1

evidence · `01M6PWYR2KTTYTY02Q26XKBKHB` · actor [claude-code-e4071629](../actor/claude-code-e4071629.md)

## result

PASS

## statement

QA at dfd42a2b: a scratch build (go build ./cmd/rhiza) ran a real long-running rhiza daemon (paused, --max-runs 0) in a fresh scratch repo with the run env unset, and curl hit it. After rhiza daemon --once with the setting unset, rhiza daemon status showed Web UI: http://127.0.0.1:7787 and --json web http://127.0.0.1:7787; after rhiza daemon set rhiza.web.port off and another --once it showed Web UI: off and web off. A running daemon status followed the 18787, 18788 and off moves.

## Links

### Out

- verifies [web-url-in-daemon-status](../capability/rhiza/web-monitor/web-url-in-daemon-status.md)

## History

- `01M6PFEXSHWC6TR4F3GFKHV09R` Web server: the daemon serves a read-only web UI on 127.0.0.1:7787 (rhiza.web.port, or off), with read-only internal/api operations as JSON, server-sent events and embedded templates and assets
