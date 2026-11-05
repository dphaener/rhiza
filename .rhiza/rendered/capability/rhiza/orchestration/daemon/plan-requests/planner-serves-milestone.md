# planner-serves-milestone

capability · `01M6JE3KEGWGBYW35XHSH0JJK1` · parent [plan-requests](../plan-requests.md)

## acceptance

internal/guide plan_test asserts the plan guide names RHIZA_MILESTONE, says the serves link must reach a milestone directly or through an intent under one, and rules out umbrella and area intents under no milestone; make check passes with the skill copy in step

## statement

The plan guide (rhiza guide plan and the rhiza-plan skill) tells the planner that an intention change set must serve a milestone, directly or through an intent under one: the milestone in RHIZA_MILESTONE when set, else the open milestone it fits, never an umbrella or area intent under no milestone

## Links

### Out

- decided-by [request-plans-serve-a-milestone](../../../../../decision/request-plans-serve-a-milestone.md)

### In

- verifies [planner-serves-milestone-verified](../../../../../evidence/planner-serves-milestone-verified.md)
- verifies [planner-serves-milestone-qa-1](../../../../../evidence/planner-serves-milestone-qa-1.md)
- verifies [pr-270-fix-plan-runs-started-from-the](../../../../../evidence/pr-270-fix-plan-runs-started-from-the.md)

## History

- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
