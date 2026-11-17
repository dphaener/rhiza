# plan-reviser-applies-rulings-qa-2

evidence · `01M7MT0G5RG3A1Y00BZTHY9WG7` · actor [claude-code-fc7b7aa0](../actor/claude-code-fc7b7aa0.md)

## result

PASS

## statement

QA round 2 at db1b1c925, scratch build (go build -trimpath -o $SCRATCH/rhiza-x ./cmd/rhiza): rhiza-x guide plan-reviser prints the rulings section saying to link decided-by in the change set holding the question and its Decision, --record only for a landed record. Re-ran round 1's scenario in a hermetic scratch repo: landed record keep-flat, open change set with Decision settle-x; rhiza-x apply link keep-flat decided-by settle-x --record exits 4 (no node settle-x), the same link with the open change set exits 0 and rhiza-x why keep-flat lists settle-x as decided by; rhiza-x show settle-x names its change set (Created/View what-if). go test ./internal/guide -run TestPlanReviserText ok.

## Links

### Out

- verifies [plan-reviser-applies-rulings](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-reviser-applies-rulings.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
