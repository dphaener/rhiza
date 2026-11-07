# backlog-skip-in-landed-daemon-scripts

decision · `01M6VZN3B2XGXST83BDNFHY084`

## reasoning

Both scripts use their no-milestone change set (Backlog, Orphan) on purpose to test the skip path: daemon-picks.txtar's header says a backlog change set is not started, and its disk-guard and permission-mode checks count exactly one start; daemon-status-projects.txtar (capability daemon-status-across-projects) is about listing each project's skipped change sets with their reason. Setting backlog = "skip" keeps what they test meaningful and exercises the skip value end to end in a real repo, matching how daemon-status-skips already settled daemon-status.txtar (a policy.toml with backlog = "skip"). Rewriting them to expect the backlog pick would drop their only skip-reason coverage across projects and tangle the disk-guard counts with backlog ordering, which the new testscript already covers. The work belongs on daemon-picks-backlog because that leaf changes the default that breaks them.

## statement

daemon-picks-backlog also covers the two landed testscripts: cmd/rhiza/testdata/script/daemon-picks.txtar gives its repo a .rhiza/policy.toml with [schedule] backlog = "skip", and cmd/rhiza/testdata/script/daemon-status-projects.txtar gives both projects (repo and second) the same file. Their expectations stay as they are: one 'Started a run of' per poll in daemon-picks.txtar, and each project's Orphan skipped with reason no-milestone / 'serves no milestone' in daemon-status-projects.txtar. The default (after) is tested by daemon-picks-backlog's own new testscript and Go test, not by rewriting these two.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6vy7e-r3-f1](../question/review-01m6vy7e-r3-f1.md)

## History

- `01M6VY7GRV7BAJ9MMFVAY9GPZ2` Scheduling as policy: whether accepted change sets that serve no milestone are built, and in what order relative to milestone work, is a repo rule, not daemon code; today the daemon always skips them as 'serves no milestone' (SkipNoMilestone, internal/api/daemon.go:1159). The shipped default builds them after milestone work, oldest accepted first, so bug fixes, small features and repos with no milestones at all ship. A request's plan links a milestone only when one fits or the request names one with --milestone (relaxes the landed fix 01M6JDMJ, which made every request's plan serve a milestone). Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
