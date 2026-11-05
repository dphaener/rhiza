# source-summary-names-milestone-check

decision · `01M6JF07XPB8FXFSR147TMP9VT`

## reasoning

The Decision is the ruling and the leaf plan-step-milestone-check already builds the command step; it gives the concrete reason the until check cannot work (the engine only fails an agent step on a failed until once its session exits, so a live session that said done is just polled again). The Source is the Decision's only provenance, since its wording is assumed, so a summary describing the rejected mechanism misleads anyone tracing why the guard exists. Fixing the Source's summary, and recording that the mechanism was revised in review rather than silently rewriting history, settles the contradiction without reopening the guard's design, which the change set's intent (fail loudly before acceptance instead of silently skipping) is served by either way.

## statement

Source plan-request-plans-serve-milestone's summary is corrected to match Decision request-plans-serve-a-milestone: the guard is a separate command step, milestone-check, right after the plan's interview and before review, failing the run when a planned change set serves no milestone directly or through an intent under one; it is not the interview's until/completion check. The summary notes that the guard's mechanism was revised from the planner's first wording (the plan step's completion check) during plan review. The Decision and leaf plan-step-milestone-check stay as they are.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6jdmg-r3-f1](../question/review-01m6jdmg-r3-f1.md)

## History

- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
