# review-question-in-node-change-set-qa-3

evidence · `01M7MTBDXVDJ6QD3HA389X58SY` · actor [claude-code-60928770](../actor/claude-code-60928770.md)

## result

PASS

## statement

QA round 3 at db1b1c925, scratch build (go build -trimpath, run by full path, outside PATH) driving autonomous plan runs by hand in scratch Rhiza roots with the plan-workflow-autonomous* fake rex/claude/gh: with the run's change set open, round 3's question on export-csv (a serving draft's leaf) went into the draft and its asked journal entry has change=<draft>; the one on landed keep-flat stayed in the run's change set with no change field; the fake arbiter logged RHIZA_CHANGE=<draft> and <run's change set> respectively and each Decision was written there.

## Links

### Out

- verifies [review-question-in-node-change-set](../capability/review-question-in-node-change-set.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
