# plan-merge-loop-six-tries

capability · `01M7AQFB7ZKZJHBNDRAAS8NVHK` · parent [plan-workflow](../plan-workflow.md)

## acceptance

In internal/workflow/workflow_test.go, TestPlanBuiltin's caps map has rebase, wait-rebased, rebased-green and merge at 6 (wait-pr 3, pr-green 4 and every other cap unchanged) and still checks merge fail routes to rebase, and fails if any of the four is 3; TestBuiltinPlanMilestone still finds plan-milestone's steps equal to plan's; TestPlanRerunsCIFlakes, for both the built-in plan and plan-milestone workflows, wants max_visits 6 for wait-rebased and rebased-green (wait-pr 3 and pr-green 4 unchanged); go test -trimpath ./internal/workflow -run 'TestPlanBuiltin|TestBuiltinPlanMilestone|TestPlanRerunsCIFlakes' passes; cmd/rhiza/testdata/script/plan-workflow-merge-main-moved.txtar still passes as written; a new testscript cmd/rhiza/testdata/script/plan-workflow-merge-main-moved-often.txtar, modelled on it, has the fake gh push a commit to origin's main after each of the first three rebased-green checks: the run's log shows merge (visits 1, 2 and 3) fail then rebase, merge (visit 4) pass then close-out, gh.log has exactly one pr merge call, pinned to the fourth rebased head, no 'went over max_visits' appears, and the run ends succeeded. Files: internal/workflow/builtin/plan.toml, internal/workflow/builtin/plan-milestone.toml (the steps and the 'up to 3 times' comments), internal/workflow/workflow_test.go (TestPlanBuiltin, TestPlanRerunsCIFlakes), cmd/rhiza/testdata/script/plan-workflow-merge-main-moved-often.txtar

## statement

When a plan or plan-milestone run's merge finds origin/main has moved past the PR's checked head, the run goes back to rebase onto the new main (settling rendered-only conflicts by rendering again), waits for the rebased head's CI and merges it, up to 6 times (rebase, wait-rebased, rebased-green and merge each max_visits 6, ship-change's cap), instead of ending the run with 'Not merged'

## Links

### Out

- decided-by [plan-merge-loop-cap-six](../../../../../decision/plan-merge-loop-cap-six.md)
- supersedes plan-run-merge-rebases (retired)

### In

- verifies [plan-merge-loop-six-tries-verified](../../../../../evidence/plan-merge-loop-six-tries-verified.md)
- verifies [plan-merge-loop-six-tries-qa-1](../../../../../evidence/plan-merge-loop-six-tries-qa-1.md)
- verifies [pr-558-fix-watchdog-unmerged-pr-at-01m7aq89397m](../../../../../evidence/pr-558-fix-watchdog-unmerged-pr-at-01m7aq89397m.md)

## History

- `01M7AQ89397MPNWGFSAYG3BK77` Fix watchdog unmerged PR at step rebase: run 01M7AH6RHKMZRAHGJKJFM4VTNT; log: PR green at 6ce8f26 and unmerged after step rebase failed: step "rebase" went over max_visits 3
