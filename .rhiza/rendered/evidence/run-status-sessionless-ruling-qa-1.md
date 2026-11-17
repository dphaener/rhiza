# run-status-sessionless-ruling-qa-1

evidence · `01M7MSTQC3MTJAN2BPSFMMGV7R` · actor [claude-code-08e6675b](../actor/claude-code-08e6675b.md)

## result

PASS

## statement

QA: built rhiza from 2c96751c5 to a scratch path and ran an autonomous plan run (fake claude/rex, real rhiza) whose round 3 left two review questions the arbiter passed. rhiza run status printed each as '— decision settle-…' with no 'not yet delivered', and run status --json questions carried no session field. go test ./internal/api -run TestRunInfoAcceptsAndQuestions and TestScripts/describe-golden pass.

## Links

### Out

- verifies [run-status-sessionless-ruling](../capability/rhiza/workflows/built-in-workflows/plan-workflow/run-status-sessionless-ruling.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
