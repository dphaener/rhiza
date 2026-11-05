# inbox-lists-owner-run-proposals-qa-1

evidence · `01M6PX8QARTZHE111DYVG0RTGP` · actor [claude-code-145392df](../actor/claude-code-145392df.md)

## result

PASS

## statement

At 5b3ccbe9 with a scratch build (go build -trimpath) in a fresh scratch repo, RHIZA_* unset, hold workflow keeping runs unfinished: after a normal plan run's gate with accept owner left a change set proposed, rhiza inbox listed it as 'accept run <run>, change set <id> / waiting for acceptance / next: rhiza accept --rule owner --change <id>'; running that next accepted it (Rule: owner) and the item left the inbox. A normal run whose gate left a proposal for want of a review round (accept auto) also listed it. inbox-owner-run-accept testscript passes.

## Links

### Out

- verifies [inbox-lists-owner-run-proposals](../capability/rhiza/orchestration/inbox/inbox-lists-owner-run-proposals.md)

## History

- `01M6PQTKN5N2T2S3MGCK00YZMQ` Fix An owner-run plan's gate leaves its change sets proposed even when rhiza.autonomy.accept is auto, and the Inbox doesn't list them (it lists proposals only when accept is owner), so a planned milestone silently stalls. Seen 2026-11-05: milestone 28's five web UI change sets sat proposed after plan PR #309 merged; the Inbox said nothing waits. Owner: "Yes those should never have been waiting for my acceptance". Milestone planning needs the owner (interview), but acceptance follows rhiza.autonomy.accept: with accept auto, the gate accepts what an owner-run plan proposed, as it does for autonomous plans; with accept owner, the proposals are listed in the Inbox as now.
