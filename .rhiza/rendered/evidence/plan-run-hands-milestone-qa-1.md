# plan-run-hands-milestone-qa-1

evidence · `01M6MJ7ZQEKSDPH7PJZASJYJC7` · actor [claude-code-cfbfe790](../actor/claude-code-cfbfe790.md)

## result

PASS

## statement

QA at e443109e, built binary with fake rex/claude: two requests (one --milestone ms, one without); rhiza daemon --once started both plan runs; the fake claude logged RHIZA_MILESTONE=<ms ID> for the first run and unset for the second, and the session records carry milestone accordingly.

## Links

### Out

- verifies [plan-run-hands-milestone](../capability/rhiza/orchestration/daemon/plan-requests/plan-run-hands-milestone.md)

## History

- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
