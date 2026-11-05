# request-names-milestone-qa-2

evidence · `01M6MK3DKRKDBM4GS65WJEPRZ7` · actor [claude-code-c5ae7d72](../actor/claude-code-c5ae7d72.md)

## result

PASS

## statement

QA at ac9cf15f, built binary in a scratch repo: rhiza request 'Add paging' --milestone ms queued it and request ls --json shows milestone=<ms ID>; --milestone area (an intent with no milestone field) and --milestone nosuch each exit 2 with a usage error and queue nothing (request ls holds only the one request).

## Links

### Out

- verifies [request-names-milestone](../capability/rhiza/orchestration/daemon/plan-requests/request-names-milestone.md)

## History

- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
