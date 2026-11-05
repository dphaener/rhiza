# request-names-milestone

capability · `01M6JE32Q9Y16GT0CG6254C38Y` · parent [plan-requests](../plan-requests.md)

## acceptance

A testscript queues a request with --milestone <handle> and sees the milestone's ID in request ls --json; --milestone naming an intent with no milestone field, or an unknown node, exits with a usage error and queues nothing

## statement

rhiza request (and the request op) takes an optional --milestone naming a milestone intent by handle or ID, refuses with a usage error one that is not a live milestone intent, stores its ID on the request and shows it in request ls

## Links

### Out

- decided-by [request-plans-serve-a-milestone](../../../../../decision/request-plans-serve-a-milestone.md)

### In

- depends-on [plan-run-hands-milestone](plan-run-hands-milestone.md)
- verifies [request-names-milestone-verified](../../../../../evidence/request-names-milestone-verified.md)
- verifies [request-names-milestone-qa-1](../../../../../evidence/request-names-milestone-qa-1.md)
- verifies [request-names-milestone-qa-2](../../../../../evidence/request-names-milestone-qa-2.md)
- verifies [pr-270-fix-plan-runs-started-from-the](../../../../../evidence/pr-270-fix-plan-runs-started-from-the.md)

## History

- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
