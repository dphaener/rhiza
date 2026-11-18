# maxruns-auto-shown-verified

evidence · `01M7QYSCWAAQYK9WQCY7FY9BG0` · actor [claude-code-f1c47b5d](../actor/claude-code-f1c47b5d.md)

## result

PASS

## statement

Verified b48a3b310 (diff vs 8eea0ef15). Re-ran every Checks: command, each exit 0: describe-golden -update (no golden diff); go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/(daemon-status|health-runs-cap|describe-golden)$' -count=1 ok; go vet ./internal/api ok; go test -p 2 ./internal/api -run 'Health|DaemonStatus' ok; rhiza ship generate (0 changed); rhiza render (only this verifier's own actor). Acceptance: daemon-status.txtar adds an auto case asserting 'Max runs: auto (3 this poll)', 'runs: N live (auto cap 3)' and '"max_runs_auto": true', and the numeric case still asserts '^Max runs: 3$' and no max_runs_auto; health-runs-cap.txtar adds an auto case (ok, 'runs: 3 live (auto cap 2), 1 starting', no 'over rhiza.daemon.maxRuns') while the existing over-cap case still goes bad; describe-golden carries max_runs_auto. Diff: DaemonStatusResponse.MaxRunsAuto from daemon state, shared runsAgainstCap line, health bad only when !MaxRunsAuto. No JSON field dropped or renamed, no web/Rex app change, no TODOs.

## Links

### Out

- verifies [maxruns-auto-shown](../capability/rhiza/orchestration/daemon/daemon-status/maxruns-auto-shown.md)

## History

- `01M7NVMXZ1YEGZB56ZEAMPX44B` rhiza.daemon.maxRuns=auto lets the daemon start another run only while the pre-start checks (load, memory, free tmp, disk) have room; a number still pins the limit
