# plan-run-stops-open-verified-3

evidence · `01M69GSBMQ9QWT5T5SBKD816TE` · actor [claude-code-160e8070](../actor/claude-code-160e8070.md)

## result

PASS

## statement

At commit e6adac1a (attempt 3, cycle 2 fix of review findings): go test ./cmd/rhiza -run 'TestScripts/(plan-workflow-open|plan-workflow|run-gate|run-review-check-landed)$' -count=1 passes all four. plan-workflow-open covers a parked question: the run does not succeed, no gh pr create is logged, run status and inbox --json name the question, and after a Decision answers it, run resume goes on to success. run-gate now expects an autonomous gate with an open review question to stop with exit 75 and name it (once in the list, once on the stop line), accepting nothing. go test ./internal/run -run Stop passes (TestCommandStopExitStopsRun). go test ./... passes, and go vet is clean. rhiza check at this commit, built from it with the run's RHIZA_ROOT/RHIZA_CHANGE/RHIZA_RUN unset, reports ok. Run inside the run env, it flags only rendered actor files of this run's live sessions in the shared worktree. Read the diffs: gateReview now reports open review questions rather than erroring, and runRunGate commits the response and stops with nothing accepted. This fixes the earlier verifier finding (run-gate count 1 -> 2).

## Links

### Out

- verifies [plan-run-stops-open](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-gate/plan-run-stops-open.md)

## History

- `01M625W0AJZ6AMCCWGFXQW19FN` A plan run that ends with confirmed review findings or parked questions still open does not report success: it stops for the owner with them listed, so nothing has to be fixed by hand after a succeeded run.
