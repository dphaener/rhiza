# maxruns-auto-shown

capability · `01M7QXK34ZQT66CGBGZ180M4EJ` · parent [daemon-status](../daemon-status.md)

## acceptance

go test -trimpath ./cmd/rhiza -run 'TestScripts/(daemon-status|health-runs-cap|describe-golden)': daemon-status gains a case whose saved state has max_runs_auto true and prints 'Max runs: auto (3 this poll)' and '"max_runs_auto": true'; health-runs-cap gains an auto case whose detail names the auto cap; describe-golden regenerated; (failure case) a state without max_runs_auto still prints 'Max runs: 3' and health still goes bad over a numeric cap

## files

internal/api/daemon.go, internal/api/health.go, cmd/rhiza/testdata/script/daemon-status.txtar, cmd/rhiza/testdata/script/health-runs-cap.txtar, cmd/rhiza/testdata/script/describe-golden.txtar

## must_not

Change the web UI or the Rex apps' display; drop or rename an existing --json field

## statement

rhiza daemon status says 'Max runs: auto (N this poll)' and its --json carries max_runs_auto when the last poll ran under auto, and rhiza health's daemon part writes 'runs: L live (auto cap N), S starting' and is never bad for being over the cap under auto

## Links

### Out

- depends-on [maxruns-auto-poll](../daemon-picks/maxruns-auto-poll.md)
- decided-by [maxruns-auto-one-per-poll](../../../../../decision/maxruns-auto-one-per-poll.md)

### In

- verifies [maxruns-auto-shown-verified](../../../../../evidence/maxruns-auto-shown-verified.md)
- verifies [maxruns-auto-shown-qa-1](../../../../../evidence/maxruns-auto-shown-qa-1.md)

## History

- `01M7NVMXZ1YEGZB56ZEAMPX44B` rhiza.daemon.maxRuns=auto lets the daemon start another run only while the pre-start checks (load, memory, free tmp, disk) have room; a number still pins the limit
