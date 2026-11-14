# plan-merge-loop-six-tries-qa-1

evidence · `01M7AY43FAA2YJ0JZE45RJW2FB` · actor [claude-code-e33614b1](../actor/claude-code-e33614b1.md)

## result

PASS

## statement

QA at e0a5e9bcf: built ./cmd/rhiza to a scratch path and ran it by hand on plan-workflow-merge-main-moved-often.txtar's fakes, outside testscript (own HOME and git config, a fake gh that moves origin/main). With main moving after the first three rebased-green checks, rhiza run start --workflow plan: merge failed at visits 1-3 ('Not merged: origin/main has moved past PR #7's checked head'), each followed by rebase; merge passed at visit 4, then close-out; exactly one gh pr merge; the run succeeded. With main moving after every check (failure case), plan and plan-milestone each made six rebase/wait-rebased/rebased-green/merge rounds, then failed cleanly with 'step "rebase" went over max_visits 6' and no gh pr merge. rhiza workflow check plan and plan-milestone: valid. go test -trimpath ./internal/workflow -run 'TestPlanBuiltin|TestBuiltinPlanMilestone|TestPlanRerunsCIFlakes': ok. TestScripts/plan-workflow-merge-main-moved and -often: PASS.

## Links

### Out

- verifies [plan-merge-loop-six-tries](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-merge-loop-six-tries.md)

## History

- `01M7AQ89397MPNWGFSAYG3BK77` Fix watchdog unmerged PR at step rebase: run 01M7AH6RHKMZRAHGJKJFM4VTNT; log: PR green at 6ce8f26 and unmerged after step rebase failed: step "rebase" went over max_visits 3
