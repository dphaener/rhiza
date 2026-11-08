# milestone-check-follows-policy-verified

evidence · `01M6W1YA25KXTCSPA0E8YQZ0Y2` · actor [claude-code-2ed25316](../actor/claude-code-2ed25316.md)

## result

PASS

## statement

At 8f4eb820, read the diff of internal/api/run_milestone.go: the check reads .rhiza/policy.toml, requires the run's started milestone (RHIZA_MILESTONE) when set, requires some milestone only when backlog = "skip", and otherwise passes, with errors naming the change set and a Fix. Ran go test -trimpath ./internal/api -run TestRunMilestoneCheck from a scratch worktree outside ~/code: TestRunMilestoneCheck (backlog = skip: serving nothing or a loose intent fails naming the change set; serving a milestone directly or through an intent passes) and TestRunMilestoneCheckBacklogAndNamedMilestone (no policy and backlog = before pass; with a milestone, serving it directly or through an intent passes; serving another milestone or none fails naming the change set and milestone) both PASS. Ran go test -trimpath ./cmd/rhiza -run TestScripts for plan-milestone-backlog (fake agent, change set serves an intent under no milestone, no policy.toml: milestone-check passes and the run enters review), plan-milestone-check (same with backlog = skip: fails at milestone-check, never reviewed), plan-workflow and describe-golden: ok.

## Links

### Out

- verifies [milestone-check-follows-policy](../capability/rhiza/orchestration/daemon/plan-requests/plan-step-milestone-check/milestone-check-follows-policy.md)

## History

- `01M6VY7GRV7BAJ9MMFVAY9GPZ2` Scheduling as policy: whether accepted change sets that serve no milestone are built, and in what order relative to milestone work, is a repo rule, not daemon code; today the daemon always skips them as 'serves no milestone' (SkipNoMilestone, internal/api/daemon.go:1159). The shipped default builds them after milestone work, oldest accepted first, so bug fixes, small features and repos with no milestones at all ship. A request's plan links a milestone only when one fits or the request names one with --milestone (relaxes the landed fix 01M6JDMJ, which made every request's plan serve a milestone). Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
