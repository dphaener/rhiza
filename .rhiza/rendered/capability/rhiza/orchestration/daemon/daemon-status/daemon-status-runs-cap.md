# daemon-status-runs-cap

capability · `01M6XVG9R7SPC15Y2DKK5M3JEZ` · parent [daemon-status](../daemon-status.md)

## acceptance

cmd/rhiza/testdata/script/daemon-status.txtar shows 'runs: <live>/<max> live, <n> starting' and the --json fields live, starting and max; internal/tui/daemonapp and internal/web unit tests render the same count from a daemon status with live runs and one starting; go test -trimpath ./internal/tui/daemonapp/... ./internal/web/... and go test -trimpath ./cmd/rhiza -run 'TestScripts/(daemon-status|daemon-status-projects)' pass

## statement

rhiza daemon status (human and --json), the Daemon app and the web monitor's daemon overview show the machine-wide runs against the cap in the same form as health: live runs of rhiza.daemon.maxRuns and the starts not yet journaled

## Links

### Out

- decided-by [plan-runs-share-max-runs](../../../../../decision/plan-runs-share-max-runs.md)

### In

- verifies [daemon-status-runs-cap-verified](../../../../../evidence/daemon-status-runs-cap-verified.md)
- verifies [daemon-status-runs-cap-qa-1](../../../../../evidence/daemon-status-runs-cap-qa-1.md)
- verifies [pr-398-fix-request-plan-runs-start](../../../../../evidence/pr-398-fix-request-plan-runs-start.md)

## History

- `01M6XV4WGSPC346CQEV9DGYWR3` Fix Request plan runs start outside rhiza.daemon.maxRuns, so the machine runs far more agents than its cap: on 2026-11-08 ~13:05Z, with maxRuns 4 (then 3), the root daemon started the plan runs of requests 01M6XHDV2P65 and 01M6XKT75PHW in the same burst as ship runs of 01M6XHCX4FRE and 01M6WGEPSYWV, leaving 5 running and 1 interrupted (daemon status) and load at 23 on the 4-core VM. Earlier the same day 5 ran at cap 3. Plan runs should count against the machine-wide cap (or a small separate plan cap in rhiza.toml, as policy): a waiting request may reserve the next free slot (#226) but must not exceed the cap. rhiza health should show live runs against the cap.
