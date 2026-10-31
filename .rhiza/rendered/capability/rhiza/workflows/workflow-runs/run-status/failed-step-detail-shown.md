# failed-step-detail-shown

capability · `01M686TEEXFC829G7CJ4EA044C` · parent [run-status](../run-status.md)

## acceptance

In run-command-step-output.txtar, run status prints - b (visit 1): fail, then end: exit 3 and an indented boom line under it

## statement

rhiza run status shows a failed step's kept output as indented lines under the step's line, and Rex's status line still takes only the first line

## Links

### Out

- decided-by [step-output-tail-cap](../../../../../decision/step-output-tail-cap.md)

### In

- verifies [failed-step-detail-shown-verified](../../../../../evidence/failed-step-detail-shown-verified.md)
- verifies [failed-step-detail-shown-qa-1](../../../../../evidence/failed-step-detail-shown-qa-1.md)
- verifies [pr-140-failed-steps-keep-output](../../../../../evidence/pr-140-failed-steps-keep-output.md)

## History

- `01M686RRCJSQ57M82FENWT2HNM` When a run's command step fails, its journal detail keeps the tail of the step's output, shown by run status and the inbox, so an unattended failure says why
