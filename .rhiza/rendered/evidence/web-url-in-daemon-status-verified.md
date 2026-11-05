# web-url-in-daemon-status-verified

evidence · `01M6PWAY0NX0RSJYBD2EXW9J5K` · actor [claude-code-fee64b7b](../actor/claude-code-fee64b7b.md)

## result

PASS

## statement

At 7ef54661, go test -trimpath ./cmd/rhiza -run 'TestScripts/(daemon-status|describe-golden)$' -count=1 passes: daemon-status.txtar runs rhiza daemon --once, then rhiza daemon status prints 'Web UI: http://127.0.0.1:7787' and --json has "web": "http://127.0.0.1:7787"; after rhiza daemon set rhiza.web.port off and another rhiza daemon --once, status prints 'Web UI: off' and --json "web": "off". The diff (internal/api/daemon.go webStatus) derives the line from the state file's recorded settings via web.Addr on rhiza.web.port and never dials the listener. go test -trimpath -short ./internal/api/ ./internal/web/ pass; gofmt -l and go vet ./internal/api are clean. The unshortened internal/api run timed out in TestReadBudget, which was blocked waiting for the machine-wide budget lock and a quiet machine while other worktrees ran tests (not this diff).

## Links

### Out

- verifies [web-url-in-daemon-status](../capability/rhiza/web-monitor/web-url-in-daemon-status.md)

## History

- `01M6PFEXSHWC6TR4F3GFKHV09R` Web server: the daemon serves a read-only web UI on 127.0.0.1:7787 (rhiza.web.port, or off), with read-only internal/api operations as JSON, server-sent events and embedded templates and assets
