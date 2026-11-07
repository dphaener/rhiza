# daemon-picks-schedule-qa-1

evidence · `01M6W2ESW0KH1ZR3VWTTZ9SXRN` · actor [claude-code-2e51135f](../actor/claude-code-2e51135f.md)

## result

PASS

## statement

QA with the branch build: rhiza daemon --once with a fake rex picks the milestone change set with a leaf under every backlog value, and backlog change sets only by the [schedule] backlog rule (after/before/skip/unknown checked by hand), each start logged once as rhiza run start --foreground <cs>.

## Links

### Out

- verifies [daemon-picks](../capability/rhiza/orchestration/daemon/daemon-picks.md)

## History

- `01M6VY7GRV7BAJ9MMFVAY9GPZ2` Scheduling as policy: whether accepted change sets that serve no milestone are built, and in what order relative to milestone work, is a repo rule, not daemon code; today the daemon always skips them as 'serves no milestone' (SkipNoMilestone, internal/api/daemon.go:1159). The shipped default builds them after milestone work, oldest accepted first, so bug fixes, small features and repos with no milestones at all ship. A request's plan links a milestone only when one fits or the request names one with --milestone (relaxes the landed fix 01M6JDMJ, which made every request's plan serve a milestone). Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
