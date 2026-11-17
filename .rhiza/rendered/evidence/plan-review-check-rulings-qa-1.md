# plan-review-check-rulings-qa-1

evidence · `01M7MSTQEG74KCXAWZNXJ38SKY` · actor [claude-code-08e6675b](../actor/claude-code-08e6675b.md)

## result

PASS

## statement

QA: in the built binary, rhiza run review-check --rulings from the apply-rulings step exited 1 before the reviser linked anything, listing both questions with their Decision and the finding's node (export-csv, keep-flat); after the reviser linked each node decided-by its Decision it exited 0, listing both rulings with their decided-by node IDs, and --json gave a rulings array. Unset or unknown RHIZA_RUN and a finished run are refused with exit 2. go test ./internal/api -run TestRunReviewCheckRulings passes.

## Links

### Out

- verifies [plan-review-check-rulings](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-review-check-rulings.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
