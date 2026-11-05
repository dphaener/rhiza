# owner-run-gate-accepts-verified

evidence · `01M6PX3KFYE3XDW7Y32FGKZ4WC` · actor [claude-code-177ad7c7](../actor/claude-code-177ad7c7.md)

## result

PASS

## statement

At d9367578, go test -trimpath ./cmd/rhiza -count=1 -run TestScripts/(owner-run-gate-accept|run-gate|autonomous-gate-rule|stage-autonomy|plan-workflow-autonomous|autonomous-review) passes. In the new owner-run-gate-accept.txtar, a non-autonomous hold run's gate with accept auto and a clean review round accepts both proposed change sets under rule autonomy, and the run journal has 2 accepted entries; a second gate adds none. With accept owner, and with accept auto but a confirmed finding no question was left for, the gate exits 0, leaves the change set proposed with next rhiza accept --rule owner --change <id>, and journals no accepted entry. The autonomous gate scripts are unchanged and still pass. In the diff, runRunGate keeps the autonomous path (failure when the review has not passed, open questions stopping the run, rule autonomous); an owner run uses rule autonomy, which checkAutonomousRule allows only while RHIZA_RUN names an unfinished run; the rhiza.autonomy.accept setting's comment now says it governs every plan run's gate. go test -trimpath -timeout 30m -skip 'TestReadBudget|TestWriteBudget' ./internal/api/... passes. The full package run failed only TestReadBudget, a wall-clock read budget, on a machine with load average 6-9 from parallel runs. This diff changes no read path, so that failure is load, not this change; CI's make check covers the budget.

## Links

### Out

- verifies [owner-run-gate-accepts](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-gate/owner-run-gate-accepts.md)

## History

- `01M6PQTKN5N2T2S3MGCK00YZMQ` Fix An owner-run plan's gate leaves its change sets proposed even when rhiza.autonomy.accept is auto, and the Inbox doesn't list them (it lists proposals only when accept is owner), so a planned milestone silently stalls. Seen 2026-11-05: milestone 28's five web UI change sets sat proposed after plan PR #309 merged; the Inbox said nothing waits. Owner: "Yes those should never have been waiting for my acceptance". Milestone planning needs the owner (interview), but acceptance follows rhiza.autonomy.accept: with accept auto, the gate accepts what an owner-run plan proposed, as it does for autonomous plans; with accept owner, the proposals are listed in the Inbox as now.
