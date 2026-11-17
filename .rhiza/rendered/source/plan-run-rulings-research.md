# plan-run-rulings-research

source · `01M7MPGTKQMY2H0GF1R82ZQXFH`

## quotes

```json
[]
```

## raw_sha256

8d513a0bb149d482a4979129312908b54d3b9a7104c07abb2e350782c5ba722e

## session

c3cfdf03-15ba-4ae5-a7b3-01a1e9658540

## summary

Autonomous plan run (RHIZA_AUTONOMOUS=1, no owner). Research pass on runs 01M7GG7GXSA97HESRNFEBFCZ4T and 01M7GG7HSBSB356BVGGQGH8YHB. Findings: after review round 3, rhiza run review-check (internal/api/run_review.go runRunReviewCheck) adds a question per confirmed major and journals a session-less asked entry; the engine (internal/run/engine.go arbitrateAsked) has the arbiter decide it after the check step, then the plan workflow (internal/workflow/builtin/plan.toml) goes check -> gate, so the gate accepts the plan unchanged. The arbiter's Decision is written and merged (parity-seed-bridge-in-model landed with the plan in #598), but it only links answers to the question: no plan node is changed or linked decided-by to it, so the accepted plan still carries the confirmed major (a later builder applied it by hand, entry 01M7GJA0T2EJSPHZQA08QTX19H). The intent's claim of no Decision node does not hold; the unapplied ruling does. Run status prints 'decision X, not yet delivered' for these questions (internal/api/run.go markdown, Answered is set only when the engine delivers to an asking session), but a session-less question has no session to deliver to, so the line misleads. The plan-reviser guide (internal/guide/text/plan-reviser.md) covers only confirmed findings, not rulings. Planned, as the planner's own choices from this reading: a review-check --rulings check that a passed ruling is applied (a live node links decided-by to its Decision), a rulings step before the gate with an apply-rulings reviser visit, failing the run when a ruling stays unapplied; the reviser guide section; and the status line fix.

## title

Autonomous plan of the bug: plan-run arbiter rulings never applied

## Links

### In

- derived-from [plan-rulings-applied-before-gate](../decision/plan-rulings-applied-before-gate.md)
- derived-from [ruling-applied-means-decided-by](../decision/ruling-applied-means-decided-by.md)
- derived-from [sessionless-ruling-not-undelivered](../decision/sessionless-ruling-not-undelivered.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
