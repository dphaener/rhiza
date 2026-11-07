# review-01m6vy7e-r3-f1

question · `01M6VZK56QGXB742TV3WNV21K6`

## statement

Plan review round 3 left finding f1 (contradiction on daemon-picks-backlog): The new default (backlog = "after", used with no policy.toml) contradicts two landed testscripts that no leaf of the plan names. (1) cmd/rhiza/testdata/script/daemon-picks.txtar accepts a 'Backlog' change set with a leaf serving no milestone, then polls with --min-free 0 under the default max-runs 2 and expects exactly one 'Started a run of' (stdout -count=1, around lines 64-70). Under 'after' the Backlog change set would also start. (2) cmd/rhiza/testdata/script/daemon-status-projects.txtar (capability daemon-status-across-projects) expects each project's accepted no-milestone 'Orphan' to be skipped with reason no-milestone / 'serves no milestone' (lines 73-82). Under the default it would be picked. daemon-picks-backlog's acceptance covers only daemon_picks_test.go, and daemon-status-skips covers only daemon-status.txtar. A leaf's acceptance should say how these two scripts change (set backlog = "skip" in their repos, or expect the backlog pick). The critic confirmed it: No leaf's acceptance names daemon-picks.txtar or daemon-status-projects.txtar (daemon-status-skips names only daemon-status.txtar, daemon-picks-backlog only daemon_picks_test.go and a new testscript). In daemon-picks.txtar the Backlog change set has a leaf (backlog-leaf), serves only the non-milestone intent plain and is accepted; with no policy.toml and the default max-runs of 2, the 20 GiB poll expects stdout -count=1 '^Started a run of ', which fails once Backlog is picked under the default after. daemon-status-projects.txtar expects each project's accepted, leafed Orphan skipped with reason no-milestone / 'serves no milestone', which the default also turns into a pick. The plan contradicts these landed scripts without saying how they change. How should it be settled?

## Links

### Out

- blocks [daemon-picks-backlog](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-picks-backlog.md)

### In

- answers [backlog-skip-in-landed-daemon-scripts](../decision/backlog-skip-in-landed-daemon-scripts.md)

## History

- `01M6VY7GRV7BAJ9MMFVAY9GPZ2` Scheduling as policy: whether accepted change sets that serve no milestone are built, and in what order relative to milestone work, is a repo rule, not daemon code; today the daemon always skips them as 'serves no milestone' (SkipNoMilestone, internal/api/daemon.go:1159). The shipped default builds them after milestone work, oldest accepted first, so bug fixes, small features and repos with no milestones at all ship. A request's plan links a milestone only when one fits or the request names one with --milestone (relaxes the landed fix 01M6JDMJ, which made every request's plan serve a milestone). Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
