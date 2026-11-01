# run-status-open-cycle-qa-1

evidence · `01M6B5CPB0Y2VGZ2QNT8TVPVHM` · actor [claude-code-46ead0c1](../actor/claude-code-46ead0c1.md)

## result

PASS

## statement

QA at d6d0e740 with a freshly built binary in a scratch git repo, hand-written journals. OPEN1 (alpha cycle 1 failed, cycle 2 open; beta cycle 1 passed): run status prints '- alpha (step review): open in cycle 2 (cycle 1 failed)' with no 'after N cycles'; --json alpha has no outcome, open_cycle 2, prior_cycle 1, prior_outcome failed, its cycle 1 keeps outcome failed; run summary heads '### alpha: open' and '### beta: passed'; the live app (rhiza tui app live, driven in a pty) shows 'checks 1 passed, 1 running'. DONE1 (same plus alpha cycle 2 passed): status 'passed after 2 cycles', JSON outcome passed with no open_cycle, summary '### alpha: passed', live app 'checks 2 passed'. go test of internal/run, internal/api and internal/tui/liveapp passes.

## Links

### Out

- verifies [run-status-open-cycle](../capability/rhiza/workflows/workflow-runs/run-status/run-status-open-cycle.md)

## History

- `01M6ATASW7ZS2DPMJM2Q12GHM3` A review check whose latest cycle is open reads as open in run status, the live app and run summary
