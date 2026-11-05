# inbox-lists-owner-run-proposals

capability · `01M6PR136HA6VBMENME6VDFANQ` · parent [inbox](../inbox.md)

## acceptance

A testscript (extend inbox-list.txtar or add inbox-owner-run-accept.txtar) gates a non-autonomous plan run with rhiza.autonomy.accept owner and shows rhiza inbox listing each proposed change set as an accept item with next rhiza accept --rule owner --change <id>; once accepted, the item is gone; existing autonomous accept items still list

## statement

rhiza inbox lists each change set any plan run's gate proposed that is still proposed as an accept item waiting for acceptance, owner-run plan runs as well as autonomous ones, so a planned milestone never stalls unseen

## Links

### Out

- decided-by [owner-run-gate-auto-accept](../../../../decision/owner-run-gate-auto-accept.md)

### In

- verifies [inbox-lists-owner-run-proposals-verified](../../../../evidence/inbox-lists-owner-run-proposals-verified.md)
- verifies [inbox-lists-owner-run-proposals-qa-1](../../../../evidence/inbox-lists-owner-run-proposals-qa-1.md)
- verifies [pr-315-fix-an-owner-run-plan-s-gate](../../../../evidence/pr-315-fix-an-owner-run-plan-s-gate.md)

## History

- `01M6PQTKN5N2T2S3MGCK00YZMQ` Fix An owner-run plan's gate leaves its change sets proposed even when rhiza.autonomy.accept is auto, and the Inbox doesn't list them (it lists proposals only when accept is owner), so a planned milestone silently stalls. Seen 2026-11-05: milestone 28's five web UI change sets sat proposed after plan PR #309 merged; the Inbox said nothing waits. Owner: "Yes those should never have been waiting for my acceptance". Milestone planning needs the owner (interview), but acceptance follows rhiza.autonomy.accept: with accept auto, the gate accepts what an owner-run plan proposed, as it does for autonomous plans; with accept owner, the proposals are listed in the Inbox as now.
