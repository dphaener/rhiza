# ruling-applied-means-decided-by

decision · `01M7MPK6RXY3HBZHQ2EJM5JV5D`

## reasoning

The model has no other mechanical sign that a node followed a Decision; decided-by is the link every leaf uses for its Decisions and modellint.UnappliedDecisions already reads it the same way. Field-diff checks would need to guess which field the ruling changes. The finding's node can be a Decision in a landed record (plan-workflow-autonomous-left.txtar's keep-flat), so the check accepts any live node kind, written in any change set.

## statement

A ruling counts as applied when a live node links decided-by to the arbiter's Decision. The arbiter writes only the answers link to the question, so the reviser who changes the finding's node per the ruling links that node decided-by the Decision; the link is what rhiza run review-check --rulings tests. Only rulings whose arbitration passed count; open questions (a normal run) and failed arbitrations (which already fail the run) are not looked at.

## Links

### Out

- derived-from [plan-run-rulings-research](../source/plan-run-rulings-research.md)

### In

- decided-by [plan-review-check-rulings](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-review-check-rulings.md)
- decided-by [plan-reviser-applies-rulings](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-reviser-applies-rulings.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
