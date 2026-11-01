# ship-check-reruns-failed-once-verified

evidence · `01M6ANNF8N73405PVC2YSZ1TVC` · actor [claude-code-778edc55](../actor/claude-code-778edc55.md)

## result

PASS

## statement

At d804f8a2 ran go test ./cmd/rhiza -run TestScripts/(run-ship-rerun-pass|run-ship-rerun-fail|run-ship-change|describe-golden): all PASS. run-ship-rerun-pass uses a fake gh whose pr view reports a FAILURE check with detailsUrl .../actions/runs/42/job/9 until gh run rerun is called; it asserts gh.log holds exactly one 'run rerun 42 --failed' and run status shows pr-green pass. run-ship-rerun-fail keeps the check failing after the rerun; run start fails with step pr-green outcome fail and exactly one rerun in gh.log. go test ./internal/api -run TestShipCheck(MainRerunsFailedOnce|NoRunIDNotRerun) PASS: --main reads databaseId 99 from run list, reruns once, records rerun_head/reruns, a second failure on the same head fails with no second rerun, and a check without an Actions run id is not rerun. Read the diff of internal/api/ship.go (actionsRunID, parseRunList databaseId, shipRerunFailed recording before rerunning). make check exit 0.

## Links

### Out

- verifies [ship-check-reruns-failed-once](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-reruns-failed-once.md)

## History

- `01M6A0JQ2WHGH7DDFKJGRXJD85` A failed CI check on a run's PR or main is rerun once per head before the run fails
