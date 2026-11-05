# plan-step-milestone-check

capability · `01M6JE3T47QE8E8YPMBHR38MAE` · parent [plan-requests](../plan-requests.md)

## acceptance

A testscript runs a plan run whose interview (fake agent, left live after `rhiza session done`, not exited) links the change set serves to an intent under no milestone: the run ends failed at step milestone-check, `rhiza run status --json` shows that step's detail naming the change set, and the run never enters review; the same run serving an intent under a milestone passes milestone-check and enters review. A unit test of `rhiza run milestone-check` covers serving a milestone directly, through an intent under one, and serving none

## statement

The plan workflow gains a command step, milestone-check, between interview and review (interview's on.pass goes to it, its on.pass to review, and it has no on.fail): it runs `rhiza run milestone-check`, which exits non-zero, printing each planned change set (the intention itself, or each change set serving an intent a record adds) that serves no milestone directly or through an intent under one, and the fix, when any does. A failed command step with no on.fail ends the run failed with the command's output in the ended entry's detail (internal/workflow Definition.Next, internal/run/engine.go command), so the failure is loud whether or not the interview session exited, and the plan never reaches review or the gate

## Links

### Out

- decided-by [request-plans-serve-a-milestone](../../../../../decision/request-plans-serve-a-milestone.md)

### In

- verifies [plan-step-milestone-check-verified](../../../../../evidence/plan-step-milestone-check-verified.md)
- verifies [plan-step-milestone-check-qa-1](../../../../../evidence/plan-step-milestone-check-qa-1.md)
- verifies [plan-step-milestone-check-verified-2](../../../../../evidence/plan-step-milestone-check-verified-2.md)
- verifies [plan-step-milestone-check-qa-2](../../../../../evidence/plan-step-milestone-check-qa-2.md)
- verifies [pr-270-fix-plan-runs-started-from-the](../../../../../evidence/pr-270-fix-plan-runs-started-from-the.md)

## History

- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
