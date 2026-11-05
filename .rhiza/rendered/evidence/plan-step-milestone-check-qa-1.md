# plan-step-milestone-check-qa-1

evidence · `01M6MJ8M5DWFSC6B2A176BFT6F` · actor [claude-code-cfbfe790](../actor/claude-code-cfbfe790.md)

## result

FAIL

## statement

QA at e443109e, built binary, built-in plan workflow, fake claude as planner: serves a plain intent (area) or no serves link -> run ends failed at milestone-check naming the change set, never enters review; serves ms directly, or a landed intent under ms -> milestone-check passes, review entered. But serving a NEW intent under ms written into the plan change set itself (as the plan guide allows) -> milestone-check fails with serve no milestone: the check reads only the landed view, so the unlanded intent has no path to the milestone.

## Links

### Out

- verifies [plan-step-milestone-check](../capability/rhiza/orchestration/daemon/plan-requests/plan-step-milestone-check.md)

## History

- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
