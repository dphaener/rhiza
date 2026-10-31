# plan-gate-open-items

capability · `01M6857XR66P2FBADYVDRMMTWX` · parent [plan-gate](../plan-gate.md)

## acceptance

run-gate.txtar adds a planner-written question in a draft and an unanswered review question; run gate --json names both under open, and once a Decision answers each the list is empty; describe-golden.txtar is regenerated

## statement

rhiza run gate lists every unanswered question the plan left, the review's questions and any question node created in the plan's change sets with no live answering Decision, in an open field of its response

## Links

### Out

- decided-by [plan-stops-before-pr-on-open-items](../../../../../../decision/plan-stops-before-pr-on-open-items.md)

### In

- verifies [plan-gate-open-items-verified](../../../../../../evidence/plan-gate-open-items-verified.md)
- verifies [plan-gate-open-items-qa-1](../../../../../../evidence/plan-gate-open-items-qa-1.md)
- verifies [pr-137-plan-run-open-findings](../../../../../../evidence/pr-137-plan-run-open-findings.md)

## History

- `01M625W0AJZ6AMCCWGFXQW19FN` A plan run that ends with confirmed review findings or parked questions still open does not report success: it stops for the owner with them listed, so nothing has to be fixed by hand after a succeeded run.
