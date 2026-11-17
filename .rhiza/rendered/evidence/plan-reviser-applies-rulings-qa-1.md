# plan-reviser-applies-rulings-qa-1

evidence · `01M7MSTQHMRV4574B0XD753HRV` · actor [claude-code-08e6675b](../actor/claude-code-08e6675b.md)

## result

PASS

## statement

QA: rhiza guide plan-reviser from the built binary prints 'Applying the arbiter's rulings' under Must do, before Command budget, naming rhiza run status --json, the three decided-by cases, --record, never editing the Decision, and review-check --rulings. go test ./internal/guide -run TestPlanReviserText passes.

## Links

### Out

- verifies [plan-reviser-applies-rulings](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-reviser-applies-rulings.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
