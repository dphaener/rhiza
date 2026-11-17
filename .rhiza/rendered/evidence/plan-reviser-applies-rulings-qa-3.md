# plan-reviser-applies-rulings-qa-3

evidence · `01M7MTBDSRZRH4GAB2W07CAWFZ` · actor [claude-code-60928770](../actor/claude-code-60928770.md)

## result

PASS

## statement

QA round 3 at db1b1c925, scratch build (go build -trimpath, run by full path, outside PATH) driving autonomous plan runs by hand in scratch Rhiza roots with the plan-workflow-autonomous* fake rex/claude/gh: rhiza guide plan-reviser prints the 'Applying the arbiter's rulings' section (run status --json, decided-by from the finding's node / link from node / each replacing node, --record for a landed record, never edit the Decision); go test ./internal/guide -run TestPlanReviserText ok. One minor: run status --json does not name the question's change set as the guide says.

## Links

### Out

- verifies [plan-reviser-applies-rulings](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-reviser-applies-rulings.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
