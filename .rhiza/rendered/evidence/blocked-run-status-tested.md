# blocked-run-status-tested

evidence · `01M5WGZJ1DK8JK0WD8EF7QSMAQ`

## command

go test ./cmd/rhiza -run TestScripts/run-blocked

## result

pass

## statement

A testscript with a blocked fake session shows the tab and the question in run status text and --json, and counts one OSC 7501 blocked report.

## tests

cmd/rhiza/testdata/script/run-blocked.txtar

## Links

### Out

- verifies [blocked-run-status](../capability/rhiza/orchestration/notifications/blocked-run-status.md)

## History

- `01M5SZJ46HYE2DDJ0STCDZE1WP` Blocked runs tell the owner: a run waiting on the owner sends a Rex notification, and genesis run status names the tab and the question.
