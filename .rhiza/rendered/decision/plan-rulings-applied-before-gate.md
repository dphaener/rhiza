# plan-rulings-applied-before-gate

decision · `01M7MPK6PCG649WAMS3GN7ST9G`

## reasoning

Runs 01M7GG7GXSA97HESRNFEBFCZ4T and 01M7GG7HSBSB356BVGGQGH8YHB went check -> gate right after the arbiter wrote its Decision (internal/workflow/builtin/plan.toml), so the accepted plan still carried the confirmed major and a builder later applied the ruling by hand. The arbiter guide forbids the arbiter writing anything but its Decision, so applying it is the reviser's job, as revising is for rounds 1 and 2. The intent offers the run failing instead of accepting as the alternative; failing after one apply visit gives both. A re-review is not added: review is capped at 3 visits and the rulings check is mechanical.

## statement

In an autonomous plan run, a confirmed major left after review round 3 is decided by the arbiter and then applied to the plan before the gate: the check step passes to a rulings step (rhiza run review-check --rulings), which fails while a passed ruling is unapplied and sends the run to one apply-rulings visit of the plan reviser; a ruling still unapplied after it ends the run failed before the gate, so the gate never accepts a plan whose confirmed major the arbiter ruled on but nobody applied. An arbiter that cannot decide still fails the run naming the question. A normal run, whose review questions stay open for the owner, is unchanged.

## Links

### Out

- supersedes [plan-review-left-arbitrated](plan-review-left-arbitrated.md)
- derived-from [plan-run-rulings-research](../source/plan-run-rulings-research.md)

### In

- decided-by [plan-review-check-rulings](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-review-check-rulings.md)
- decided-by [plan-workflow-applies-rulings](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-workflow-applies-rulings.md)
- decided-by [plan-reviser-applies-rulings](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-reviser-applies-rulings.md)
- decided-by [review-question-in-node-change-set](../capability/review-question-in-node-change-set.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
