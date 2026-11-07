# daemon-picks-backlog-verified

evidence · `01M6W1C5S219KHSARHCEEF5A68` · actor [claude-code-449ee551](../actor/claude-code-449ee551.md)

## result

PASS

## statement

Verified at 4aa87845 as the verifier. Read the diff: daemonCandidates takes a daemonSchedule (Backlog, Detail, Since); backlog groups' accepted change sets go after milestone candidates by default, ahead under before, sorted by acceptedSince then ID; under skip or an unknown value (Backlog "", policy.Read's error in Detail) they are skipped no-milestone with that detail; daemonFold reads .rhiza/policy.toml and daemonStamp hashes it. Ran go test -trimpath ./internal/api -run 'TestDaemonCandidatesBacklog|TestDaemonPicks$|TestUmbrellaNotReady' -v: all PASS. TestDaemonCandidatesBacklog covers two milestone and two backlog change sets with different accepted times: after gives m1,m2,b2,b1; before gives b2,b1,m1,m2; skip gives m1,m2 with 2 skips; unknown gives m1,m2 with both skips no-milestone, detail carried. TestDaemonPicks now sets backlog skip. go test ./internal/policy PASS (unknown value error names the value). go test -trimpath ./cmd/rhiza -run 'TestScripts/(daemon-picks-backlog|daemon-picks|daemon-status|daemon-status-projects)$' ok: the new daemon-picks-backlog.txtar starts the only, no-milestone change set with no policy.toml, does not start it under backlog = "skip" (status says serves no milestone), and shows 'sometimes' in daemon status --json for an unknown value; the landed daemon scripts set backlog = "skip" as review-01m6vy7e-r3-f1's answer settled. The broader api run was flagged only for a write to ~/code/.rhiza-root/usage.live.json, which live Claude Code statusline sessions write, not this change.

## Links

### Out

- verifies [daemon-picks-backlog](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-picks-backlog.md)

## History

- `01M6VY7GRV7BAJ9MMFVAY9GPZ2` Scheduling as policy: whether accepted change sets that serve no milestone are built, and in what order relative to milestone work, is a repo rule, not daemon code; today the daemon always skips them as 'serves no milestone' (SkipNoMilestone, internal/api/daemon.go:1159). The shipped default builds them after milestone work, oldest accepted first, so bug fixes, small features and repos with no milestones at all ship. A request's plan links a milestone only when one fits or the request names one with --milestone (relaxes the landed fix 01M6JDMJ, which made every request's plan serve a milestone). Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
