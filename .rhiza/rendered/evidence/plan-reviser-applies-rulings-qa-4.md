# plan-reviser-applies-rulings-qa-4

evidence · `01M7MTS4GVZJ30B9SWV33SW2XD` · actor [claude-code-a51c8304](../actor/claude-code-a51c8304.md)

## result

FAIL

## statement

QA cycle 4: followed rhiza guide plan-reviser 'Applying the arbiter's rulings' literally inside the apply-rulings visit of a scratch autonomous plan run (RHIZA_CHANGE = the run's draft change set). rhiza run status --json listed both questions with their decisions; rhiza apply --record changed the landed node; rhiza show <decision> named the run's change set for the landed-node ruling; but for the serving-draft leaf ruling rhiza show <decision> failed with 'no node or link matches', so the guide's way to find the change set does not work in that case.

## Links

### Out

- verifies [plan-reviser-applies-rulings](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-reviser-applies-rulings.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
