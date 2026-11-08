# plan-step-milestone-check

capability · `01M6JE3T47QE8E8YPMBHR38MAE` · parent [plan-requests](../plan-requests.md)

## acceptance

A unit test of rhiza run milestone-check covers: serving no milestone with no policy.toml passes; serving no milestone with backlog = "skip" fails naming the change set; a run with a milestone whose change set serves that milestone directly, or through an intent under it, passes; one serving another milestone or none fails naming it. A testscript runs a plan run with a fake agent (left live after rhiza session done) whose change set serves an intent under no milestone: it enters review with no policy.toml, and ends failed at milestone-check with backlog = "skip"

## statement

The plan workflow keeps its command step milestone-check between interview and review (no on.fail): `rhiza run milestone-check` passes a planned change set (the intention itself, or each change set serving an intent a record adds) that serves no milestone, except that it fails, printing each offending change set and the fix, when the run was started with a milestone (the request's --milestone, the run's RHIZA_MILESTONE) and the change set does not serve that milestone directly or through an intent under it, or when the repo's [schedule] backlog is "skip" and the change set serves no milestone (today's message, since that work would never be built); a failure ends the run failed with the output in the ended entry's detail and the plan never reaches review or the gate. A repo that wants no check drops the step from its own plan workflow

## Links

### Out

- decided-by [scheduling-policy-shape](../../../../../decision/scheduling-policy-shape.md)

### In

- verifies [plan-step-milestone-check-verified](../../../../../evidence/plan-step-milestone-check-verified.md)
- verifies [plan-step-milestone-check-qa-1](../../../../../evidence/plan-step-milestone-check-qa-1.md)
- verifies [plan-step-milestone-check-verified-2](../../../../../evidence/plan-step-milestone-check-verified-2.md)
- verifies [plan-step-milestone-check-qa-2](../../../../../evidence/plan-step-milestone-check-qa-2.md)
- verifies [pr-270-fix-plan-runs-started-from-the](../../../../../evidence/pr-270-fix-plan-runs-started-from-the.md)
- verifies [plan-step-milestone-check-policy-qa-1](../../../../../evidence/plan-step-milestone-check-policy-qa-1.md)
- verifies [pr-358-scheduling-as-policy](../../../../../evidence/pr-358-scheduling-as-policy.md)

## History

- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
- `01M6VY7GRV7BAJ9MMFVAY9GPZ2` Scheduling as policy: whether accepted change sets that serve no milestone are built, and in what order relative to milestone work, is a repo rule, not daemon code; today the daemon always skips them as 'serves no milestone' (SkipNoMilestone, internal/api/daemon.go:1159). The shipped default builds them after milestone work, oldest accepted first, so bug fixes, small features and repos with no milestones at all ship. A request's plan links a milestone only when one fits or the request names one with --milestone (relaxes the landed fix 01M6JDMJ, which made every request's plan serve a milestone). Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
