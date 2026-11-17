# question-013bd336-1

question · `01M7MR2HXJGK247YZV324525NC`

## session

013bd336-c320-4bd6-b37b-01a1e9730922

## statement

Leaf plan-workflow-applies-rulings is blocked (question ruling-on-serving-draft-node-unappliable). In an autonomous plan run, the arbiter writes its Decision into the run's intention change set, but a finding on a leaf targets a node in the serving draft change set. Neither change set's view holds both nodes, so `rhiza apply` refuses the decided-by link either way (not_found). `rhiza run review-check --rulings` can then never pass for a leaf ruling, so every real plan run with a leaf ruling would end failed. How should this change set resolve it?

## Links

### In

- answers [review-ruling-in-node-change-set](../decision/review-ruling-in-node-change-set.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
