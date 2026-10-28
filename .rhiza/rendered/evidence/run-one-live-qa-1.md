# run-one-live-qa-1

evidence · `01M624NTWVMCCY7Q7RJTPBJTHW` · actor [claude-code-7e6e3e8f](../actor/claude-code-7e6e3e8f.md)

## result

PASS

## statement

Built rhiza from the branch head and drove it in a scratch git repo with real hold (wait) and slow/flaky (command) workflows. With a run stopped at a wait via --no-wait, a second run start (human, --json and --dry-run) exited 2 naming that run with next: rhiza run resume <id>, and no second journal appeared. A run whose process was SIGKILLed mid-step still blocked a new start, exit 2 naming it. Six concurrent run start processes on one change set: exactly one began (exit 0), five exited 2 naming it. With a failed run and a second run stopped at a wait, run resume --retry (and --dry-run) on the failed run exited 2 naming the second. Two concurrent resume --retry of two failed runs: one succeeded, the other exited 2 naming it. After the live run finished, a new start succeeded.

## Links

### Out

- verifies [run-one-live](../capability/rhiza/workflows/workflow-runs/run-one-live.md)

## History

- `01M5ZNV4E4AJR52EX631180DMS` A change set has at most one live run: rhiza run start refuses a second run on a change set whose run is still going, naming that run.
