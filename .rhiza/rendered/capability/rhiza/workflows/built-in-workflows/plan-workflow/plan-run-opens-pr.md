# plan-run-opens-pr

capability · `01M6857TPTS3BZMRX3S51FBSXT` · parent [plan-workflow](../plan-workflow.md)

## acceptance

A new plan-workflow-pr.txtar with a bare origin.git and a fake gh, as run-ship-change.txtar has, asserts gh.log holds one pr create with the change set's title, the body file holds the round and finding lines, origin has the branch head and the ship state records the PR number; a Go test shows ship open-pr still refuses a plan run with no review round

## statement

A plan run pushes its branch and opens a pull request titled by the change set, its body a plan summary of the review rounds, confirmed and dropped findings, questions and the gate's moves; a second visit pushes and replaces the body

## Links

### Out

- decided-by [plan-pr-merges-without-landing](../../../../../decision/plan-pr-merges-without-landing.md)

### In

- verifies [plan-run-opens-pr-verified](../../../../../evidence/plan-run-opens-pr-verified.md)
- verifies [plan-run-opens-pr-qa-1](../../../../../evidence/plan-run-opens-pr-qa-1.md)

## History

- `01M625RCFY4RB9HC0B69Y08VFV` A plan run ends with its plan on a pushed branch and an open PR, as ship-change does, so no one renders, commits or opens the plan's PR by hand.
