# plan-fix-queued-milestone-request

source · `01M7678ZK6T2FQHD4MYTNWKWRE`

## quotes

```json
[]
```

## raw_sha256

ea8763729b1da69fa8b58164e7b8fa80578bf5412f39b5e75ce9c5b1a629e7e9

## session

a77d1ef4-7b1c-471a-b80e-01a1cc70c1d1

## summary

Autonomous plan run for bug change set 01M767125EJM98D7H4ME2T67H5. No owner was present. The planner read the bug intent, internal/api/daemon_projects.go daemonMilestoneWaits, internal/api/request.go daemonPlanRequests' milestoneLive, the plan-milestone workflow (interview start step, then autonomous agent and command steps), run.State (Current, Stopped, SignoffWait) and the inbox's interview item (already ends when the interview step ends). Two questions went to the arbiter: which intent the change set serves, and which live milestone plan runs block a queued milestone request. The arbiter's Decision milestone-wait-owner-steps-only: serve milestone-25; block only on an unjournaled start within daemonStartGrace, or a run whose current step is owner-interactive (the interview start step, or a step paused waiting on the owner such as a gate or sign-off under owner autonomy); maxRuns still caps every start. The planner folded the rule into the existing leaf milestone-request-starts-at-once, whose statement said any live milestone run keeps the next one queued.

## title

Plan of Fix A queued milestone request (autonomous)

## Links

### In

- derived-from [milestone-wait-owner-steps-only](../decision/milestone-wait-owner-steps-only.md)
- derived-from [milestone-wait-owner-steps](../capability/rhiza/orchestration/daemon/plan-requests/milestone-wait-owner-steps.md)

## History

- `01M767125EJM98D7H4ME2T67H5` Fix A queued milestone request waits for another milestone's whole plan run to end, not just its owner interview, so the owner can't start the next interview for an hour or more after finishing one. Evidence 2026-11-11: milestone request 01M75YW54KYRCF7TMNHSETKTEM (release channels) was queued at 18:45Z while plan run 01M75W8YT6M1F8W05GAJZSEC0H (milestone 32, workflow cleanup) was at its interview; the owner finished that interview at about 20:19Z and the run moved on to milestone-check (then review, PR, CI, merge), but the release request stays queued because daemonMilestoneWaits (internal/api/daemon_projects.go:423, and daemonPlanRequests' milestoneLive) treats any live milestone plan run as blocking. Raising rhiza.daemon.maxRuns does not help. The one-at-a-time rule exists so the owner holds one interview at a time; it should count only milestone plan runs whose current step is an owner-interactive step (the interview, or a gate waiting on the owner), so the next milestone interview can start as soon as the previous one's interview step has passed. Test: a testscript with two milestone requests where the first run is past its interview step and the second starts.
