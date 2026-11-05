# pr-270-fix-plan-runs-started-from-the

evidence · `01M6MKW77GQZ7ZQGVTX1HBH33W` · actor [dphaener](../actor/dphaener.md)

## head

7a04679884c0ab3d2a698867c13095b73d90d53a

## pr

https://github.com/dphaener/rhiza/pull/270

## result

1 check passed: check SUCCESS

## statement

Pull request #270 (https://github.com/dphaener/rhiza/pull/270) passed CI at 7a04679

## status

verified

## Links

### Out

- verifies [request-names-milestone](../capability/rhiza/orchestration/daemon/plan-requests/request-names-milestone.md)
- verifies [plan-run-hands-milestone](../capability/rhiza/orchestration/daemon/plan-requests/plan-run-hands-milestone.md)
- verifies [planner-serves-milestone](../capability/rhiza/orchestration/daemon/plan-requests/planner-serves-milestone.md)
- verifies [plan-step-milestone-check](../capability/rhiza/orchestration/daemon/plan-requests/plan-step-milestone-check.md)
- verifies [requests-app-milestone](../capability/rhiza/orchestration/daemon/plan-requests/requests-app-milestone.md)
- verifies [plan-workflow-file](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-workflow-file.md)

## History

- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
