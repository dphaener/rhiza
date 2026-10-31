# plan-run-merges-closes-out-verified-2

evidence · `01M6961SM8ETP8QB31VEF354E3` · actor [claude-code-fb939bcf](../actor/claude-code-fb939bcf.md)

## result

PASS

## statement

Verified at 75effb4 (attempt 2, fixing quality f1). Diff: plan.toml now runs pr-green -> rebase (rhiza ship rebase) -> wait-rebased -> rebased-green -> merge -> close-out --forget, as ship-change does; TestPlanBuiltin and plan-workflow-pr.txtar assert the new edges. Ran: go test ./cmd/rhiza -run 'TestScripts/(plan-workflow-pr|plan-workflow)$' (ok): with green fake checks the run calls gh pr merge 7 --squash --match-head-commit pinned to the checked head, the change set stays proposed, branch and worktree are gone locally and on origin, ship state is cleared, and a following ship-change calls gh pr merge 8 (a fresh merge, AlreadyMerged false). go test ./internal/api -run TestShipMergeAfterPlanMerge (PASS: from a ship state holding the plan's merge commit, open-pr then merge still merges). go test ./internal/workflow -run TestPlanBuiltin (ok). make check exit 0.

## Links

### Out

- verifies [plan-run-merges-closes-out](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-run-merges-closes-out.md)

## History

- `01M625RCFY4RB9HC0B69Y08VFV` A plan run ends with its plan on a pushed branch and an open PR, as ship-change does, so no one renders, commits or opens the plan's PR by hand.
