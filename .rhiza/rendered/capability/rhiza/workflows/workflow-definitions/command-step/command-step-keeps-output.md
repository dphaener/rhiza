# command-step-keeps-output

capability · `01M686TD7X0X91G79FG5S2Q1GJ` · parent [command-step](../command-step.md)

## acceptance

A new run-command-step-output.txtar runs a step echo out1; echo boom >&2; exit 3 and run status --json's detail holds exit 3, boom and out1; a step printing 500 lines keeps only the tail (line 500, not line 1); a passing step's detail stays exit 0; a unit test in internal/run/engine_test.go shows a runner's detail is kept

## statement

When a command step exits non-zero, its journal detail keeps exit N followed by the last lines of the step's interleaved stdout and stderr, capped at 20 lines and 2 KiB

## Links

### Out

- decided-by [step-output-tail-cap](../../../../../decision/step-output-tail-cap.md)

### In

- verifies [command-step-keeps-output-verified](../../../../../evidence/command-step-keeps-output-verified.md)
- verifies [command-step-keeps-output-qa-1](../../../../../evidence/command-step-keeps-output-qa-1.md)
- verifies [pr-140-failed-steps-keep-output](../../../../../evidence/pr-140-failed-steps-keep-output.md)

## History

- `01M686RRCJSQ57M82FENWT2HNM` When a run's command step fails, its journal detail keeps the tail of the step's output, shown by run status and the inbox, so an unattended failure says why
