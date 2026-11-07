# pr-358-scheduling-as-policy

evidence · `01M6W34EWT3G45Q8EZ6ND27A62` · actor [dphaener](../actor/dphaener.md)

## head

333b9d914c008171755a20556285506d7cab7f25

## pr

https://github.com/dphaener/rhiza/pull/358

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #358 (https://github.com/dphaener/rhiza/pull/358) passed CI at 333b9d9

## status

verified

## Links

### Out

- verifies [schedule-policy-file](../capability/rhiza/orchestration/daemon/schedule-policy-file.md)
- verifies [daemon-picks-backlog](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-picks-backlog.md)
- verifies [daemon-picks](../capability/rhiza/orchestration/daemon/daemon-picks.md)
- verifies [plan-step-milestone-check](../capability/rhiza/orchestration/daemon/plan-requests/plan-step-milestone-check.md)
- verifies [planner-serves-milestone](../capability/rhiza/orchestration/daemon/plan-requests/planner-serves-milestone.md)
- verifies [milestone-check-follows-policy](../capability/rhiza/orchestration/daemon/plan-requests/plan-step-milestone-check/milestone-check-follows-policy.md)
- verifies [plan-guide-milestone-optional](../capability/rhiza/orchestration/daemon/plan-requests/planner-serves-milestone/plan-guide-milestone-optional.md)
- verifies [plan-workflow-file](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-workflow-file.md)
- verifies [daemon-status-skips](../capability/rhiza/orchestration/daemon/daemon-status/daemon-status-skips.md)

## History

- `01M6VY7GRV7BAJ9MMFVAY9GPZ2` Scheduling as policy: whether accepted change sets that serve no milestone are built, and in what order relative to milestone work, is a repo rule, not daemon code; today the daemon always skips them as 'serves no milestone' (SkipNoMilestone, internal/api/daemon.go:1159). The shipped default builds them after milestone work, oldest accepted first, so bug fixes, small features and repos with no milestones at all ship. A request's plan links a milestone only when one fits or the request names one with --milestone (relaxes the landed fix 01M6JDMJ, which made every request's plan serve a milestone). Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
