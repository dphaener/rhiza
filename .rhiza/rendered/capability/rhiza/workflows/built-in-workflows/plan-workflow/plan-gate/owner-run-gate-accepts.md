# owner-run-gate-accepts

capability · `01M6PR133TKYZKSJRVSXNQ4Y9T` · parent [plan-gate](../plan-gate.md)

## acceptance

cmd/rhiza/testdata/script/run-gate.txtar (or a new owner-run-gate-accept.txtar) shows a non-autonomous plan run's gate with accept auto and a passed review accepting each proposed change set under rule autonomy and journalling an accepted entry each; with accept owner, and with accept auto but a confirmed finding left without a question, the same gate leaves them proposed with the owner accept next; autonomous-run gate behaviour is unchanged; go test -trimpath ./internal/api/... passes

## statement

An owner-run (non-autonomous) plan run's gate, after proposing its plan's drafts, accepts each proposed change set under rule autonomy when rhiza.autonomy.accept is auto and the plan's review has passed by gateReview's check; when accept is owner or the review has not passed, it accepts nothing and leaves each proposal with rhiza accept --rule owner as its next step; the rhiza.autonomy.accept setting's description says it governs every plan run's gate

## Links

### Out

- decided-by [owner-run-gate-auto-accept](../../../../../../decision/owner-run-gate-auto-accept.md)

### In

- verifies [owner-run-gate-accepts-verified](../../../../../../evidence/owner-run-gate-accepts-verified.md)
- verifies [owner-run-gate-accepts-qa-1](../../../../../../evidence/owner-run-gate-accepts-qa-1.md)
- verifies [owner-run-gate-accepts-verified-2](../../../../../../evidence/owner-run-gate-accepts-verified-2.md)
- verifies [pr-315-fix-an-owner-run-plan-s-gate](../../../../../../evidence/pr-315-fix-an-owner-run-plan-s-gate.md)

## History

- `01M6PQTKN5N2T2S3MGCK00YZMQ` Fix An owner-run plan's gate leaves its change sets proposed even when rhiza.autonomy.accept is auto, and the Inbox doesn't list them (it lists proposals only when accept is owner), so a planned milestone silently stalls. Seen 2026-11-05: milestone 28's five web UI change sets sat proposed after plan PR #309 merged; the Inbox said nothing waits. Owner: "Yes those should never have been waiting for my acceptance". Milestone planning needs the owner (interview), but acceptance follows rhiza.autonomy.accept: with accept auto, the gate accepts what an owner-run plan proposed, as it does for autonomous plans; with accept owner, the proposals are listed in the Inbox as now.
