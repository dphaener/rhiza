# schedule-policy-file

capability · `01M6VYK4NBNG3NG4FKXD1VP896` · parent [daemon](../daemon.md)

## acceptance

A Go unit test of the policy reader covers a missing file, a missing key, each of after, before and skip, an unknown value and malformed TOML; a testscript (policy-check.txtar) runs rhiza check on a repo with backlog = "sometimes" and sees a problem naming .rhiza/policy.toml and the value, and passes once the value is "skip"

## statement

A repo sets its scheduling rule in a checked-in .rhiza/policy.toml, table [schedule], key backlog: "after" (the default, also used when the file or key is missing), "before" or "skip"; rhiza check reports a policy.toml that does not parse, or a backlog value outside those three, as a problem naming the file and the bad value, and repo.json stays identity only

## Links

### Out

- decided-by [scheduling-policy-shape](../../../../decision/scheduling-policy-shape.md)

### In

- depends-on [milestone-check-follows-policy](plan-requests/plan-step-milestone-check/milestone-check-follows-policy.md)
- depends-on [daemon-picks-backlog](daemon-picks/daemon-picks-backlog.md)
- verifies [schedule-policy-file-verified](../../../../evidence/schedule-policy-file-verified.md)
- verifies [schedule-policy-file-qa-1](../../../../evidence/schedule-policy-file-qa-1.md)
- verifies [pr-358-scheduling-as-policy](../../../../evidence/pr-358-scheduling-as-policy.md)

## History

- `01M6VY7GRV7BAJ9MMFVAY9GPZ2` Scheduling as policy: whether accepted change sets that serve no milestone are built, and in what order relative to milestone work, is a repo rule, not daemon code; today the daemon always skips them as 'serves no milestone' (SkipNoMilestone, internal/api/daemon.go:1159). The shipped default builds them after milestone work, oldest accepted first, so bug fixes, small features and repos with no milestones at all ship. A request's plan links a milestone only when one fits or the request names one with --milestone (relaxes the landed fix 01M6JDMJ, which made every request's plan serve a milestone). Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
