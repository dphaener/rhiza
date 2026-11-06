# web-events-stream-qa-1

evidence · `01M6PWYQZ6QK13CSAV9PVH2JXC` · actor [claude-code-e4071629](../actor/claude-code-e4071629.md)

## result

PASS

## statement

QA at dfd42a2b: a scratch build (go build ./cmd/rhiza) ran a real long-running rhiza daemon (paused, --max-runs 0) in a fresh scratch repo with the run env unset, and curl hit it. curl -N /events answered text/event-stream with a ready comment and keep-alives. rhiza add to an existing change set and rhiza change new each produced event: change with that change set ID; toggling rhiza.daemon.paused produced event: daemon; raising rhiza.daemon.minFree to 1000T raised a disk-low inbox item and an event: inbox, and lowering it back sent event: inbox again as the item left. Run journal and driver log events were not exercised hands-on (no agent run was started); the package test covers them.

## Links

### Out

- verifies [web-events-stream](../capability/rhiza/web-monitor/web-events-stream.md)

## History

- `01M6PFEXSHWC6TR4F3GFKHV09R` Web server: the daemon serves a read-only web UI on 127.0.0.1:7787 (rhiza.web.port, or off), with read-only internal/api operations as JSON, server-sent events and embedded templates and assets
