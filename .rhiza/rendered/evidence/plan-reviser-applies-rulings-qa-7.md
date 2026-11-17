# plan-reviser-applies-rulings-qa-7

evidence · `01M7MWP0V7J42AK0Q016GF6ZVD` · actor [claude-code-6f696e6c](../actor/claude-code-6f696e6c.md)

## result

PASS

## statement

Scratch build of HEAD 94911c803 at /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza-fix-plan-run-arbiter-rulings/6f696e6c-9f58-4698-b518-01a1e9c78b82/scratchpad/bin/rhiza in two scratch autonomous plan runs (fake rex/claude/gh, own HOME/root): a reviser following the guide's 'Applying the arbiter's rulings' literally (run status --json questions with decision; question change else run's change set; rhiza show <decision> --change <it>; apply --change <cs> / --record for the node; link decided-by in <it>) succeeded for a serving-draft leaf, a link the serving draft created, a landed node, a landed link and a leaf the ruling retired and split; every command exited 0.

## Links

### Out

- verifies [plan-reviser-applies-rulings](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-reviser-applies-rulings.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
