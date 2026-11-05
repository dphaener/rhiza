# autonomy-accept-rule-verified

evidence · `01M6PVZVASM6TKS5CPQZGCV5B6` · actor [claude-code-a3c54cf1](../actor/claude-code-a3c54cf1.md)

## result

PASS

## statement

At b14779d5 ran go test -trimpath -count=1 -v ./cmd/rhiza -run 'TestScripts/(autonomy-accept-rule|autonomous-gate-rule|owner-run-gate-accept|stage-autonomy|run-gate)$': all PASS. autonomy-accept-rule.txtar shows rhiza accept --rule autonomy with rhiza.autonomy.accept=owner refused as usage (code usage, exit 2, change set stays proposed), rule autonomous in the non-autonomous run still refused ('is not autonomous', exit 2), and with accept=auto the accept succeeds and change show records rule autonomy; fold-ok passes. Also ran go test -trimpath -count=1 ./internal/api -run 'TestLifecycle|...|Autonomy': ok. Read the diff: checkAutonomousRule in internal/api/lifecycle.go refuses rule autonomy via Usagef unless autonomyStage(accept) is auto, before the RHIZA_RUN check; rule autonomous is unchanged.

## Links

### Out

- verifies [autonomy-accept-rule](../capability/rhiza/workflows/autonomous-runs/autonomy-accept-rule.md)

## History

- `01M6PQTKN5N2T2S3MGCK00YZMQ` Fix An owner-run plan's gate leaves its change sets proposed even when rhiza.autonomy.accept is auto, and the Inbox doesn't list them (it lists proposals only when accept is owner), so a planned milestone silently stalls. Seen 2026-11-05: milestone 28's five web UI change sets sat proposed after plan PR #309 merged; the Inbox said nothing waits. Owner: "Yes those should never have been waiting for my acceptance". Milestone planning needs the owner (interview), but acceptance follows rhiza.autonomy.accept: with accept auto, the gate accepts what an owner-run plan proposed, as it does for autonomous plans; with accept owner, the proposals are listed in the Inbox as now.
