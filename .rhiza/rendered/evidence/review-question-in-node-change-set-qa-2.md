# review-question-in-node-change-set-qa-2

evidence · `01M7MT0GGSWSA4XJ9CPWNCFAKM` · actor [claude-code-fc7b7aa0](../actor/claude-code-fc7b7aa0.md)

## result

PASS

## statement

QA round 2 at db1b1c925, scratch build (go build -trimpath -o $SCRATCH/rhiza-x ./cmd/rhiza): the plan-workflow-autonomous-left/-unapplied scripts pass, and by hand a Decision in the open change set links from a landed node in that same change set (apply exit 0).

## Links

### Out

- verifies [review-question-in-node-change-set](../capability/review-question-in-node-change-set.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
