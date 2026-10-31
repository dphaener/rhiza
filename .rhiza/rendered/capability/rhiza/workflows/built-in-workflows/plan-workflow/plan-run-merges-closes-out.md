# plan-run-merges-closes-out

capability · `01M6857VP1T62T5SQHWD7R38F6` · parent [plan-workflow](../plan-workflow.md)

## acceptance

plan-workflow-pr.txtar continues: with the fake gh reporting green checks the run merges with --match-head-commit, the change set's state is unchanged by the run, the branch and worktree are gone, and a following ship-change on the change set calls gh pr merge (AlreadyMerged false); a Go test in internal/api/ship_test.go starts from a ship state holding a plan's merge commit and shows open-pr then merge still merge

## statement

A plan run waits for its PR's CI, merges it pinned to the checked head without landing the change set, then closes out the branch and worktree and clears its ship state, so a later ship run on the same change set opens and merges its own PR

## Links

### Out

- decided-by [plan-pr-merges-without-landing](../../../../../decision/plan-pr-merges-without-landing.md)

### In

- verifies [plan-run-merges-closes-out-verified](../../../../../evidence/plan-run-merges-closes-out-verified.md)
- verifies [plan-run-merges-closes-out-qa-1](../../../../../evidence/plan-run-merges-closes-out-qa-1.md)
- verifies [plan-run-merges-closes-out-verified-2](../../../../../evidence/plan-run-merges-closes-out-verified-2.md)

## History

- `01M625RCFY4RB9HC0B69Y08VFV` A plan run ends with its plan on a pushed branch and an open PR, as ship-change does, so no one renders, commits or opens the plan's PR by hand.
