# health-report-verified

evidence · `01M6DFVTASA1J7J3JX31RJJ1KK` · actor [claude-code-ed01284a](../actor/claude-code-ed01284a.md)

## result

PASS

## statement

At 0b903288, with RHIZA_ROOT/CHANGE/RUN unset: `go test ./cmd/rhiza -run 'TestScripts/(health|describe-golden)$' -count=1 -v` passes both. health.txtar covers each acceptance point: no daemon is warn with 'start it with rhiza daemon'; RHIZA_TEST_FREE_BYTES=100 makes disk bad and health exits 1; a dirty main is warn; a failed inbox run is warn ('1 waiting: 1 failed'); a model problem makes check bad; health --json and api call health are compared with cmp (ok and bad cases). describe-golden.txtar holds the health op. Attempt 1's finding is fixed: 0b903288 lands the missing actor records and re-renders, and `make check` now exits 0 with a clean tree afterwards.

## Links

### Out

- verifies [health-report](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-report.md)

## History

- `01M6CRZSRFQ94MAJ8CRTEE44P2` rhiza health and a health app show the whole system's state on one screen
