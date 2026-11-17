# review-question-in-node-change-set-qa-1

evidence · `01M7MSTQP6ANVZV4647S6VKMH8` · actor [claude-code-08e6675b](../actor/claude-code-08e6675b.md)

## result

PASS

## statement

QA: in a hand-run autonomous plan run with the built binary, the arbiter for the round-3 question on export-csv (a leaf the serving draft created) started with RHIZA_CHANGE set to the serving draft, and the one on the landed Decision keep-flat started in the run's open change set. go test ./internal/run -run 'TestCommandAsked|TestArbitrate' and the run-review-check*, run-arbitrate-asked scripts pass.

## Links

### Out

- verifies [review-question-in-node-change-set](../capability/review-question-in-node-change-set.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
