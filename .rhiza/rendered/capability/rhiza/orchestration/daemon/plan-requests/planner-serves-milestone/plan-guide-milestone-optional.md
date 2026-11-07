# plan-guide-milestone-optional

capability · `01M6VYMQBZFQMEAXX3YSKEYS5J` · parent [planner-serves-milestone](../planner-serves-milestone.md)

## acceptance

internal/guide plan_test asserts the plan guide names RHIZA_MILESTONE, says a milestone is linked only when RHIZA_MILESTONE is set or one fits, names .rhiza/policy.toml's backlog rule and its "skip" exception, and no longer says every change set must serve a milestone; the installed skill copy matches the guide text

## statement

The plan guide (rhiza guide plan and the rhiza-plan skill) says an intention change set serves the intent it belongs under and links a milestone only when RHIZA_MILESTONE names one (then the serves link must reach it) or an open milestone fits; work no milestone fits is built as backlog by the repo's .rhiza/policy.toml [schedule] backlog rule, unless that rule is "skip", when it must serve a milestone as before

## Links

### Out

- decided-by [scheduling-policy-shape](../../../../../../decision/scheduling-policy-shape.md)

### In

- verifies [plan-guide-milestone-optional-verified](../../../../../../evidence/plan-guide-milestone-optional-verified.md)
- verifies [plan-guide-milestone-optional-qa-1](../../../../../../evidence/plan-guide-milestone-optional-qa-1.md)
- verifies [pr-358-scheduling-as-policy](../../../../../../evidence/pr-358-scheduling-as-policy.md)

## History

- `01M6VY7GRV7BAJ9MMFVAY9GPZ2` Scheduling as policy: whether accepted change sets that serve no milestone are built, and in what order relative to milestone work, is a repo rule, not daemon code; today the daemon always skips them as 'serves no milestone' (SkipNoMilestone, internal/api/daemon.go:1159). The shipped default builds them after milestone work, oldest accepted first, so bug fixes, small features and repos with no milestones at all ship. A request's plan links a milestone only when one fits or the request names one with --milestone (relaxes the landed fix 01M6JDMJ, which made every request's plan serve a milestone). Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
