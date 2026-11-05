# owner-run-gate-accepts-verified-2

evidence · `01M6PY6N692EM95TK77FE46GS9` · actor [claude-code-72caff76](../actor/claude-code-72caff76.md)

## result

PASS

## statement

Verified attempt 2 at 8ff443e5 (fix for blocker f1). owner-run-gate-accept.txtar covers the owner-run gate: accept auto with a clean review accepts both change sets under rule autonomy and journals one accepted entry each (and a re-run adds none); with accept owner, and with accept auto but a confirmed finding that has no question, both stay proposed and rhiza accept --rule owner is next. The fix commit updates plan-workflow, plan-workflow-pr and plan-workflow-open to expect accepted/rule autonomy and an accepted Gate line in the PR body. The rhiza.autonomy.accept comment now says it governs every plan run's gate. Ran go test -trimpath ./cmd/rhiza -run TestScripts/ on owner-run-gate-accept, run-gate, plan-workflow, plan-workflow-pr, plan-workflow-open, plan-workflow-merge-main-moved, plan-workflow-autonomous, plan-workflow-autonomous-left, autonomous-gate-rule, autonomy-accept-rule, stage-autonomy, tui-autonomy-app, describe-golden, inbox-owner-run-accept and every other script that starts a plan run (18 more): all PASS. The autonomous gate scripts are unchanged and pass. go test -trimpath ./internal/api/... passes.

## Links

### Out

- verifies [owner-run-gate-accepts](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-gate/owner-run-gate-accepts.md)

## History

- `01M6PQTKN5N2T2S3MGCK00YZMQ` Fix An owner-run plan's gate leaves its change sets proposed even when rhiza.autonomy.accept is auto, and the Inbox doesn't list them (it lists proposals only when accept is owner), so a planned milestone silently stalls. Seen 2026-11-05: milestone 28's five web UI change sets sat proposed after plan PR #309 merged; the Inbox said nothing waits. Owner: "Yes those should never have been waiting for my acceptance". Milestone planning needs the owner (interview), but acceptance follows rhiza.autonomy.accept: with accept auto, the gate accepts what an owner-run plan proposed, as it does for autonomous plans; with accept owner, the proposals are listed in the Inbox as now.
