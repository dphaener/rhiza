# plan-run-stops-open-verified-5

evidence · `01M69JD1ZB8T4518WERMN1RXFC` · actor [claude-code-e21a6a6e](../actor/claude-code-e21a6a6e.md)

## result

PASS

## statement

At 47b72737 (attempt 5) I ran go test ./cmd/rhiza -run 'TestScripts/(plan-workflow-open|plan-workflow|run-gate)$' -count=1: PASS. plan-workflow-open.txtar uses the fake rex/claude/gh: the planner parks park-disk, the run output lacks succeeded, gh.log is absent (no pr create), run status --json shows state stopped at gate with 'open questions: park-disk' and 'rhiza run resume', inbox --json has kind open with the question and next 'rhiza run resume'; after a Decision answers park-disk, run resume exits 0, succeeds and logs one pr create. go test ./internal/run -run TestCommandStopExitStopsRun: PASS (Go test of the stop outcome). The retry fix (questionOpen shared by gateReview and gateOpenQuestions, so a retired review question is closed for the stop and for resp.Open) is covered by the new run-gate.txtar case, which passes. go vet ./... clean, gofmt -l empty, go test ./... passes (modelapp's timing test failed once under full-suite load at 61ms vs 50ms and passed when rerun alone).

## Links

### Out

- verifies [plan-run-stops-open](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-gate/plan-run-stops-open.md)

## History

- `01M625W0AJZ6AMCCWGFXQW19FN` A plan run that ends with confirmed review findings or parked questions still open does not report success: it stops for the owner with them listed, so nothing has to be fixed by hand after a succeeded run.
