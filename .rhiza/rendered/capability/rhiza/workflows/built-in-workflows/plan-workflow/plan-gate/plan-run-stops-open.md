# plan-run-stops-open

capability · `01M6857YNN8RZG1ZHVPBF9YP4J` · parent [plan-gate](../plan-gate.md)

## acceptance

A new plan-workflow-open.txtar with the fake rex and claude of plan-workflow.txtar: when the planner parks a question the run does not end succeeded, no gh pr create is logged, run status and inbox --json name the question; after a Decision answers it, run resume goes on and the run succeeds; a Go test in internal/run covers the stop outcome

## statement

A plan run whose gate finds open questions does not end succeeded and opens no PR: it stops for the owner, run status and rhiza inbox name each open question and the command to resume, and resuming after they are answered goes on to ship the plan

## Links

### Out

- decided-by [plan-stops-before-pr-on-open-items](../../../../../../decision/plan-stops-before-pr-on-open-items.md)

### In

- verifies [plan-run-stops-open-verified](../../../../../../evidence/plan-run-stops-open-verified.md)
- verifies [plan-run-stops-open-qa-1](../../../../../../evidence/plan-run-stops-open-qa-1.md)
- verifies [plan-run-stops-open-verified-3](../../../../../../evidence/plan-run-stops-open-verified-3.md)
- verifies [plan-run-stops-open-verified-5](../../../../../../evidence/plan-run-stops-open-verified-5.md)
- verifies [pr-137-plan-run-open-findings](../../../../../../evidence/pr-137-plan-run-open-findings.md)

## History

- `01M625W0AJZ6AMCCWGFXQW19FN` A plan run that ends with confirmed review findings or parked questions still open does not report success: it stops for the owner with them listed, so nothing has to be fixed by hand after a succeeded run.
