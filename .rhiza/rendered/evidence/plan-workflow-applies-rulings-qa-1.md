# plan-workflow-applies-rulings-qa-1

evidence · `01M7MSTQT9TRDWV649QQH1BTFP` · actor [claude-code-08e6675b](../actor/claude-code-08e6675b.md)

## result

PASS

## statement

QA: built binary, hand-run plan workflow: with a reviser that applies the rulings, check visit 3 passed to rulings, rulings visit 1 failed, apply-rulings visit 1 passed, rulings visit 2 passed, then gate. With a reviser that writes nothing, the run failed after rulings visit 2 (step apply-rulings over max_visits 1), with no gate step, and run status named both unapplied Decisions. go test ./internal/workflow -run 'TestPlanBuiltin|TestBuiltinPlanMilestone|TestBuiltinPlanReview' and the plan-workflow* scripts pass.

## Links

### Out

- verifies [plan-workflow-applies-rulings](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-workflow-applies-rulings.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
