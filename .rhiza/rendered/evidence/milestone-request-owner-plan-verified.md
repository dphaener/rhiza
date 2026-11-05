# milestone-request-owner-plan-verified

evidence · `01M6N3MSBJHH9ZKCQ5XK8FWZ0S` · actor [claude-code-84889f26](../actor/claude-code-84889f26.md)

## result

PASS

## statement

At 33055a5c, go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/plan-requests$' passes. With a fake rex and rhiza.autonomy.plan=auto, the daemon pass detaches the milestone request 'Ship v2' as run start ... --milestone (line ending there, and ! stdout 'Ship v2.*--autonomous'), the second feature request 'Add search' as ... --autonomous and the bug as ... --autonomous. Starting the milestone run as the daemon said leaves its journal's started entry with no "autonomous":true (count 0). I confirmed in a scratch copy of the script that request 3 is the milestone and that an autonomous run's journal does contain "autonomous":true, so the 0 check means something. The diff in internal/api/request.go daemonStartPlan sets autonomous=false for KindMilestone only, so feature and bug requests still follow rhiza.autonomy.plan.

## Links

### Out

- verifies [milestone-request-owner-plan](../capability/rhiza/orchestration/daemon/plan-requests/request-queue/milestone-request-owner-plan.md)

## History

- `01M6MY8BMPFK9ZWD3TDQHKT64J` Milestone planning always needs a human: a plan run for a milestone (a record change set, or a request of kind milestone) never runs autonomously, even when the daemon or run is autonomous (rhiza.autonomy.plan=auto); it waits for the owner's interview instead of the arbiter answering. Owner, 2026-11-05: "milestone planning should never be autonomous even in autonomous mode. that always requires a human."
