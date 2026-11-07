# milestone-check-follows-policy

capability · `01M6VYMQ92BYC34XV88RA2A320` · parent [plan-step-milestone-check](../plan-step-milestone-check.md)

## acceptance

A unit test of rhiza run milestone-check covers: no milestone served and no policy.toml passes; no milestone served with backlog = "skip" fails naming the change set; a run with a milestone whose change set serves it directly, or through an intent under it, passes; one serving another milestone or none fails naming it. A testscript runs a plan run with a fake agent (left live after rhiza session done) whose change set serves an intent under no milestone: it enters review with no policy.toml, and ends failed at milestone-check with backlog = "skip"

## statement

rhiza run milestone-check passes a planned change set that serves no milestone, and fails, naming it and the fix, only when the run was started with a milestone (request --milestone, RHIZA_MILESTONE) and the change set does not serve that milestone directly or through an intent under it, or when the repo's .rhiza/policy.toml [schedule] backlog is "skip" (with today's message, since that work would never be built)

## Links

### Out

- depends-on [schedule-policy-file](../../schedule-policy-file.md)
- decided-by [scheduling-policy-shape](../../../../../../decision/scheduling-policy-shape.md)

### In

- depends-on [plan-workflow-file](../../../../workflows/built-in-workflows/plan-workflow/plan-workflow-file.md)
- verifies [milestone-check-follows-policy-verified](../../../../../../evidence/milestone-check-follows-policy-verified.md)
- verifies [milestone-check-follows-policy-qa-1](../../../../../../evidence/milestone-check-follows-policy-qa-1.md)
- verifies [pr-358-scheduling-as-policy](../../../../../../evidence/pr-358-scheduling-as-policy.md)

## History

- `01M6VY7GRV7BAJ9MMFVAY9GPZ2` Scheduling as policy: whether accepted change sets that serve no milestone are built, and in what order relative to milestone work, is a repo rule, not daemon code; today the daemon always skips them as 'serves no milestone' (SkipNoMilestone, internal/api/daemon.go:1159). The shipped default builds them after milestone work, oldest accepted first, so bug fixes, small features and repos with no milestones at all ship. A request's plan links a milestone only when one fits or the request names one with --milestone (relaxes the landed fix 01M6JDMJ, which made every request's plan serve a milestone). Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
