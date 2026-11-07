# daemon-picks-backlog

capability · `01M6VYK4TRY2VN9XHGKRD9V4XZ` · parent [daemon-picks](../daemon-picks.md)

## acceptance

A Go test of daemonCandidates over a fixture with two milestone change sets and two backlog change sets with leaves, accepted at different times, gives milestone then backlog-oldest-first under after, backlog-oldest-first then milestone under before, only the milestone ones with both backlog ones skipped no-milestone under skip, and under an unknown value the skip with the value in its detail; the existing case in internal/api/daemon_picks_test.go that expects SkipNoMilestone for a backlog change set with no policy now expects it picked (or sets backlog = "skip"); a testscript runs rhiza daemon --once with a fake Rex on a repo whose only accepted change set serves no milestone, sees it started with no policy.toml, and not started with backlog = "skip"

## statement

The daemon picks accepted, ready change sets with a leaf that serve no milestone by the repo's [schedule] backlog rule: after (the default), after every milestone group's candidates; before, ahead of them; skip, never, skipped as 'serves no milestone' as today. Within the backlog it picks the oldest accepted first (by the time each change set was last accepted or reopened). With an unknown backlog value it skips every backlog change set as 'serves no milestone' and puts the bad value in the skip's detail, so rhiza daemon status shows it

## Links

### Out

- depends-on [schedule-policy-file](../schedule-policy-file.md)
- decided-by [scheduling-policy-shape](../../../../../decision/scheduling-policy-shape.md)

### In

- depends-on [daemon-status-skips](../daemon-status/daemon-status-skips.md)
- blocks [review-01m6vy7e-r3-f1](../../../../../question/review-01m6vy7e-r3-f1.md)
- verifies [daemon-picks-backlog-verified](../../../../../evidence/daemon-picks-backlog-verified.md)
- verifies [daemon-picks-backlog-qa-1](../../../../../evidence/daemon-picks-backlog-qa-1.md)
- verifies [pr-358-scheduling-as-policy](../../../../../evidence/pr-358-scheduling-as-policy.md)

## History

- `01M6VY7GRV7BAJ9MMFVAY9GPZ2` Scheduling as policy: whether accepted change sets that serve no milestone are built, and in what order relative to milestone work, is a repo rule, not daemon code; today the daemon always skips them as 'serves no milestone' (SkipNoMilestone, internal/api/daemon.go:1159). The shipped default builds them after milestone work, oldest accepted first, so bug fixes, small features and repos with no milestones at all ship. A request's plan links a milestone only when one fits or the request names one with --milestone (relaxes the landed fix 01M6JDMJ, which made every request's plan serve a milestone). Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
