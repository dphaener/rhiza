# owner-run-gate-plan

source · `01M6PR0PP63YAFMJWV8M41DNNK`

## quotes

```json
[]
```

## raw_sha256

7cb4f7da8296a18402e86ada7b0dde7503633cb4380b66b256fe719e3b9a6be3

## session

627fef4e-be7d-42fa-908e-01a1ad7d72e6

## summary

Autonomous plan run, no owner. The planner read the gate (internal/api/run_review.go runRunGate: a non-autonomous run proposes its plan's drafts and returns before any accept) and the Inbox (internal/api/inbox.go: accept items only for autonomous runs' gated proposals). It raised three questions to the arbiter: which intent the change set serves, which rule an owner-run gate accepts under when rhiza.autonomy.accept is auto, and whether that accept needs the plan's review to have passed. The arbiter's Decision owner-run-gate-auto-accept: serve milestone-25; accept under a new rule 'autonomy', allowed only while rhiza.autonomy.accept is auto, with 'autonomous' keeping its meaning; accept only after the review has passed by gateReview's check, else keep the proposals waiting as with accept owner. The planner wrote three leaves: the owner-run gate accepts under rule autonomy, rule autonomy itself, and the Inbox listing any run's gated proposals that still wait.

## title

Plan of Fix An owner-run plan's gate (autonomous, arbiter)

## Links

### In

- derived-from [owner-run-gate-auto-accept](../decision/owner-run-gate-auto-accept.md)

## History

- `01M6PQTKN5N2T2S3MGCK00YZMQ` Fix An owner-run plan's gate leaves its change sets proposed even when rhiza.autonomy.accept is auto, and the Inbox doesn't list them (it lists proposals only when accept is owner), so a planned milestone silently stalls. Seen 2026-11-05: milestone 28's five web UI change sets sat proposed after plan PR #309 merged; the Inbox said nothing waits. Owner: "Yes those should never have been waiting for my acceptance". Milestone planning needs the owner (interview), but acceptance follows rhiza.autonomy.accept: with accept auto, the gate accepts what an owner-run plan proposed, as it does for autonomous plans; with accept owner, the proposals are listed in the Inbox as now.
