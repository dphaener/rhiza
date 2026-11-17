# plan-reviser-applies-rulings-verified

evidence · `01M7MQT97AHP7KP0TARPAV1X9W` · actor [claude-code-2e407393](../actor/claude-code-2e407393.md)

## result

PASS

## statement

At 4d8c1fa98 re-ran the Checks: block: go test -p 2 -trimpath -count=1 ./internal/guide -run 'TestPlanReviserText' (PASS, exit 0), go test -p 2 -trimpath -count=1 ./internal/guide (ok, exit 0), rhiza ship generate (exit 0, 0 files changed), rhiza render (exit 0, no tracked change beyond concurrent sessions' actor index lines). Read the diff: plan-reviser.md gains '### Applying the arbiter's rulings' under Must do naming rhiza run status --json, rhiza apply in the change set holding the node with --record for a landed record, decided-by from the finding's node / the link's from node / each replacing node (every new leaf of a split), never edit, retire or re-link the Decision, change nothing a ruling does not name; TestPlanReviserText asserts each phrase in that subsection and the six ## headings in order. Confirmed-findings text and other guide topics untouched; no stubs or TODOs.

## Links

### Out

- verifies [plan-reviser-applies-rulings](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-reviser-applies-rulings.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
