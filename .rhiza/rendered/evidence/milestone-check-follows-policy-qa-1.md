# milestone-check-follows-policy-qa-1

evidence · `01M6W2ET4G4CTN4CW4EDDQRKWT` · actor [claude-code-2e51135f](../actor/claude-code-2e51135f.md)

## result

PASS

## statement

QA with the branch build as rhiza on PATH, the plan-milestone-check.txtar fixtures (fake claude/rex/gh) run by hand: rhiza run start --workflow plan on a plan serving an intent under no milestone ended failed at milestone-check naming the change set and policy.toml's backlog = "skip" when that file was present, and with no policy.toml passed milestone-check, entered review and the run succeeded.

## Links

### Out

- verifies [milestone-check-follows-policy](../capability/rhiza/orchestration/daemon/plan-requests/plan-step-milestone-check/milestone-check-follows-policy.md)

## History

- `01M6VY7GRV7BAJ9MMFVAY9GPZ2` Scheduling as policy: whether accepted change sets that serve no milestone are built, and in what order relative to milestone work, is a repo rule, not daemon code; today the daemon always skips them as 'serves no milestone' (SkipNoMilestone, internal/api/daemon.go:1159). The shipped default builds them after milestone work, oldest accepted first, so bug fixes, small features and repos with no milestones at all ship. A request's plan links a milestone only when one fits or the request names one with --milestone (relaxes the landed fix 01M6JDMJ, which made every request's plan serve a milestone). Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
