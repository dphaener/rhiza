# milestone-request-starts-at-once-qa-owner-steps-1

evidence · `01M76MXSBWB72GDD7325PC70Q8` · actor [claude-code-187fc5bc](../actor/claude-code-187fc5bc.md)

## result

PASS

## statement

QA at 42f7d662d with the built binary in a scratch root with a fake rex: the oldest queued milestone request started in a free slot; a second milestone request stayed queued while the first's plan run held the owner (start unjournaled within grace, live at interview, stopped at the gate) and with the slots full (--max-runs 1), and started in the next free slot once the first was past its interview.

## Links

### Out

- verifies [milestone-request-starts-at-once](../capability/rhiza/orchestration/daemon/plan-requests/milestone-request-starts-at-once.md)

## History

- `01M767125EJM98D7H4ME2T67H5` Fix A queued milestone request waits for another milestone's whole plan run to end, not just its owner interview, so the owner can't start the next interview for an hour or more after finishing one. Evidence 2026-11-11: milestone request 01M75YW54KYRCF7TMNHSETKTEM (release channels) was queued at 18:45Z while plan run 01M75W8YT6M1F8W05GAJZSEC0H (milestone 32, workflow cleanup) was at its interview; the owner finished that interview at about 20:19Z and the run moved on to milestone-check (then review, PR, CI, merge), but the release request stays queued because daemonMilestoneWaits (internal/api/daemon_projects.go:423, and daemonPlanRequests' milestoneLive) treats any live milestone plan run as blocking. Raising rhiza.daemon.maxRuns does not help. The one-at-a-time rule exists so the owner holds one interview at a time; it should count only milestone plan runs whose current step is an owner-interactive step (the interview, or a gate waiting on the owner), so the next milestone interview can start as soon as the previous one's interview step has passed. Test: a testscript with two milestone requests where the first run is past its interview step and the second starts.
