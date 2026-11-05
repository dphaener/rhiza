# request-names-milestone-qa-1

evidence · `01M6MJ7RTZQA0JXRVNRQ3KTA9P` · actor [claude-code-cfbfe790](../actor/claude-code-cfbfe790.md)

## result

PASS

## statement

QA at e443109e, built binary in a scratch repo: rhiza request --milestone ms, by full ID and by short ID prefix each queued with the milestone ID shown in request ls --json (and --dry-run --json, and api call request); --milestone naming a plain intent, a capability, a retired milestone (by ID) or an unknown handle exited 2 with a usage error and queued nothing. The human request ls lines do not show the milestone.

## Links

### Out

- verifies [request-names-milestone](../capability/rhiza/orchestration/daemon/plan-requests/request-names-milestone.md)

## History

- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
