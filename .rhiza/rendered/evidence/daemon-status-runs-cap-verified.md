# daemon-status-runs-cap-verified

evidence · `01M6YE0CYAW3ABDM3W9MX5JATN` · actor [claude-code-36e17af1](../actor/claude-code-36e17af1.md)

## result

PASS

## statement

At 243aac07 I read the diff. DaemonStatusResponse gains live, starting (always present now) and max. daemonRunsHeld sums non-interrupted runs and unjournaled starts over the status's runs and each project's part. Human prints 'runs: <live>/<max> live, <n> starting'. The Daemon app subtitle and the web daemon panel show the same line, and web Slots = live + starting. Health takes the summed counts instead of adding each project's starts a second time. daemon-status.txtar asserts the human line '^runs: [0-9]+/3 live, [0-9]+ starting$' and the JSON fields live, starting and max=3. daemonapp TestShowsRunsAgainstCap (a fake status with 2 live, 1 starting, max 4) gives 'runs: 2/4 live, 1 starting'. Web TestOverviewReadsDaemonFromRepo (1 running, 1 interrupted, Starting 1, max 4) renders 'runs: 2/4 live, 1 starting' and 3 of 4 slots. With RHIZA_ROOT unset and GOTMPDIR outside the repo, these passed: go test -count=1 -trimpath ./internal/tui/daemonapp/... ./internal/web/...; go test -count=1 -trimpath ./cmd/rhiza -run 'TestScripts/(daemon-status|daemon-status-projects|describe-golden|health.*)$' (daemon-status, daemon-status-projects, describe-golden, health, health-runs-cap, health-projects and health-machine-part all PASS); go test -trimpath ./internal/api -run 'Health|DaemonStatus|WebDaemon'.

## Links

### Out

- verifies [daemon-status-runs-cap](../capability/rhiza/orchestration/daemon/daemon-status/daemon-status-runs-cap.md)

## History

- `01M6XV4WGSPC346CQEV9DGYWR3` Fix Request plan runs start outside rhiza.daemon.maxRuns, so the machine runs far more agents than its cap: on 2026-11-08 ~13:05Z, with maxRuns 4 (then 3), the root daemon started the plan runs of requests 01M6XHDV2P65 and 01M6XKT75PHW in the same burst as ship runs of 01M6XHCX4FRE and 01M6WGEPSYWV, leaving 5 running and 1 interrupted (daemon status) and load at 23 on the 4-core VM. Earlier the same day 5 ran at cap 3. Plan runs should count against the machine-wide cap (or a small separate plan cap in rhiza.toml, as policy): a waiting request may reserve the next free slot (#226) but must not exceed the cap. rhiza health should show live runs against the cap.
