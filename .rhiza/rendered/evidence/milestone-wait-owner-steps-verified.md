# milestone-wait-owner-steps-verified

evidence · `01M76MQGXDEJ95ATR12XW7EX5R` · actor [claude-code-fe4fed09](../actor/claude-code-fe4fed09.md)

## result

PASS

## statement

Verifier at 2932c7a22: ran go test -trimpath ./internal/api -run MilestoneWaits -v (TestDaemonMilestoneWaitsOwnerSteps and TestDaemonMilestoneWaitsAfterEndedRun PASS) and go test -trimpath ./cmd/rhiza -run 'TestScripts/daemon-milestone-request' (PASS). The unit test covers the queued second milestone request waiting while the first's run is unjournaled within daemonStartGrace, live on another host at the interview step, and stopped waiting on the owner at the gate; and not waiting when live at milestone-check or review. The built-in plan-milestone has no signoff-kind step: its owner sign-off stop is the gate (rhiza run gate), which is what the test exercises; milestoneHoldsOwner treats KindSignoff steps the same way. The txtar adds the case: first run live at interview keeps the second queued; after it enters milestone-check, --max-runs 2 (full) starts nothing and --max-runs 5 starts the second's plan-milestone run exactly once. Diff read: daemonMilestoneWaits and daemonPlanRequests' milestoneLive both call the new milestoneHoldsOwner, so one rule applies, and the maxRuns cap is untouched.

## Links

### Out

- verifies [milestone-wait-owner-steps](../capability/rhiza/orchestration/daemon/plan-requests/milestone-wait-owner-steps.md)

## History

- `01M767125EJM98D7H4ME2T67H5` Fix A queued milestone request waits for another milestone's whole plan run to end, not just its owner interview, so the owner can't start the next interview for an hour or more after finishing one. Evidence 2026-11-11: milestone request 01M75YW54KYRCF7TMNHSETKTEM (release channels) was queued at 18:45Z while plan run 01M75W8YT6M1F8W05GAJZSEC0H (milestone 32, workflow cleanup) was at its interview; the owner finished that interview at about 20:19Z and the run moved on to milestone-check (then review, PR, CI, merge), but the release request stays queued because daemonMilestoneWaits (internal/api/daemon_projects.go:423, and daemonPlanRequests' milestoneLive) treats any live milestone plan run as blocking. Raising rhiza.daemon.maxRuns does not help. The one-at-a-time rule exists so the owner holds one interview at a time; it should count only milestone plan runs whose current step is an owner-interactive step (the interview, or a gate waiting on the owner), so the next milestone interview can start as soon as the previous one's interview step has passed. Test: a testscript with two milestone requests where the first run is past its interview step and the second starts.
