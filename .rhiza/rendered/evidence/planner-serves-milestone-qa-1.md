# planner-serves-milestone-qa-1

evidence · `01M6MJ8E1Y2T763TJDNMMS66ZD` · actor [claude-code-cfbfe790](../actor/claude-code-cfbfe790.md)

## result

PASS

## statement

QA at e443109e: rhiza guide plan from the built binary says the serves link must reach a milestone directly or through an intent under one, naming RHIZA_MILESTONE, the open milestone it fits, and never an umbrella or area intent under no milestone; go test ./internal/guide passes.

## Links

### Out

- verifies [planner-serves-milestone](../capability/rhiza/orchestration/daemon/plan-requests/planner-serves-milestone.md)

## History

- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
