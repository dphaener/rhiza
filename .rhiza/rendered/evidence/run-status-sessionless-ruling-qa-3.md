# run-status-sessionless-ruling-qa-3

evidence · `01M7MTBDKV9F9DMKMK5WZY6D9R` · actor [claude-code-60928770](../actor/claude-code-60928770.md)

## result

PASS

## statement

QA round 3 at db1b1c925, scratch build (go build -trimpath, run by full path, outside PATH) driving autonomous plan runs by hand in scratch Rhiza roots with the plan-workflow-autonomous* fake rex/claude/gh: after round 3, rhiza run status showed both review questions as '— decision settle-…' with no 'not yet delivered' and their --json had no session field; a session's question (planner's, plan-workflow-autonomous fakes) carried session in --json, and with its answered journal entry removed printed '— decision assume-streaming, not yet delivered'.

## Links

### Out

- verifies [run-status-sessionless-ruling](../capability/rhiza/workflows/built-in-workflows/plan-workflow/run-status-sessionless-ruling.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
