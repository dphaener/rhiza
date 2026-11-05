# autonomy-accept-rule

capability · `01M6PR135E00DF7NHH46BYZF87` · parent [autonomous-runs](../autonomous-runs.md)

## acceptance

A testscript (extend autonomous-gate-rule.txtar or add autonomy-accept-rule.txtar) shows rhiza accept --rule autonomy refused as usage while rhiza.autonomy.accept is owner and accepted (the log entry records rule autonomy) while it is auto; rule autonomous outside an autonomous run is still refused as before

## statement

rhiza accept takes rule autonomy, recording that the owner's standing rhiza.autonomy.accept setting made the move; it is refused as usage unless rhiza.autonomy.accept is auto, and rule autonomous keeps meaning a move an unfinished autonomous run made

## Links

### Out

- decided-by [owner-run-gate-auto-accept](../../../../decision/owner-run-gate-auto-accept.md)

### In

- verifies [autonomy-accept-rule-verified](../../../../evidence/autonomy-accept-rule-verified.md)
- verifies [autonomy-accept-rule-qa-1](../../../../evidence/autonomy-accept-rule-qa-1.md)
- verifies [pr-315-fix-an-owner-run-plan-s-gate](../../../../evidence/pr-315-fix-an-owner-run-plan-s-gate.md)

## History

- `01M6PQTKN5N2T2S3MGCK00YZMQ` Fix An owner-run plan's gate leaves its change sets proposed even when rhiza.autonomy.accept is auto, and the Inbox doesn't list them (it lists proposals only when accept is owner), so a planned milestone silently stalls. Seen 2026-11-05: milestone 28's five web UI change sets sat proposed after plan PR #309 merged; the Inbox said nothing waits. Owner: "Yes those should never have been waiting for my acceptance". Milestone planning needs the owner (interview), but acceptance follows rhiza.autonomy.accept: with accept auto, the gate accepts what an owner-run plan proposed, as it does for autonomous plans; with accept owner, the proposals are listed in the Inbox as now.
