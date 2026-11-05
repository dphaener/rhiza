# review-01m6jdmg-r3-f1

question · `01M6JEYNH0TXDW0AEPG7TEXZ7F`

## statement

Plan review round 3 left finding f1 (contradiction on plan-request-plans-serve-milestone): The Source's summary says the guard is 'the plan step's completion check fails the step when a planned change set serves no milestone', but Decision request-plans-serve-a-milestone (and leaf plan-step-milestone-check) rule the guard is a separate command step, milestone-check, between interview and review, and explicitly 'Not the interview's until check'. The Source records the earlier wording and now contradicts the Decision it backs; its summary should describe the milestone-check command step. The critic confirmed it: Source plan-request-plans-serve-milestone's summary still says 'the plan step's completion check fails the step when a planned change set serves no milestone', while the Decision it backs (request-plans-serve-a-milestone, via derived-from) now rules the guard is a command step right after the interview and explicitly 'Not the interview's until check'. The Source is the Decision's only provenance (wording is assumed), so its summary describes the mechanism the Decision rejects; it should describe the command step (or note the round-2 revision). How should it be settled?

## Links

### In

- answers [source-summary-names-milestone-check](../decision/source-summary-names-milestone-check.md)

## History

- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
