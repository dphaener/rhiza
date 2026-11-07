# planner-serves-milestone

capability · `01M6JE3KEGWGBYW35XHSH0JJK1` · parent [plan-requests](../plan-requests.md)

## acceptance

internal/guide plan_test asserts the plan guide names RHIZA_MILESTONE, says a milestone is linked only when RHIZA_MILESTONE is set or one fits, names .rhiza/policy.toml's backlog rule and its "skip" exception, and no longer says every change set must serve a milestone; the installed skill copy matches the guide text

## statement

The plan guide (rhiza guide plan and the rhiza-plan skill) tells the planner that an intention change set serves the intent it belongs under, and links a milestone only when the run hands one in RHIZA_MILESTONE (then the serves link must reach that milestone, directly or through an intent under it) or an open milestone fits the work; work no milestone fits serves its intent and is built as backlog by the repo's [schedule] rule, unless that rule is "skip", when it must serve a milestone as before

## Links

### Out

- decided-by [scheduling-policy-shape](../../../../../decision/scheduling-policy-shape.md)

### In

- verifies [planner-serves-milestone-verified](../../../../../evidence/planner-serves-milestone-verified.md)
- verifies [planner-serves-milestone-qa-1](../../../../../evidence/planner-serves-milestone-qa-1.md)
- verifies [pr-270-fix-plan-runs-started-from-the](../../../../../evidence/pr-270-fix-plan-runs-started-from-the.md)
- verifies [planner-serves-milestone-policy-qa-1](../../../../../evidence/planner-serves-milestone-policy-qa-1.md)
- verifies [pr-358-scheduling-as-policy](../../../../../evidence/pr-358-scheduling-as-policy.md)

## History

- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
- `01M6VY7GRV7BAJ9MMFVAY9GPZ2` Scheduling as policy: whether accepted change sets that serve no milestone are built, and in what order relative to milestone work, is a repo rule, not daemon code; today the daemon always skips them as 'serves no milestone' (SkipNoMilestone, internal/api/daemon.go:1159). The shipped default builds them after milestone work, oldest accepted first, so bug fixes, small features and repos with no milestones at all ship. A request's plan links a milestone only when one fits or the request names one with --milestone (relaxes the landed fix 01M6JDMJ, which made every request's plan serve a milestone). Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
