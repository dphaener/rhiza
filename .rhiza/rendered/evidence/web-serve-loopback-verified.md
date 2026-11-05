# web-serve-loopback-verified

evidence · `01M6PTDFXVBXX31Z9YPPK5CF32` · actor [claude-code-e101ac60](../actor/claude-code-e101ac60.md)

## result

PASS

## statement

Verified at eb7a1c45 (attempt 2). go test -trimpath ./internal/web/... ./internal/daemon/... pass: TestParsePort accepts 1..65535 and off and rejects 0, 70000, abc; TestAddrIsAlwaysLoopback gives 127.0.0.1:7787 by default, 127.0.0.1:8080 when changed, and an empty address for off; TestHandlerServesReadOnly drives the handler through httptest with no socket bound; TestManagerFollowsSetting covers start, unchanged, move, stop, restart and retry after a failed start. go test -trimpath -race -count=10 ./internal/api -run TestDaemonWeb passed all 10 runs with no race reported. TestDaemonWebListenerFollowsPort now feeds each rhiza.web.port read over a channel and checks after the next tick, so polls see default 7787, then 8080, off, 7787 and 7787 again in order; at the off poll it asserts the 8080 fake is closed and nothing is open (attempt 1's findings f1 and f2 are fixed), and it checks no listener is left open after the interrupt. TestDaemonWebNoListenerOnceOrTestNow shows that --once and a set RHIZA_TEST_NOW open no listener. go test -trimpath ./cmd/rhiza -run 'TestScripts/(daemon-settings|describe-golden)$' passes: rhiza.web.port is listed as 7787 (default), set to 8080 and off is accepted, and 0, 70000 and abc are rejected with stderr naming rhiza.web.port. Read the diff: daemon.go builds the web.Manager only when not --once and RHIZA_TEST_NOW is unset, applies conf.webPort at each poll, and closes the manager on exit.

## Links

### Out

- verifies [web-serve-loopback](../capability/rhiza/web-monitor/web-serve-loopback.md)

## History

- `01M6PFEXSHWC6TR4F3GFKHV09R` Web server: the daemon serves a read-only web UI on 127.0.0.1:7787 (rhiza.web.port, or off), with read-only internal/api operations as JSON, server-sent events and embedded templates and assets
