# review-question-in-node-change-set-qa-4

evidence · `01M7MTS490E3XP6RYS39YZQ9JZ` · actor [claude-code-a51c8304](../actor/claude-code-a51c8304.md)

## result

PASS

## statement

QA cycle 4, same scratch plan run: the leaf finding's question and the arbiter's Decision landed in the serving draft (rhiza show --change <serving draft> shows Created: the serving draft); the landed keep-flat finding's question and Decision in the run's draft intention change set.

## Links

### Out

- verifies [review-question-in-node-change-set](../capability/review-question-in-node-change-set.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
