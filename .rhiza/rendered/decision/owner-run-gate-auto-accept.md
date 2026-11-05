# owner-run-gate-auto-accept

decision · `01M6PQZSD9GN7MMRNFBNK03TGE`

## reasoning

(1) The change set's intent is a bug/polish fix that unblocks the daemon building planned work, which is milestone-25's 'bugs and polish' scope; the question itself notes recent Fix change sets serve it. (2) internal/api/lifecycle.go documents rule autonomous as 'always means such a run made the move' and guards it on RHIZA_RUN naming an unfinished autonomous run; reusing it for an owner-run plan would break that invariant, and rule 'owner' would falsely claim the owner accepted (the owner's point was these 'should never have been waiting for my acceptance'). A distinct rule naming the setting keeps the record honest about who accepted and why, and mirrors the autonomous guard. (3) The plan workflow runs the capped plan review in normal runs as well (evidence plan-review-arbitrated-verified: a normal run leaves review questions open rather than arbitrating them), so gateReview has a round to check. Accepting without the owner should carry the same safety as the autonomous gate; accepting an unreviewed plan or one with open review questions would bypass the review the owner relies on. Reusing gateReview keeps one rule for 'may the gate accept without the owner'.

## statement

(1) Change set 'Fix An owner-run plan's gate' serves milestone-25, as the other recent Fix change sets do. (2) When an owner-run (non-autonomous) plan's gate accepts because rhiza.autonomy.accept is auto, it accepts under a new rule 'autonomy' (the owner's standing rhiza.autonomy.accept setting accepted it), not rule 'autonomous' and not rule 'owner'; rule 'autonomy' is allowed only while rhiza.autonomy.accept is auto, and 'autonomous' keeps its meaning (an unfinished autonomous run made the move). (3) Yes: the owner-run gate auto-accepts only after the plan's review has passed by the same check as the autonomous gate (gateReview: a review round exists, every confirmed finding has a question, no review question open); if it has not passed, it accepts nothing, keeps the proposals, and they wait as they do with accept owner.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-627fef4e-1](../question/question-627fef4e-1.md)
- derived-from [owner-run-gate-plan](../source/owner-run-gate-plan.md)

### In

- decided-by [owner-run-gate-accepts](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-gate/owner-run-gate-accepts.md)
- decided-by [autonomy-accept-rule](../capability/rhiza/workflows/autonomous-runs/autonomy-accept-rule.md)
- decided-by [inbox-lists-owner-run-proposals](../capability/rhiza/orchestration/inbox/inbox-lists-owner-run-proposals.md)

## History

- `01M6PQTKN5N2T2S3MGCK00YZMQ` Fix An owner-run plan's gate leaves its change sets proposed even when rhiza.autonomy.accept is auto, and the Inbox doesn't list them (it lists proposals only when accept is owner), so a planned milestone silently stalls. Seen 2026-11-05: milestone 28's five web UI change sets sat proposed after plan PR #309 merged; the Inbox said nothing waits. Owner: "Yes those should never have been waiting for my acceptance". Milestone planning needs the owner (interview), but acceptance follows rhiza.autonomy.accept: with accept auto, the gate accepts what an owner-run plan proposed, as it does for autonomous plans; with accept owner, the proposals are listed in the Inbox as now.
