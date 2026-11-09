# daemon-status-runs-cap-qa-1

evidence · `01M6YEGHTY2163DNDJW7S7F3M9` · actor [claude-code-53686c59](../actor/claude-code-53686c59.md)

## result

PASS

## statement

QA at 3686b6ce, with the built binary in a scratch root. rhiza daemon status printed 'runs: 1/3 live, 2 starting' (and '1/2 live, 1 starting' or '0/3 live, 3 starting' in the resume scenarios), and daemon status --json gave live, starting and max matching it. The web monitor came from a long-lived paused daemon with rhiza.web.port 17991, no RHIZA_TEST_NOW and nothing startable; GET /p/repo/ rendered 'Run slots 1 / 3' and 'runs: 1/3 live, 0 starting', which matched daemon status. I did not run the Daemon app under Rex. Its runs line comes from daemonapp.go:148, the same format, and go test -trimpath ./internal/tui/daemonapp/... passed.

## Links

### Out

- verifies [daemon-status-runs-cap](../capability/rhiza/orchestration/daemon/daemon-status/daemon-status-runs-cap.md)

## History

- `01M6XV4WGSPC346CQEV9DGYWR3` Fix Request plan runs start outside rhiza.daemon.maxRuns, so the machine runs far more agents than its cap: on 2026-11-08 ~13:05Z, with maxRuns 4 (then 3), the root daemon started the plan runs of requests 01M6XHDV2P65 and 01M6XKT75PHW in the same burst as ship runs of 01M6XHCX4FRE and 01M6WGEPSYWV, leaving 5 running and 1 interrupted (daemon status) and load at 23 on the 4-core VM. Earlier the same day 5 ran at cap 3. Plan runs should count against the machine-wide cap (or a small separate plan cap in rhiza.toml, as policy): a waiting request may reserve the next free slot (#226) but must not exceed the cap. rhiza health should show live runs against the cap.
