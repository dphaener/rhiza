# plan-pr-merges-without-landing

decision · `01M68555Q989EZZN42FPMGG5S5`

## statement

A plan run ships its plan as ship-change ships code, minus landing: it renders and commits, opens the PR, waits for green CI, merges it and closes out the branch and worktree; it never lands the change set, whose acceptance stays the owner's gate

## wording

(assumed under the owner's standing instruction to decide on their behalf)

## Links

### Out

- derived-from [m6-overnight-interview](../source/m6-overnight-interview.md)

### In

- decided-by [plan-run-renders-branch](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-run-renders-branch.md)
- decided-by [plan-run-opens-pr](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-run-opens-pr.md)
- decided-by [plan-run-merges-closes-out](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-run-merges-closes-out.md)

## History

- `01M684XWDT3EKJD4JPXT0JVK2M` Plan milestone 6: Rhiza runs unattended overnight, with the run machinery fixes that needed the orchestrator's hand in milestones 4 and 5
