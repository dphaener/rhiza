# milestone-wait-owner-steps

capability · `01M767A0H2FFW8H8DY4N0DNY8Y` · parent [plan-requests](../plan-requests.md)

## acceptance

A unit test in internal/api/request_test.go covers daemonMilestoneWaits with a queued milestone request and another milestone request's run: live at the interview step (waits), live at a later step such as milestone-check or review (does not wait), stopped at a sign-off wait (waits), unjournaled within daemonStartGrace (waits); cmd/rhiza/testdata/script/daemon-milestone-request.txtar adds a case where the first milestone request's run journal has ended its interview step and is live at a later step, and the next daemon --once starts the second milestone request's plan run, while with --max-runs full it still starts nothing; go test -trimpath ./internal/api -run MilestoneWaits and go test -trimpath ./cmd/rhiza -run 'TestScripts/daemon-milestone-request' pass

## statement

A queued milestone request waits only while another milestone request's plan run holds the owner: its start is not journaled yet and is within daemonStartGrace, or the run is live with its current step the interview (start) step, or it is stopped waiting on the owner at a gate or sign-off step under owner autonomy; a live milestone plan run on any other step (milestone-check, review, critic, check, revise, open-pr, wait-pr, CI, rebase, merge, close-out) blocks nothing, so the next milestone interview starts as soon as the previous one's interview step has ended. daemonMilestoneWaits and daemonPlanRequests' milestoneLive apply this one rule, and rhiza.daemon.maxRuns still caps every start

## Links

### Out

- decided-by [milestone-wait-owner-steps-only](../../../../../decision/milestone-wait-owner-steps-only.md)
- derived-from [plan-fix-queued-milestone-request](../../../../../source/plan-fix-queued-milestone-request.md)

### In

- verifies [milestone-wait-owner-steps-verified](../../../../../evidence/milestone-wait-owner-steps-verified.md)
- verifies [milestone-wait-owner-steps-qa-1](../../../../../evidence/milestone-wait-owner-steps-qa-1.md)
- verifies [pr-453-fix-a-queued-milestone-request](../../../../../evidence/pr-453-fix-a-queued-milestone-request.md)

## History

- `01M767125EJM98D7H4ME2T67H5` Fix A queued milestone request waits for another milestone's whole plan run to end, not just its owner interview, so the owner can't start the next interview for an hour or more after finishing one. Evidence 2026-11-11: milestone request 01M75YW54KYRCF7TMNHSETKTEM (release channels) was queued at 18:45Z while plan run 01M75W8YT6M1F8W05GAJZSEC0H (milestone 32, workflow cleanup) was at its interview; the owner finished that interview at about 20:19Z and the run moved on to milestone-check (then review, PR, CI, merge), but the release request stays queued because daemonMilestoneWaits (internal/api/daemon_projects.go:423, and daemonPlanRequests' milestoneLive) treats any live milestone plan run as blocking. Raising rhiza.daemon.maxRuns does not help. The one-at-a-time rule exists so the owner holds one interview at a time; it should count only milestone plan runs whose current step is an owner-interactive step (the interview, or a gate waiting on the owner), so the next milestone interview can start as soon as the previous one's interview step has passed. Test: a testscript with two milestone requests where the first run is past its interview step and the second starts.
