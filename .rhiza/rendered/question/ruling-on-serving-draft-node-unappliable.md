# ruling-on-serving-draft-node-unappliable

question · `01M7MQZE172EGW6NDZ6R3Z2PFZ`

## statement

How is an arbiter ruling applied when the finding's node lives in a serving draft change set and the arbiter's Decision lives in the run's intention change set? Neither change set's view holds both nodes, so rhiza apply refuses the decided-by link in either (not_found, by handle and by ID), and review-check --rulings can never pass for such a ruling. plan-workflow-autonomous-left.txtar hits exactly this: finding f1 is on leaf export-csv in the serving draft, its question and the arbiter's Decision settle-... land in the intention change set; f2 on landed keep-flat applies fine. In real autonomous plan runs leaves live in serving drafts, so every leaf ruling would end the run failed. Options: (a) the arbiter writes its Decision in the change set holding the finding's node; (b) review-check leaves the question there; (c) --rulings accepts another witness; (d) narrow this leaf's script to rulings on nodes the run's change set sees. The workflow steps (plan.toml, plan-milestone.toml), the unit test and the five 'then rulings' script assertions are built and pass in the working tree, uncommitted.

## Links

### Out

- blocks [plan-workflow-applies-rulings](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-workflow-applies-rulings.md)

### In

- answers [review-ruling-in-node-change-set](../decision/review-ruling-in-node-change-set.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
