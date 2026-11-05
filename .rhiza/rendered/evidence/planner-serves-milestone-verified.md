# planner-serves-milestone-verified

evidence · `01M6MF63GC42KXENSSNP4XYG3V` · actor [claude-code-19152758](../actor/claude-code-19152758.md)

## result

PASS

## statement

At dffc827d: go test -trimpath -count=1 -run TestPlan ./internal/guide passes; TestPlanText's new plan-serves-milestone leaf asserts RHIZA_MILESTONE, 'reach one, directly or through an intent under one' and 'umbrella or area intent under no' milestone. A fresh build's `rhiza guide plan` prints the new paragraph in step 2; the rhiza-plan skill is generated from the same embedded plan.md, so the skill copy is in step. make self-check reports no failure except render staleness of actor files and index.md, which rhiza render regenerates after landing; fmt/vet/full suite left to the PR's CI.

## Links

### Out

- verifies [planner-serves-milestone](../capability/rhiza/orchestration/daemon/plan-requests/planner-serves-milestone.md)

## History

- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
