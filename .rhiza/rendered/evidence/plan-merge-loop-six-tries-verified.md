# plan-merge-loop-six-tries-verified

evidence · `01M7AXRR6HAEJVY6AY7GJSZ3VP` · actor [claude-code-e90c5199](../actor/claude-code-e90c5199.md)

## result

PASS

## statement

At 1586c178f re-ran the Checks block: go test -trimpath ./internal/workflow -run 'TestPlanBuiltin|TestBuiltinPlanMilestone|TestPlanRerunsCIFlakes' exit 0; go test -trimpath ./cmd/rhiza -run TestScripts/plan-workflow-merge-main-moved exit 0 (both plan-workflow-merge-main-moved and the new plan-workflow-merge-main-moved-often pass); rhiza ship generate exit 0 (0 changed); rhiza render exit 0, tree clean. Also passed the whole internal/workflow package and the related scripts plan-workflow-pr, plan-workflow-ci-flake, plan-workflow-ci-infra-wait, plan-workflow-ci-infra-timeout and workflow-check. Diff read: plan.toml and plan-milestone.toml set rebase, wait-rebased, rebased-green and merge to max_visits 6 with the comment now 'up to 6 times'; TestPlanBuiltin caps have those four at 6 (wait-pr 3, pr-green 4, others unchanged) and still check merge fail routes to rebase; TestPlanRerunsCIFlakes wants 6 for wait-rebased and rebased-green in both workflows; TestBuiltinPlanMilestone is unchanged and still compares the steps. The new script makes the fake gh move origin's main after each of the first three rebased-green checks and asserts merge visits 1-3 fail then rebase, visit 4 passes then close-out, no visit 5, no 'went over max_visits', exactly one pr merge pinned to the merged head, run succeeded. Only fakes are used, and there are no stubs or TODOs and no changes outside the leaf apart from the builder's actor record.

## Links

### Out

- verifies [plan-merge-loop-six-tries](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-merge-loop-six-tries.md)

## History

- `01M7AQ89397MPNWGFSAYG3BK77` Fix watchdog unmerged PR at step rebase: run 01M7AH6RHKMZRAHGJKJFM4VTNT; log: PR green at 6ce8f26 and unmerged after step rebase failed: step "rebase" went over max_visits 3
