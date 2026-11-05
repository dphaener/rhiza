# plan-step-milestone-check-qa-2

evidence · `01M6MK36YKRHWQXWW4NAEM6XXF` · actor [claude-code-c5ae7d72](../actor/claude-code-c5ae7d72.md)

## result

PASS

## statement

QA at ac9cf15f, built binary, built-in plan workflow via run start --foreground, fake rex/claude planner left live after session done, landed milestone ms and area intent: intention change set adding a NEW unlanded intent under ms into itself and serving it -> milestone-check pass, then review (the cycle-1 failure is fixed); serving ms directly -> pass; adding a new milestone intent with a child and serving the child -> pass; serving a new intent with no parent, or one under the landed area intent -> run ends failed at milestone-check (fail, then end) naming the change set with the Fix line, never enters review. Record plan that lands the record then serves a new change set to an intent under ms -> pass.

## Links

### Out

- verifies [plan-step-milestone-check](../capability/rhiza/orchestration/daemon/plan-requests/plan-step-milestone-check.md)

## History

- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
