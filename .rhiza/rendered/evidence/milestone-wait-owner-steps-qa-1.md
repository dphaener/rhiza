# milestone-wait-owner-steps-qa-1

evidence · `01M76MXS8C8ZXCFYBBA2XWBAEQ` · actor [claude-code-187fc5bc](../actor/claude-code-187fc5bc.md)

## result

PASS

## statement

QA at 42f7d662d with the built binary (go build -trimpath ./cmd/rhiza) in a scratch root with a fake rex and RHIZA_TEST_DETACH_LOG: queued two milestone requests; daemon --once --max-runs 5 started only the first's plan-milestone run; the next poll (start unjournaled, within grace) kept the second queued; with the first's journal live on another host at the interview step the second stayed queued; after appending entered review, --max-runs 1 started nothing and --max-runs 5 started the second's plan run (detach log: run start --workflow plan-milestone --new Second milestone ... --milestone). A third milestone request stayed queued while the first was stopped at the gate (not autonomous) and the second live at milestone-check, and started once the first's stop was at review instead. go test -trimpath ./internal/api -run MilestoneWaits and ./cmd/rhiza -run 'TestScripts/(daemon-milestone-request|plan-requests)' pass.

## Links

### Out

- verifies [milestone-wait-owner-steps](../capability/rhiza/orchestration/daemon/plan-requests/milestone-wait-owner-steps.md)

## History

- `01M767125EJM98D7H4ME2T67H5` Fix A queued milestone request waits for another milestone's whole plan run to end, not just its owner interview, so the owner can't start the next interview for an hour or more after finishing one. Evidence 2026-11-11: milestone request 01M75YW54KYRCF7TMNHSETKTEM (release channels) was queued at 18:45Z while plan run 01M75W8YT6M1F8W05GAJZSEC0H (milestone 32, workflow cleanup) was at its interview; the owner finished that interview at about 20:19Z and the run moved on to milestone-check (then review, PR, CI, merge), but the release request stays queued because daemonMilestoneWaits (internal/api/daemon_projects.go:423, and daemonPlanRequests' milestoneLive) treats any live milestone plan run as blocking. Raising rhiza.daemon.maxRuns does not help. The one-at-a-time rule exists so the owner holds one interview at a time; it should count only milestone plan runs whose current step is an owner-interactive step (the interview, or a gate waiting on the owner), so the next milestone interview can start as soon as the previous one's interview step has passed. Test: a testscript with two milestone requests where the first run is past its interview step and the second starts.
