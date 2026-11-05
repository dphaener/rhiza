# plan-run-hands-milestone

capability · `01M6JE39K3W6P0FZ3AMZV7W550` · parent [plan-requests](../plan-requests.md)

## acceptance

A testscript starts the plan run of a request queued with --milestone (run start --request, with a fake agent that prints its environment) and sees RHIZA_MILESTONE set to the milestone's ID; for a request without --milestone it is unset

## statement

The plan run of a request that names a milestone hands that milestone to the planner: its interview step's environment carries RHIZA_MILESTONE with the milestone's ID, absent when the request names none

## Links

### Out

- depends-on [request-names-milestone](request-names-milestone.md)
- decided-by [request-plans-serve-a-milestone](../../../../../decision/request-plans-serve-a-milestone.md)

### In

- verifies [plan-run-hands-milestone-verified](../../../../../evidence/plan-run-hands-milestone-verified.md)
- verifies [plan-run-hands-milestone-qa-1](../../../../../evidence/plan-run-hands-milestone-qa-1.md)
- verifies [pr-270-fix-plan-runs-started-from-the](../../../../../evidence/pr-270-fix-plan-runs-started-from-the.md)

## History

- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
