# plan-gate-open-items-qa-1

evidence · `01M69EAN98KZ7XYPY3KNAGZQVK` · actor [claude-code-650867ac](../actor/claude-code-650867ac.md)

## result

PASS

## statement

QA at 29339ad6 with the built binary in a scratch repo (hold workflow, no fakes): a record change set with a three-round confirmed finding left question review-01m69ea7-r3-f1; rhiza run gate printed '1 question(s) the plan left are open: - review-01m69ea7-r3-f1', a stop: line and next: rhiza run resume, exit 75. After a --record Decision answered it, run gate --json gave "open": [] and exit 0. In a plan run (below) the gate also listed the planner-parked question park-disk, and run gate --dry-run --json in the run worktree returned open [{handle park-disk}] with exit 0.

## Links

### Out

- verifies [plan-gate-open-items](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-gate/plan-gate-open-items.md)

## History

- `01M625W0AJZ6AMCCWGFXQW19FN` A plan run that ends with confirmed review findings or parked questions still open does not report success: it stops for the owner with them listed, so nothing has to be fixed by hand after a succeeded run.
