# inbox-failed-step-output

capability · `01M686TFM9P7BREP2N91XSK6AG` · parent [inbox-list](../inbox-list.md)

## acceptance

inbox-list.txtar or run-abandon-failed.txtar gains a run failing at a command step that prints boom; rhiza inbox shows boom indented under the item and inbox --json's detail holds it

## statement

The inbox item of a run that failed at a command step shows that step's kept output

## Links

### Out

- decided-by [step-output-tail-cap](../../../../../decision/step-output-tail-cap.md)

### In

- verifies [inbox-failed-step-output-verified](../../../../../evidence/inbox-failed-step-output-verified.md)
- verifies [inbox-failed-step-output-qa-1](../../../../../evidence/inbox-failed-step-output-qa-1.md)
- verifies [pr-140-failed-steps-keep-output](../../../../../evidence/pr-140-failed-steps-keep-output.md)

## History

- `01M686RRCJSQ57M82FENWT2HNM` When a run's command step fails, its journal detail keeps the tail of the step's output, shown by run status and the inbox, so an unattended failure says why
