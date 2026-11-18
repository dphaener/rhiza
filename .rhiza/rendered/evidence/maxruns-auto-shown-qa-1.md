# maxruns-auto-shown-qa-1

evidence · `01M7QZEJEFHT2G6141JPQ9AJ7G` · actor [claude-code-34e635d3](../actor/claude-code-34e635d3.md)

## result

PASS

## statement

QA by hand holding the root daemon lock with flock over a saved state of 2 live runs and 1 starting: max_runs 1 with max_runs_auto true gives daemon status 'Max runs: auto (1 this poll)' and 'runs: 2 live (auto cap 1), 1 starting', --json max_runs_auto true, and health 'ok daemon' with 'runs: 2 live (auto cap 1), 1 starting'; the same state without auto gives 'Max runs: 1', health bad 'over rhiza.daemon.maxRuns 1' naming each live run, and --json has no max_runs_auto; a state with no field and max_runs 3 prints 'Max runs: 3', 'runs: 2/3 live'. TestScripts daemon-status, health-runs-cap, describe-golden pass.

## Links

### Out

- verifies [maxruns-auto-shown](../capability/rhiza/orchestration/daemon/daemon-status/maxruns-auto-shown.md)

## History

- `01M7NVMXZ1YEGZB56ZEAMPX44B` rhiza.daemon.maxRuns=auto lets the daemon start another run only while the pre-start checks (load, memory, free tmp, disk) have room; a number still pins the limit
