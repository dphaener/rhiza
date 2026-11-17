# plan-review-check-rulings-qa-4

evidence · `01M7MTS3ZR8B69B6ENK5FGCNZH` · actor [claude-code-a51c8304](../actor/claude-code-a51c8304.md)

## result

PASS

## statement

QA cycle 4, scratch binary (go build -trimpath) in a scratch autonomous plan run with fake agents (two round-3 findings: serving-draft leaf export-csv, landed Decision keep-flat): rhiza run review-check --rulings failed listing both question, Decision and node handles while unapplied, and passed ('2 arbiter ruling(s) applied') once each node linked decided-by its Decision; a run whose apply-rulings visit wrote nothing ended failed before the gate.

## Links

### Out

- verifies [plan-review-check-rulings](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-review-check-rulings.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
