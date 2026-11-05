# milestone-plan-refuses-autonomous-verified

evidence · `01M6N3FQPZVRM3K01REKABQS35` · actor [claude-code-376a1866](../actor/claude-code-376a1866.md)

## result

PASS

## statement

At 6f24552b ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-start-milestone-autonomous|plan-workflow-autonomous|plan-workflow-autonomous-left|arbiter-agent-tabs|describe-golden)$': all PASS; go vet ./internal/api clean. run-start-milestone-autonomous shows run start --workflow plan --autonomous with --new --milestone and on a record change set each exits 2 with 'milestone planning needs the owner', leaving the change set count unchanged and no .git/rhiza/runs; the same record start without --autonomous succeeds, an --new intention with --autonomous succeeds with autonomous true, and run resume of the record run's journal edited to autonomous exits 2 with the same error. plan-workflow-autonomous and -left now start --new intention runs (intent without milestone=1 under a landed ms milestone; -left's keep-flat Decision landed with --record before the run) and still check the arbiter's answers, review rounds, the check's questions in the draft and the gate accepting with rule autonomous. Diff read: internal/api/run.go refuses before newPlanChange for --milestone (which milestone requests pass), after resolving for a record change set, and in runRunResume for an autonomous journal on a record change set.

## Links

### Out

- verifies [milestone-plan-refuses-autonomous](../capability/rhiza/workflows/autonomous-runs/milestone-plan-refuses-autonomous.md)

## History

- `01M6MY8BMPFK9ZWD3TDQHKT64J` Milestone planning always needs a human: a plan run for a milestone (a record change set, or a request of kind milestone) never runs autonomously, even when the daemon or run is autonomous (rhiza.autonomy.plan=auto); it waits for the owner's interview instead of the arbiter answering. Owner, 2026-11-05: "milestone planning should never be autonomous even in autonomous mode. that always requires a human."
