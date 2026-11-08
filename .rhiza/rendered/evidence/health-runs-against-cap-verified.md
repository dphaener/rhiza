# health-runs-against-cap-verified

evidence · `01M6Y1HQGXA1RPHNTNJYJBHRRX` · actor [claude-code-e530f80a](../actor/claude-code-e530f80a.md)

## result

PASS

## statement

Verifier claude-code-e530f80a at 84bda078: go test -trimpath ./cmd/rhiza -run 'TestScripts/(health|health-projects|health-runs-cap|describe-golden|watchdog-health|daemon-version-in-health|daemon-status.*)$' -count=1 passes, with GOTMPDIR moved outside ~/code/rhiza (the run's GOTMPDIR sits inside the main repo, so every script's rhiza init otherwise finds its .rhiza). The new health-runs-cap.txtar uses two projects' fake daemon state. With max_runs 4, two live runs in alpha, one in beta and starting 1, health prints 'ok daemon running' and 'runs: 3/4 live, 1 starting', and health --json carries live 3, starting 1, max 4. With max_runs 2 it prints 'bad daemon running', 'runs: 3/2 live, 1 starting' and 'live run: alpha, RUNA1' / 'alpha, RUNA2' / 'beta, RUNB1', and exits 1. Diff read: daemon.State gains Starting, set at saveState from pending starts, request starts and resumes within daemonStartGrace; healthRunsAgainstCap sums live (an interrupted run takes no slot) and starting across projects and is bad when live+starting > maxRuns, adding one detail line per live run (project, run, change set title, which falls back to the ID). go vet is clean. go test ./internal/api -run 'Health|Daemon|Start': every test PASS. The package was then marked FAIL only by the enclosing-root guard on ~/code/.rhiza-root/usage.live.json, which the live root daemon rewrites (mtime changed with no test running); the leaf touches no usage code.

## Links

### Out

- verifies [health-runs-against-cap](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-runs-against-cap.md)

## History

- `01M6XV4WGSPC346CQEV9DGYWR3` Fix Request plan runs start outside rhiza.daemon.maxRuns, so the machine runs far more agents than its cap: on 2026-11-08 ~13:05Z, with maxRuns 4 (then 3), the root daemon started the plan runs of requests 01M6XHDV2P65 and 01M6XKT75PHW in the same burst as ship runs of 01M6XHCX4FRE and 01M6WGEPSYWV, leaving 5 running and 1 interrupted (daemon status) and load at 23 on the 4-core VM. Earlier the same day 5 ran at cap 3. Plan runs should count against the machine-wide cap (or a small separate plan cap in rhiza.toml, as policy): a waiting request may reserve the next free slot (#226) but must not exceed the cap. rhiza health should show live runs against the cap.
