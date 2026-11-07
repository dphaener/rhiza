# daemon-status-skips

capability · `01M686TKZGT5S1SHQSCK4QPPRR` · parent [daemon-status](../daemon-status.md)

## acceptance

daemon-status.txtar sets up an accepted change set serving no milestone in a repo whose .rhiza/policy.toml sets [schedule] backlog = "skip", one with no leaves and one depending on an unlanded change set, runs daemon --once, and checks Skipped: lines with serves no milestone and the JSON skipped list with reason not-ready and its reasons; describe-golden.txtar is regenerated

## statement

rhiza daemon status and its --json list each skipped change set by ID and title with its reason, and its readiness reasons when it is not ready, as of the daemon's last poll

## Links

### Out

- depends-on [daemon-picks-backlog](../daemon-picks/daemon-picks-backlog.md)
- decided-by [daemon-skips-saved-at-poll](../../../../../decision/daemon-skips-saved-at-poll.md)

### In

- verifies [daemon-status-skips-verified](../../../../../evidence/daemon-status-skips-verified.md)
- verifies [daemon-status-skips-qa-1](../../../../../evidence/daemon-status-skips-qa-1.md)
- verifies [pr-129-daemon-explains-skips](../../../../../evidence/pr-129-daemon-explains-skips.md)
- verifies [pr-358-scheduling-as-policy](../../../../../evidence/pr-358-scheduling-as-policy.md)

## History

- `01M686RS47EK7F4H2WY15TNAQ5` rhiza daemon status names each accepted change set the daemon does not start and why
- `01M6VY7GRV7BAJ9MMFVAY9GPZ2` Scheduling as policy: whether accepted change sets that serve no milestone are built, and in what order relative to milestone work, is a repo rule, not daemon code; today the daemon always skips them as 'serves no milestone' (SkipNoMilestone, internal/api/daemon.go:1159). The shipped default builds them after milestone work, oldest accepted first, so bug fixes, small features and repos with no milestones at all ship. A request's plan links a milestone only when one fits or the request names one with --milestone (relaxes the landed fix 01M6JDMJ, which made every request's plan serve a milestone). Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
