# owner-run-gate-accepts-qa-1

evidence · `01M6PX8Q669Z2909TVHB5K70XX` · actor [claude-code-145392df](../actor/claude-code-145392df.md)

## result

PASS

## statement

At 5b3ccbe9 with a scratch build (go build -trimpath) in a fresh scratch repo, RHIZA_* unset, hold workflow keeping runs unfinished: a normal (non-autonomous) plan run with a clean review round and accept auto (default): RHIZA_RUN=<run> rhiza run gate exited 0, said 'Gated the plan of normal run' and both served change sets went proposed then accepted; change show records rule autonomy, reason 'rhiza.autonomy.accept is auto, run <id>'. With git config rhiza.autonomy.accept owner the same gate exited 0, left the change set proposed with next rhiza accept --rule owner --change <id>. With accept auto and no review round at all, the gate exited 0 and left it proposed with the owner accept next. rhiza check passed after. Scoped testscripts owner-run-gate-accept, run-gate, autonomous-gate-rule, stage-autonomy pass.

## Links

### Out

- verifies [owner-run-gate-accepts](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-gate/owner-run-gate-accepts.md)

## History

- `01M6PQTKN5N2T2S3MGCK00YZMQ` Fix An owner-run plan's gate leaves its change sets proposed even when rhiza.autonomy.accept is auto, and the Inbox doesn't list them (it lists proposals only when accept is owner), so a planned milestone silently stalls. Seen 2026-11-05: milestone 28's five web UI change sets sat proposed after plan PR #309 merged; the Inbox said nothing waits. Owner: "Yes those should never have been waiting for my acceptance". Milestone planning needs the owner (interview), but acceptance follows rhiza.autonomy.accept: with accept auto, the gate accepts what an owner-run plan proposed, as it does for autonomous plans; with accept owner, the proposals are listed in the Inbox as now.
