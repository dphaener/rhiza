# pr-315-fix-an-owner-run-plan-s-gate

evidence · `01M6Q043E552F498VDBDB1G5RT` · actor [dphaener](../actor/dphaener.md)

## head

37a8ef2053acf61598cb923cc51aebc0b059b8fc

## pr

https://github.com/dphaener/rhiza/pull/315

## result

1 check passed: check SUCCESS

## statement

Pull request #315 (https://github.com/dphaener/rhiza/pull/315) passed CI at 37a8ef2

## status

verified

## Links

### Out

- verifies [owner-run-gate-accepts](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-gate/owner-run-gate-accepts.md)
- verifies [autonomy-accept-rule](../capability/rhiza/workflows/autonomous-runs/autonomy-accept-rule.md)
- verifies [inbox-lists-owner-run-proposals](../capability/rhiza/orchestration/inbox/inbox-lists-owner-run-proposals.md)
- verifies [plan-gate](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-gate.md)

## History

- `01M6PQTKN5N2T2S3MGCK00YZMQ` Fix An owner-run plan's gate leaves its change sets proposed even when rhiza.autonomy.accept is auto, and the Inbox doesn't list them (it lists proposals only when accept is owner), so a planned milestone silently stalls. Seen 2026-11-05: milestone 28's five web UI change sets sat proposed after plan PR #309 merged; the Inbox said nothing waits. Owner: "Yes those should never have been waiting for my acceptance". Milestone planning needs the owner (interview), but acceptance follows rhiza.autonomy.accept: with accept auto, the gate accepts what an owner-run plan proposed, as it does for autonomous plans; with accept owner, the proposals are listed in the Inbox as now.
