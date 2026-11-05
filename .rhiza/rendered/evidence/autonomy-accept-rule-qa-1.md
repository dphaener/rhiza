# autonomy-accept-rule-qa-1

evidence · `01M6PX8Q8P3G0M0GT56EK34X2T` · actor [claude-code-145392df](../actor/claude-code-145392df.md)

## result

PASS

## statement

At 5b3ccbe9 with a scratch build (go build -trimpath) in a fresh scratch repo, RHIZA_* unset, hold workflow keeping runs unfinished: rhiza accept --rule autonomy inside an unfinished normal run with rhiza.autonomy.accept owner was refused, exit 2, 'rhiza.autonomy.accept is owner, so only the owner accepts'; with accept auto but RHIZA_RUN unset refused, exit 2, 'only for a move made inside a run, and RHIZA_RUN is not set'; rule autonomous in the normal run still refused, exit 2, 'is not autonomous'; the gate's accept under accept auto recorded rule autonomy. autonomy-accept-rule testscript passes.

## Links

### Out

- verifies [autonomy-accept-rule](../capability/rhiza/workflows/autonomous-runs/autonomy-accept-rule.md)

## History

- `01M6PQTKN5N2T2S3MGCK00YZMQ` Fix An owner-run plan's gate leaves its change sets proposed even when rhiza.autonomy.accept is auto, and the Inbox doesn't list them (it lists proposals only when accept is owner), so a planned milestone silently stalls. Seen 2026-11-05: milestone 28's five web UI change sets sat proposed after plan PR #309 merged; the Inbox said nothing waits. Owner: "Yes those should never have been waiting for my acceptance". Milestone planning needs the owner (interview), but acceptance follows rhiza.autonomy.accept: with accept auto, the gate accepts what an owner-run plan proposed, as it does for autonomous plans; with accept owner, the proposals are listed in the Inbox as now.
