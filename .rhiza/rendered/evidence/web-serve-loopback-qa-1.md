# web-serve-loopback-qa-1

evidence · `01M6PWYQV2J46MFZT128F0E1RB` · actor [claude-code-e4071629](../actor/claude-code-e4071629.md)

## result

PASS

## statement

QA at dfd42a2b: a scratch build (go build ./cmd/rhiza) ran a real long-running rhiza daemon (paused, --max-runs 0) in a fresh scratch repo with the run env unset, and curl hit it. With rhiza.web.port unset the daemon listened on 127.0.0.1:7787 only (ss) and GET / gave 200. rhiza daemon set rhiza.web.port rejected 0, 70000 and abc (must be a port number from 1 to 65535, or off) and accepted 18787. With the daemon running, setting 18788, then off, then 18787 moved, stopped and restarted the listener within 4-5 s each (ss shows only 127.0.0.1:<port>). A port already held by another process was logged (bind: address already in use), the daemon stayed up and bound it once freed. strace -e bind shows rhiza daemon --once and a long-running daemon with RHIZA_TEST_NOW set make no bind call.

## Links

### Out

- verifies [web-serve-loopback](../capability/rhiza/web-monitor/web-serve-loopback.md)

## History

- `01M6PFEXSHWC6TR4F3GFKHV09R` Web server: the daemon serves a read-only web UI on 127.0.0.1:7787 (rhiza.web.port, or off), with read-only internal/api operations as JSON, server-sent events and embedded templates and assets
