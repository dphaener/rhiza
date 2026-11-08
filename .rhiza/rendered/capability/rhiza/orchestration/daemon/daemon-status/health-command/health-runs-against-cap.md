# health-runs-against-cap

capability · `01M6XVG9PE300YRA4FA4A530WC` · parent [health-command](../health-command.md)

## acceptance

cmd/rhiza/testdata/script/health.txtar (or a new health-runs-cap.txtar) with a fake daemon state: with max-runs 4, three live runs and one unjournaled start, health shows 'runs: 3/4 live, 1 starting' and the daemon part ok; with max-runs 2 and three live runs it shows the part bad, names the three runs and health exits 1; health --json carries live, starting and max; go test -trimpath ./cmd/rhiza -run 'TestScripts/(health|health-projects|health-runs-cap)' passes

## statement

rhiza health's daemon part shows the runs against the cap, summed across the daemon's projects, as 'runs: <live>/<maxRuns> live, <n> starting' (starts not yet journaled within daemonStartGrace shown apart); the part is ok while live plus starting is at or under rhiza.daemon.maxRuns and bad when it exceeds it, with a detail line naming each live run (project, run, change set title)

## Links

### Out

- decided-by [plan-runs-share-max-runs](../../../../../../decision/plan-runs-share-max-runs.md)

### In

- verifies [health-runs-against-cap-verified](../../../../../../evidence/health-runs-against-cap-verified.md)
- verifies [health-runs-against-cap-qa-1](../../../../../../evidence/health-runs-against-cap-qa-1.md)
- verifies [pr-398-fix-request-plan-runs-start](../../../../../../evidence/pr-398-fix-request-plan-runs-start.md)

## History

- `01M6XV4WGSPC346CQEV9DGYWR3` Fix Request plan runs start outside rhiza.daemon.maxRuns, so the machine runs far more agents than its cap: on 2026-11-08 ~13:05Z, with maxRuns 4 (then 3), the root daemon started the plan runs of requests 01M6XHDV2P65 and 01M6XKT75PHW in the same burst as ship runs of 01M6XHCX4FRE and 01M6WGEPSYWV, leaving 5 running and 1 interrupted (daemon status) and load at 23 on the 4-core VM. Earlier the same day 5 ran at cap 3. Plan runs should count against the machine-wide cap (or a small separate plan cap in rhiza.toml, as policy): a waiting request may reserve the next free slot (#226) but must not exceed the cap. rhiza health should show live runs against the cap.
