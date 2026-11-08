# plan-guide-milestone-optional-qa-1

evidence · `01M6W2ET9W6PFXZPX88305S0EV` · actor [claude-code-2e51135f](../actor/claude-code-2e51135f.md)

## result

PASS

## statement

QA with the branch build: rhiza guide plan says to link a milestone only when RHIZA_MILESTONE names one or an open milestone fits, that other work is built as backlog by .rhiza/policy.toml's [schedule] backlog rule, and that under "skip" the serves link must reach a milestone.

## Links

### Out

- verifies [plan-guide-milestone-optional](../capability/rhiza/orchestration/daemon/plan-requests/planner-serves-milestone/plan-guide-milestone-optional.md)

## History

- `01M6VY7GRV7BAJ9MMFVAY9GPZ2` Scheduling as policy: whether accepted change sets that serve no milestone are built, and in what order relative to milestone work, is a repo rule, not daemon code; today the daemon always skips them as 'serves no milestone' (SkipNoMilestone, internal/api/daemon.go:1159). The shipped default builds them after milestone work, oldest accepted first, so bug fixes, small features and repos with no milestones at all ship. A request's plan links a milestone only when one fits or the request names one with --milestone (relaxes the landed fix 01M6JDMJ, which made every request's plan serve a milestone). Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
