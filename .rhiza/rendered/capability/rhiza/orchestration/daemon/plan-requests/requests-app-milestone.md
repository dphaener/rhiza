# requests-app-milestone

capability · `01M6JE40PG6Z1ZG27N7DMFH3FF` · parent [plan-requests](../plan-requests.md)

## acceptance

A requestsapp unit test files a request with the milestone left at auto and gets argv with no --milestone, then one with a milestone chosen and gets --milestone <id> in the request argv; the entry view shows the milestone field

## statement

The requests Rex app's entry has an optional milestone field, defaulting to auto (the planner picks the open milestone it fits); a milestone chosen there is filed as the request's --milestone

## Links

### Out

- decided-by [request-plans-serve-a-milestone](../../../../../decision/request-plans-serve-a-milestone.md)

### In

- verifies [requests-app-milestone-verified](../../../../../evidence/requests-app-milestone-verified.md)
- verifies [requests-app-milestone-qa-1](../../../../../evidence/requests-app-milestone-qa-1.md)
- verifies [pr-270-fix-plan-runs-started-from-the](../../../../../evidence/pr-270-fix-plan-runs-started-from-the.md)

## History

- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
