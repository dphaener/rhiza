# milestone-plan-refuses-autonomous-verified-2

evidence · `01M6N5DN8W3DD3P3QKKBBS4FWP` · actor [claude-code-84d76959](../actor/claude-code-84d76959.md)

## result

PASS

## statement

Verified at 4c99a781 (attempt 2). Read the diff: runRunStart refuses --autonomous on a plan run with --milestone, with --request naming a milestone-kind request (new: requestRun now returns the request and its Kind is checked before any draft opens, --milestone or not), or on a record change set (planIsMilestone); runRunResume refuses a plan run journaled autonomous on a record change set; all use milestoneNeedsOwner, a usage error saying milestone planning needs the owner. Ran go test -trimpath ./cmd/rhiza -run TestScripts/(run-start-milestone-autonomous|plan-workflow-autonomous|plan-workflow-autonomous-left|plan-requests|describe-golden) -v: all PASS. run-start-milestone-autonomous checks exit 2 with the rule for --new --milestone, a record change set, and a milestone request without --milestone (dry run and real), no new draft and no journal; the record start without --autonomous succeeds, an autonomous --new intention succeeds with autonomous true, and run resume of the record run edited to autonomous exits 2 with the rule. plan-workflow-autonomous and -left now start on an intention (--new, no --milestone), the fake claude asserts RHIZA_AUTONOMOUS=1, and the gate accepts with rule autonomous. go test -trimpath ./internal/api/... ok.

## Links

### Out

- verifies [milestone-plan-refuses-autonomous](../capability/rhiza/workflows/autonomous-runs/milestone-plan-refuses-autonomous.md)

## History

- `01M6MY8BMPFK9ZWD3TDQHKT64J` Milestone planning always needs a human: a plan run for a milestone (a record change set, or a request of kind milestone) never runs autonomously, even when the daemon or run is autonomous (rhiza.autonomy.plan=auto); it waits for the owner's interview instead of the arbiter answering. Owner, 2026-11-05: "milestone planning should never be autonomous even in autonomous mode. that always requires a human."
