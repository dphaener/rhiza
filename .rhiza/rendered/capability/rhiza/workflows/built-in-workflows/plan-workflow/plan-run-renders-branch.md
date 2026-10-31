# plan-run-renders-branch

capability · `01M6857SSEWGCV5WVPMG1VCH6A` · parent [plan-workflow](../plan-workflow.md)

## acceptance

plan-workflow.txtar, changed from its old expectation of changes only under .rhiza/, asserts that after the run git status --porcelain in the worktree is empty, rhiza check passes there, and the last commit is the render commit carrying Change: and Actor: trailers

## statement

After its gate a plan run renders and commits .rhiza/rendered and every .rhiza/log entry the run wrote on its branch, with Change: and Actor: trailers, so the branch passes rhiza check; a run with nothing to render adds no commit

## Links

### Out

- decided-by [plan-pr-merges-without-landing](../../../../../decision/plan-pr-merges-without-landing.md)

### In

- verifies [plan-run-renders-branch-verified](../../../../../evidence/plan-run-renders-branch-verified.md)
- verifies [plan-run-renders-branch-qa-1](../../../../../evidence/plan-run-renders-branch-qa-1.md)

## History

- `01M625RCFY4RB9HC0B69Y08VFV` A plan run ends with its plan on a pushed branch and an open PR, as ship-change does, so no one renders, commits or opens the plan's PR by hand.
