# run-status-sessionless-ruling-qa-7

evidence · `01M7MWP16A322ZWAYRT79W5NA5` · actor [claude-code-6f696e6c](../actor/claude-code-6f696e6c.md)

## result

PASS

## statement

Scratch build of HEAD 94911c803 at /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza-fix-plan-run-arbiter-rulings/6f696e6c-9f58-4698-b518-01a1e9c78b82/scratchpad/bin/rhiza in two scratch autonomous plan runs (fake rex/claude/gh, own HOME/root): rhiza run status printed each session-less passed question as '— decision settle-...' with no 'not yet delivered', and --json questions carried no session field (null) and change only for questions outside the run's change set.

## Links

### Out

- verifies [run-status-sessionless-ruling](../capability/rhiza/workflows/built-in-workflows/plan-workflow/run-status-sessionless-ruling.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
