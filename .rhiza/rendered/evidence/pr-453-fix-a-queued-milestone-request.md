# pr-453-fix-a-queued-milestone-request

evidence · `01M76N88SQ56BXPWFWY7E88BK3` · actor [dphaener](../actor/dphaener.md)

## head

fa3f0391c017be3ee873a43309250121a5c91dc6

## pr

https://github.com/dphaener/rhiza/pull/453

## result

11 checks passed: lint SUCCESS, go-install SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #453 (https://github.com/dphaener/rhiza/pull/453) passed CI at fa3f039

## status

verified

## Links

### Out

- verifies [milestone-request-starts-at-once](../capability/rhiza/orchestration/daemon/plan-requests/milestone-request-starts-at-once.md)
- verifies [milestone-wait-owner-steps](../capability/rhiza/orchestration/daemon/plan-requests/milestone-wait-owner-steps.md)

## History

- `01M767125EJM98D7H4ME2T67H5` Fix A queued milestone request waits for another milestone's whole plan run to end, not just its owner interview, so the owner can't start the next interview for an hour or more after finishing one. Evidence 2026-11-11: milestone request 01M75YW54KYRCF7TMNHSETKTEM (release channels) was queued at 18:45Z while plan run 01M75W8YT6M1F8W05GAJZSEC0H (milestone 32, workflow cleanup) was at its interview; the owner finished that interview at about 20:19Z and the run moved on to milestone-check (then review, PR, CI, merge), but the release request stays queued because daemonMilestoneWaits (internal/api/daemon_projects.go:423, and daemonPlanRequests' milestoneLive) treats any live milestone plan run as blocking. Raising rhiza.daemon.maxRuns does not help. The one-at-a-time rule exists so the owner holds one interview at a time; it should count only milestone plan runs whose current step is an owner-interactive step (the interview, or a gate waiting on the owner), so the next milestone interview can start as soon as the previous one's interview step has passed. Test: a testscript with two milestone requests where the first run is past its interview step and the second starts.
