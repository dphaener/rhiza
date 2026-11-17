# plan-workflow-applies-rulings-qa-7

evidence · `01M7MWP0Z7YJT2AD0T48T1HZYM` · actor [claude-code-6f696e6c](../actor/claude-code-6f696e6c.md)

## result

PASS

## statement

Scratch build of HEAD 94911c803 at /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza-fix-plan-run-arbiter-rulings/6f696e6c-9f58-4698-b518-01a1e9c78b82/scratchpad/bin/rhiza in two scratch autonomous plan runs (fake rex/claude/gh, own HOME/root): both runs went check (visit 3) pass -> rulings (visit 1) fail -> apply-rulings (visit 1) pass -> rulings (visit 2) pass -> gate pass; run B continued to merge and close-out and succeeded. go test ./cmd/rhiza -run TestScripts/(plan-workflow-autonomous-left|plan-workflow-autonomous-unapplied)$ passed.

## Links

### Out

- verifies [plan-workflow-applies-rulings](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-workflow-applies-rulings.md)

## History

- `01M7MPAHJZY3RV8ANGGB660J52` Fix Plan run arbiter rulings are never delivered: a plan run whose last review round leaves a confirmed major puts it to the arbiter, the arbiter answers with a decision handle (seen: parity-seed-bridge-in-model in run 01M7GG7GXSA97HESRNFEBFCZ4T, split-check-whatif-once in run 01M7GG7HSBSB356BVGGQGH8YHB), then the run merges and accepts the plan with the question marked 'decision ..., not yet delivered' and no Decision node or plan change in the model. The accepted plan still carries the confirmed major. Expected: the ruling is written (Decision + reviser applies it) before merge, or the run fails instead of accepting.
