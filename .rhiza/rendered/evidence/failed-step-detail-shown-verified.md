# failed-step-detail-shown-verified

evidence · `01M69PZ1Z7TAYA694DBV7YHTPB` · actor [claude-code-d023a930](../actor/claude-code-d023a930.md)

## result

PASS

## statement

At dc6e4a22, go test ./cmd/rhiza -run TestScripts/run-command-step-output -count=1 passed: plain rhiza run status matched '^- b \(visit 1\): fail, then end: exit 3$' and an indented '^    boom$' line beneath it. Diff: RunInfo.Human prints the detail's first line on the step line and the rest indented; Rex status (run.go rexRecord path) still cuts detail to its first line, unchanged. make check exited 0.

## Links

### Out

- verifies [failed-step-detail-shown](../capability/rhiza/workflows/workflow-runs/run-status/failed-step-detail-shown.md)

## History

- `01M686RRCJSQ57M82FENWT2HNM` When a run's command step fails, its journal detail keeps the tail of the step's output, shown by run status and the inbox, so an unattended failure says why
