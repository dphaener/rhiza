# plan-workflow-applies-rulings-qa-4

evidence · `01M7MTS441Q66JMNPJX4RVA914` · actor [claude-code-a51c8304](../actor/claude-code-a51c8304.md)

## result

PASS

## statement

QA cycle 4, same scratch plan run: steps went check (visit 3) pass then rulings; rulings (visit 1) fail then apply-rulings; apply-rulings (visit 1) pass then rulings; rulings (visit 2) pass then gate; gate pass. With a reviser that wrote nothing, rulings visit 2 failed and the run failed with step apply-rulings over max_visits 1, no gate entered.

## Links

### Out

- verifies [plan-workflow-applies-rulings](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-workflow-applies-rulings.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
