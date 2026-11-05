# plan-step-milestone-check-verified

evidence · `01M6MG04MZQJWW6CQF1HQ2MXV3` · actor [claude-code-3da096ca](../actor/claude-code-3da096ca.md)

## result

PASS

## statement

At bb492482 I read the diff (plan.toml: interview on.pass -> milestone-check, a command step running `rhiza run milestone-check` with on.pass review and no on.fail; internal/api/run_milestone.go fails naming each plan change set serving no milestone directly or via an intent under one, with a Fix line) and ran: go test -trimpath -count=1 ./internal/api -run TestRunMilestoneCheck (PASS: fails serving nothing and serving an intent under no milestone, passes through an intent under a milestone and serving a milestone directly); go test ./internal/workflow/... (PASS, TestPlanBuiltin covers milestone-check pass->review and fail->EndFailed); go test ./cmd/rhiza -run TestScripts/ plan-milestone-check, plan-workflow, plan-workflow-pr, plan-workflow-autonomous, plan-workflow-autonomous-left, plan-workflow-open, arbiter-agent-tabs, describe-golden (ok). plan-milestone-check.txtar runs a plan run whose fake interview agent stays live (sleep) after `rhiza session done` and links its new change set to an intent under no milestone: the run exits non-zero, failed at step milestone-check, status --json names the change set, no review visit and no reviewer prompt; plan-workflow.txtar with the intent marked milestone=1 shows milestone-check pass, then review.

## Links

### Out

- verifies [plan-step-milestone-check](../capability/rhiza/orchestration/daemon/plan-requests/plan-step-milestone-check.md)

## History

- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
