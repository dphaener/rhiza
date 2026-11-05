# plan-step-milestone-check-verified-2

evidence · `01M6MJMSK038VAM9T4BBMYB13E` · actor [claude-code-22010156](../actor/claude-code-22010156.md)

## result

PASS

## statement

Attempt 2 at 7e3a3780, verifier: read the fix diff (run milestone-check now resolves each in-flight plan change set's served intents in that change set's what-if view via milestoneServed, fixing QA f1 where an intent the plan adds under a milestone in its own change set was invisible in to-be). Ran go test -trimpath ./internal/api -run TestRunMilestoneCheck -count=1: PASS, covering serving a milestone directly, through an intent under one, serving none, and the new case of an unlanded intent the change set adds under the milestone. Ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(plan-milestone-check|plan-run-milestone-env|plan-workflow.*|roadmap-milestones|describe-golden)$' -count=1: ok. plan-milestone-check.txtar: interview fake left live, run ends failed at milestone-check, status --json names the change set, never enters review; plan-workflow*.txtar pass milestone-check then enter review. plan.toml has interview on.pass=milestone-check, milestone-check runs rhiza run milestone-check with on.pass review and no on.fail.

## Links

### Out

- verifies [plan-step-milestone-check](../capability/rhiza/orchestration/daemon/plan-requests/plan-step-milestone-check.md)

## History

- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
