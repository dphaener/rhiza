# plan-guide-milestone-optional-verified

evidence · `01M6W0N96QNJNXWVY4R4GVVMDQ` · actor [claude-code-48e174b7](../actor/claude-code-48e174b7.md)

## result

PASS

## statement

Verifier claude-code-48e174b7 at 77246b9f: read the diff of internal/guide/text/plan.md and plan_test.go. go test -trimpath ./internal/guide/... passes: TestPlanText's plan-serves-milestone phrases require RHIZA_MILESTONE, 'only when `RHIZA_MILESTONE` names one', 'or an open milestone fits', '.rhiza/policy.toml', '`backlog` rule' and '"skip"'; new TestPlanNoMandatoryMilestone asserts the old 'builds only a change set that serves a milestone' sentence is gone, and grep finds no other every-change-set-must-serve-a-milestone wording in plan.md. The installed rhiza-plan SKILL.md is built from the same plan topic body (internal/api/guide_install.go skillFile("rhiza-plan","plan")), so it matches by construction; go test -trimpath ./internal/api -run Guide passes.

## Links

### Out

- verifies [plan-guide-milestone-optional](../capability/rhiza/orchestration/daemon/plan-requests/planner-serves-milestone/plan-guide-milestone-optional.md)

## History

- `01M6VY7GRV7BAJ9MMFVAY9GPZ2` Scheduling as policy: whether accepted change sets that serve no milestone are built, and in what order relative to milestone work, is a repo rule, not daemon code; today the daemon always skips them as 'serves no milestone' (SkipNoMilestone, internal/api/daemon.go:1159). The shipped default builds them after milestone work, oldest accepted first, so bug fixes, small features and repos with no milestones at all ship. A request's plan links a milestone only when one fits or the request names one with --milestone (relaxes the landed fix 01M6JDMJ, which made every request's plan serve a milestone). Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
