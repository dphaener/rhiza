# daemon-picks-backlog-qa-1

evidence · `01M6W2ESRRB9N6A3RWEGHSTW8G` · actor [claude-code-2e51135f](../actor/claude-code-2e51135f.md)

## result

PASS

## statement

QA with the branch build (go build ./cmd/rhiza) in a scratch Rhiza root, fake rex, RHIZA_TEST_DETACH_LOG: one milestone change set and two backlog change sets with leaves. rhiza daemon --once --max-runs=3 with no policy.toml started milestone first then backlog; backlog = "before" started backlog first; backlog = "skip" started only the milestone one and daemon status listed each backlog one as 'serves no milestone'; backlog = "x" did the same with the policy error as the skip detail (text and --json reason no-milestone). Backlog order follows acceptance time, not ID: a change set made first but accepted later started second.

## Links

### Out

- verifies [daemon-picks-backlog](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-picks-backlog.md)

## History

- `01M6VY7GRV7BAJ9MMFVAY9GPZ2` Scheduling as policy: whether accepted change sets that serve no milestone are built, and in what order relative to milestone work, is a repo rule, not daemon code; today the daemon always skips them as 'serves no milestone' (SkipNoMilestone, internal/api/daemon.go:1159). The shipped default builds them after milestone work, oldest accepted first, so bug fixes, small features and repos with no milestones at all ship. A request's plan links a milestone only when one fits or the request names one with --milestone (relaxes the landed fix 01M6JDMJ, which made every request's plan serve a milestone). Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
